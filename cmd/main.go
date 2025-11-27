package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"ES/internal/handlers"     // Driving Adapter
	"ES/internal/ports"        // Ports
	"ES/internal/repositories" // Driven Adapter
	"ES/internal/services"     // Core Logic
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// --- 1. ตั้งค่า Database Connection ---
	// อ่านค่า DSN จาก Environment Variable เพื่อความยืดหยุ่น
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// กำหนดค่า default สำหรับการพัฒนาในเครื่อง
		dsn = "root:123456@tcp(127.0.0.1:3306)/TTDB?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database.")

	// --- 1.5. ตั้งค่า Redis Connection ---
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // ที่อยู่ของ Redis
		Password: "",               // ไม่มีรหัสผ่าน
		DB:       0,                // ใช้ DB default
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis.")

	// --- 2. Dependency Injection (ประกอบร่าง) ---
	// Repository (Driven Adapter) -> Service (Core) -> Handler (Driving Adapter)

	// สร้าง Repository Adapter โดยส่ง connection (db) เข้าไป
	// สังเกตว่า NewMySQLRepository คืนค่าเป็น struct เพื่อให้เราสามารถ cast เป็น interface ย่อยๆ ได้
	repo := repositories.NewMySQLRepository(db)
	var branchRepo ports.BranchRepository = repo
	var interestRepo ports.InterestRepository = repo
	var productRepo ports.ProductRepository = repo
	var productOptionRepo ports.ProductOptionRepository = repo
	var outboxRepo ports.OutboxRepository = repo

	// สร้าง Service โดยส่ง db (สำหรับ transaction) และ Repository เข้าไป
	var branchSvc ports.BranchService = services.NewBranchService(db, branchRepo, outboxRepo, redisClient)
	var interestSvc ports.InterestService = services.NewInterestService(db, interestRepo)
	var productSvc ports.ProductService = services.NewProductService(db, productRepo)
	var productOptionSvc ports.ProductOptionService = services.NewProductOptionService(db, productOptionRepo)

	// สร้าง Handler โดยส่ง Service เข้าไป
	httpHandler := handlers.NewHTTPHandler(branchSvc, interestSvc, productSvc, productOptionSvc)

	// --- 3. ตั้งค่า Gin Router ---
	router := gin.Default()

	// กำหนด endpoint
	branchRoutes := router.Group("/branches")
	{
		branchRoutes.POST("/", httpHandler.CreateBranch) // Create ยังคงอยู่
		branchRoutes.GET("/:id", httpHandler.GetBranch)
		branchRoutes.PUT("/:id", httpHandler.UpdateBranch)
		branchRoutes.DELETE("/:id", httpHandler.DeleteBranch)
	}

	interestRoutes := router.Group("/interests")
	{
		interestRoutes.PUT("/:id", httpHandler.UpdateInterest)
		interestRoutes.DELETE("/:id", httpHandler.DeleteInterest)
	}

	productRoutes := router.Group("/products")
	{
		productRoutes.PUT("/:id", httpHandler.UpdateProduct)
		productRoutes.DELETE("/:id", httpHandler.DeleteProduct)
	}

	productOptionRoutes := router.Group("/product-options")
	{
		productOptionRoutes.PUT("/:id", httpHandler.UpdateProductOption)
		productOptionRoutes.DELETE("/:id", httpHandler.DeleteProductOption)
	}

	// --- 4. รันเซิร์ฟเวอร์ ---
	fmt.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
