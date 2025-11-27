package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"ES/internal/domain"
	"ES/internal/ports"

	"github.com/go-redis/redis/v8"
)

// branchService คือ implementation ของ BranchService
type branchService struct {
	db          *sql.DB
	branchRepo  ports.BranchRepository
	outboxRepo  ports.OutboxRepository
	redisClient *redis.Client
}

// NewBranchService คือ factory function สำหรับสร้าง branchService
// สังเกตว่าเรารับ *sql.DB เข้ามาด้วยเพื่อใช้จัดการ Transaction
func NewBranchService(db *sql.DB, branchRepo ports.BranchRepository, outboxRepo ports.OutboxRepository, redisClient *redis.Client) ports.BranchService {
	return &branchService{
		db:          db,
		branchRepo:  branchRepo,
		outboxRepo:  outboxRepo,
		redisClient: redisClient,
	}
}

// CreateBranchWithProducts คือเมธอดที่จัดการ business logic ทั้งหมดใน transaction เดียว
func (s *branchService) CreateBranchWithProducts(ctx context.Context, name domain.BranchNameJSON, productIDs []int) (*domain.Branch, error) {
	// 1. เริ่มต้น Transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback อัตโนมัติหากเกิด error และฟังก์ชันจบการทำงาน

	// 2. สร้างสาขา โดยส่ง `tx` เข้าไปให้ Repository
	branchID, err := s.branchRepo.CreateBranch(ctx, tx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch in transaction: %w", err)
	}

	// 3. เชื่อมโยงสินค้า โดยส่ง `tx` ตัวเดียวกันเข้าไป
	if err := s.branchRepo.LinkProductsToBranch(ctx, tx, branchID, productIDs); err != nil {
		return nil, fmt.Errorf("failed to link products in transaction: %w", err)
	}

	// 4. ถ้าทุกอย่างสำเร็จ ให้ Commit Transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &domain.Branch{ID: branchID, Name: name, ProductIDs: productIDs}, nil
}

// UpdateBranchWithProducts อัปเดตข้อมูลสาขาและสินค้าที่เชื่อมโยง
func (s *branchService) UpdateBranchWithProducts(ctx context.Context, id int64, name domain.BranchNameJSON, productIDs []int) (*domain.Branch, error) {
	log.Printf("Starting transaction to update branch ID: %d", id)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("ERROR: Failed to begin transaction for branch ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// เราจะจัดการ Rollback และ Commit เองเพื่อเพิ่ม Log ให้ชัดเจน
	// defer tx.Rollback()

	// 1. อัปเดตข้อมูลพื้นฐานของสาขา
	log.Printf("Step 1: Updating branch info for ID: %d", id)
	if err := s.branchRepo.UpdateBranch(ctx, tx, id, name); err != nil {
		log.Printf("ERROR: Step 1 failed. Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to update branch in transaction: %w", err)
	}

	// 2. ลบการเชื่อมโยงสินค้าเก่าทั้งหมด
	log.Printf("Step 2: Unlinking all old products from branch ID: %d", id)
	if err := s.branchRepo.UnlinkAllProductsFromBranch(ctx, tx, id); err != nil {
		log.Printf("ERROR: Step 2 failed. Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to unlink old products in transaction: %w", err)
	}

	// 3. สร้างการเชื่อมโยงสินค้าใหม่
	log.Printf("Step 3: Linking new products to branch ID: %d", id)
	if err := s.branchRepo.LinkProductsToBranch(ctx, tx, id, productIDs); err != nil {
		log.Printf("ERROR: Step 3 failed. Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to link new products in transaction: %w", err)
	}

	// 4. สร้าง Event สำหรับ Outbox
	log.Printf("Step 4: Creating outbox event for branch ID: %d", id)
	// ดึงข้อมูลฉบับสมบูรณ์ล่าสุดจาก DB เพื่อสร้าง payload
	richBranchData, err := s.branchRepo.GetRichBranchData(ctx, tx, id)
	if err != nil {
		log.Printf("ERROR: Step 4 failed (GetRichBranchData). Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to get rich branch data for outbox: %w", err)
	}
	payload, err := json.Marshal(richBranchData)
	if err != nil {
		log.Printf("ERROR: Step 4 failed (JSON Marshal). Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to marshal payload for outbox: %w", err)
	}
	if err := s.outboxRepo.CreateEvent(ctx, tx, strconv.FormatInt(id, 10), "branch", "updated", payload); err != nil {
		log.Printf("ERROR: Step 4 failed (Create Event). Rolling back transaction. Error: %v", err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create outbox event: %w", err)
	}

	// 5. Commit Transaction
	log.Printf("All steps successful. Committing transaction for branch ID: %d", id)
	if err := tx.Commit(); err != nil {
		log.Printf("ERROR: Failed to commit transaction. Error: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// หลังจาก Commit สำเร็จ ให้ส่ง Notification
	log.Println("Transaction committed. Publishing notification to 'outbox_channel'.")
	if err := s.redisClient.Publish(ctx, "outbox_channel", "new_event").Err(); err != nil {
		// การส่ง notification ล้มเหลวไม่ควรกระทบ logic หลัก แต่ควร log ไว้
		log.Printf("WARNING: Failed to publish notification to Redis: %v", err)
	}

	return richBranchData, nil
}

// DeleteBranch ลบสาขา
func (s *branchService) DeleteBranch(ctx context.Context, id int64) error {
	// การลบไม่จำเป็นต้องใช้ transaction ที่ซับซ้อน เพราะ ON DELETE CASCADE ใน DB จะจัดการให้
	// แต่เพื่อให้สามารถสร้าง Outbox Event ได้อย่างปลอดภัย เราจะทำใน Transaction
	log.Printf("Starting transaction to delete branch ID: %d", id)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for delete: %w", err)
	}
	defer tx.Rollback()

	// 1. สร้าง Event "deleted" ก่อน
	// Payload สำหรับการลบอาจไม่จำเป็นต้องมีข้อมูลเต็ม แค่ ID ก็เพียงพอ
	payload, _ := json.Marshal(map[string]int64{"id": id})
	if err := s.outboxRepo.CreateEvent(ctx, tx, strconv.FormatInt(id, 10), "branch", "deleted", payload); err != nil {
		return fmt.Errorf("failed to create 'deleted' event for outbox: %w", err)
	}

	// 2. ทำการลบข้อมูลจริง
	if err := s.branchRepo.DeleteBranch(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete branch in transaction: %w", err)
	}

	// 3. Commit Transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for delete: %w", err)
	}

	// 4. ส่ง Notification
	log.Println("Delete transaction committed. Publishing notification to 'outbox_channel'.")
	s.redisClient.Publish(ctx, "outbox_channel", "new_event") // ไม่ต้องเช็ค error เพื่อไม่ให้กระทบ flow หลัก

	return nil
}

// GetBranch ดึงข้อมูลสาขาแบบสมบูรณ์
func (s *branchService) GetBranch(ctx context.Context, id int64) (*domain.Branch, error) {
	// ใช้ DB connection ปกติ ไม่จำเป็นต้องใช้ transaction สำหรับการอ่าน
	return s.branchRepo.GetRichBranchData(ctx, s.db, id)
}
