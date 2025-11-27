package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("--- Starting Outbox Cleanup Process ---")

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

	// --- 2. กำหนดระยะเวลาที่จะเก็บข้อมูลไว้ (เช่น 7 วัน) ---
	retentionDays := 7
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	log.Printf("Deleting processed events older than %d days (before %s)", retentionDays, cutoffDate.Format("2006-01-02"))

	// --- 3. รันคำสั่ง DELETE ---
	query := "DELETE FROM outbox_events WHERE status = 'processed' AND created_at < ?"
	result, err := db.ExecContext(context.Background(), query, cutoffDate)
	if err != nil {
		log.Fatalf("Failed to execute delete query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	log.Printf("--- Cleanup Process Finished. Deleted %d old events. ---", rowsAffected)
}
