package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"ES/internal/repositories"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic/v7"
)

func main() {
	log.Println("--- Starting Backfill Process ---")
	startTime := time.Now()

	ctx := context.Background()

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
	log.Println("Successfully connected to MySQL.")

	// --- 2. Connect to Elasticsearch ---
	esClient, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}
	log.Println("Successfully connected to Elasticsearch.")

	// --- 3. Initialize Repository ---
	// เราจะใช้ Repository ที่มีอยู่แล้วเพื่อดึงข้อมูล
	repo := repositories.NewMySQLRepository(db)

	// --- 4. Get all Branch IDs from MySQL ---
	rows, err := db.QueryContext(ctx, "SELECT id FROM branch ORDER BY id ASC")
	if err != nil {
		log.Fatalf("Failed to query branch IDs: %v", err)
	}
	defer rows.Close()

	var branchIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("Warning: Failed to scan branch ID: %v", err)
			continue
		}
		branchIDs = append(branchIDs, id)
	}

	if len(branchIDs) == 0 {
		log.Println("No branches found in the database. Exiting.")
		return
	}

	log.Printf("Found %d branches to backfill.", len(branchIDs))

	// --- 5. Loop through each ID, get rich data, and index to Elasticsearch ---
	successCount := 0
	for _, id := range branchIDs {
		log.Printf("Backfilling branch ID: %d...", id)

		// ใช้ฟังก์ชันที่เรามีอยู่แล้วเพื่อดึงข้อมูลที่สมบูรณ์
		richBranchData, err := repo.GetRichBranchData(ctx, db, id)
		if err != nil {
			log.Printf("ERROR: Could not get rich data for branch ID %d: %v. Skipping.", id, err)
			continue
		}

		// ส่งข้อมูลไป Index ใน Elasticsearch
		_, err = esClient.Index().
			Index("branches").
			Id(strconv.FormatInt(id, 10)).
			BodyJson(richBranchData).
			Do(ctx)

		if err != nil {
			log.Printf("ERROR: Failed to index branch ID %d to Elasticsearch: %v. Skipping.", id, err)
			continue
		}
		successCount++
	}

	log.Printf("--- Backfill Process Finished ---")
	log.Printf("Successfully backfilled %d out of %d branches.", successCount, len(branchIDs))
	log.Printf("Total time taken: %v", time.Since(startTime))
}
