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

	// --- 4. Get ALL rich branch data in a SINGLE query ---
	log.Println("Fetching all rich branch data from MySQL in a single query...")
	allBranches, err := repo.GetAllRichBranchData(ctx, db)
	if err != nil {
		log.Fatalf("Failed to get all rich branch data: %v", err)
	}

	if len(allBranches) == 0 {
		log.Println("No branches found in the database. Exiting.")
		return
	}

	log.Printf("Found %d branches to backfill.", len(allBranches))

	// --- 5. Loop through the results and index to Elasticsearch ---
	successCount := 0
	for _, branchData := range allBranches {
		log.Printf("Indexing branch ID: %d...", branchData.ID)

		// ส่งข้อมูลไป Index ใน Elasticsearch
		_, err = esClient.Index().
			Index("branches").
			Id(strconv.FormatInt(branchData.ID, 10)).
			BodyJson(branchData).
			Do(ctx)

		if err != nil {
			log.Printf("ERROR: Failed to index branch ID %d to Elasticsearch: %v. Skipping.", branchData.ID, err)
			continue
		}
		successCount++
	}

	log.Printf("--- Backfill Process Finished ---")
	log.Printf("Successfully backfilled %d out of %d branches.", successCount, len(allBranches))
	log.Printf("Total time taken: %v", time.Since(startTime))
}
