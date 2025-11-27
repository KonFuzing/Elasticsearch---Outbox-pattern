package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic/v7" // ต้อง go get package นี้
)

type OutboxEvent struct {
	ID            int64
	AggregateID   string
	AggregateType string
	EventType     string
	Payload       []byte
}

func main() {
	// --- 1. Connect to MySQL ---
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/TTDB?parseTime=true"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	// --- 2. Connect to Elasticsearch ---
	esClient, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"), // URL ของ Elasticsearch
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}
	log.Println("Successfully connected to Elasticsearch.")

	// --- 3. Connect to Redis and Subscribe ---
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")

	pubsub := redisClient.Subscribe(context.Background(), "outbox_channel")
	defer pubsub.Close()

	// --- 4. Start Worker ---
	log.Println("Worker started. Waiting for notifications on 'outbox_channel'...")

	// ประมวลผลครั้งแรกเผื่อมี event ค้างอยู่ตอน worker ปิดไป
	processEvents(db, esClient)

	// รอรับ message จาก channel
	for msg := range pubsub.Channel() {
		log.Printf("Received notification: %s. Triggering event processing.", msg.Payload)
		processEvents(db, esClient)
	}
}

func processEvents(db *sql.DB, esClient *elastic.Client) {
	log.Println("--- Checking for new events... ---")
	ctx := context.Background()

	// 1. ดึง Events ที่เป็น pending
	rows, err := db.QueryContext(ctx, "SELECT id, aggregate_id, aggregate_type, event_type, payload FROM outbox_events WHERE status = 'pending' ORDER BY created_at ASC LIMIT 10")
	if err != nil {
		log.Printf("Error querying events: %v", err)
		return
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.ID, &event.AggregateID, &event.AggregateType, &event.EventType, &event.Payload); err != nil {
			log.Printf("Error scanning event row: %v", err)
			continue
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		log.Println("--- No new events found. ---")
		return
	}

	log.Printf("Found %d new events to process.", len(events))

	// 2. ประมวลผลแต่ละ Event
	for _, event := range events {
		err := handleEvent(ctx, esClient, event)

		// 3. อัปเดตสถานะ Event
		var newStatus string
		if err != nil {
			log.Printf("Failed to process event ID %d: %v", event.ID, err)
			newStatus = "failed"
		} else {
			log.Printf("Successfully processed event ID %d", event.ID)
			newStatus = "processed"
		}

		_, updateErr := db.ExecContext(ctx, "UPDATE outbox_events SET status = ? WHERE id = ?", newStatus, event.ID)
		if updateErr != nil {
			log.Printf("CRITICAL: Failed to update status for event ID %d: %v", event.ID, updateErr)
		}
	}
}

func handleEvent(ctx context.Context, esClient *elastic.Client, event OutboxEvent) error {
	// Logic การส่งข้อมูลไป Elasticsearch
	// ในตัวอย่างนี้ เราจะจัดการเฉพาะ "branch"
	if event.AggregateType != "branch" {
		return fmt.Errorf("unhandled aggregate type: %s", event.AggregateType)
	}

	indexName := "branches" // ชื่อ index ใน Elasticsearch

	switch event.EventType {
	case "created", "updated":
		var payloadData map[string]interface{}
		if err := json.Unmarshal(event.Payload, &payloadData); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		_, err := esClient.Index().
			Index(indexName).
			Id(event.AggregateID).
			BodyJson(payloadData).
			Do(ctx)
		return err

	case "deleted":
		_, err := esClient.Delete().
			Index(indexName).
			Id(event.AggregateID).
			Do(ctx)
		// สำหรับการลบ ถ้าไม่เจอ (404) ก็ถือว่าสำเร็จ
		if err != nil && elastic.IsNotFound(err) {
			log.Printf("Document with ID %s already deleted. Considering it a success.", event.AggregateID)
			return nil
		}
		return err

	default:
		return fmt.Errorf("unhandled event type: %s", event.EventType)
	}
}
