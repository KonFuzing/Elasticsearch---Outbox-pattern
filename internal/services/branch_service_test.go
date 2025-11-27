package services

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"ES/internal/domain"
	"ES/internal/ports"
	"ES/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBranchWithProducts_RollbackOnLinkError(t *testing.T) {
	// 1. --- Setup ---
	// สร้าง Mock Database และ Mock Redis
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	redisClient, redisMock := redismock.NewClientMock()

	// สร้าง Repository และ Service ที่จะทดสอบ
	// สังเกตว่าเราใช้ repo จริง แต่ส่ง mock db เข้าไป
	repo := repositories.NewMySQLRepository(db)
	var branchRepo ports.BranchRepository = repo
	var outboxRepo ports.OutboxRepository = repo

	service := NewBranchService(db, branchRepo, outboxRepo, redisClient)

	// กำหนดค่าสำหรับ Test
	branchID := int64(1)
	branchName := domain.BranchNameJSON{EN: "Test Branch", TH: "สาขาทดสอบ"}
	productIDs := []int{101, 102}
	expectedError := errors.New("simulated link product error")

	// 2. --- กำหนด Expectations ของ Mock ---
	// เราคาดหวังว่าโค้ดจะทำงานตามลำดับนี้:
	mock.ExpectBegin() // 1. เริ่ม Transaction

	// 2. อัปเดตข้อมูล Branch (คาดว่าจะสำเร็จ)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE branch SET name = ? WHERE id = ?")).
		WithArgs(`{"en":"Test Branch","th":"สาขาทดสอบ"}`, branchID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 3. ลบการเชื่อมโยงสินค้าเก่า (คาดว่าจะสำเร็จ)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM branches_products WHERE branch_id = ?")).
		WithArgs(branchID).
		WillReturnResult(sqlmock.NewResult(0, 5)) // สมมติว่าลบไป 5 รายการ

	// 4. สร้างการเชื่อมโยงสินค้าใหม่ (จำลองให้ขั้นตอนนี้ "ล้มเหลว")
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO branches_products (branch_id, product_id) VALUES (?, ?), (?, ?)")).
		WithArgs(branchID, productIDs[0], branchID, productIDs[1]).
		WillReturnError(expectedError)

	// 5. คาดหวังว่าจะมีการเรียก Rollback!
	mock.ExpectRollback()

	// เราไม่คาดหวังว่าจะมีการเรียกใช้ GetRichBranchData, CreateEvent, Commit, หรือ Redis Publish
	// เพราะโค้ดควรจะล้มเหลวและ Rollback ไปก่อน
	redisMock.ExpectPublish("outbox_channel", "new_event").SetVal(0) // ตั้งค่าเผื่อไว้ แต่ไม่ควรถูกเรียก

	// 3. --- เรียกใช้ฟังก์ชันที่ต้องการทดสอบ ---
	ctx := context.Background()
	_, err = service.UpdateBranchWithProducts(ctx, branchID, branchName, productIDs)

	// 4. --- ตรวจสอบผลลัพธ์ ---
	// ตรวจสอบว่าฟังก์ชัน return error กลับมาจริง
	assert.Error(t, err)
	// ตรวจสอบว่า error message ที่ได้ มีข้อความจาก error ที่เราจำลองขึ้น
	assert.Contains(t, err.Error(), "failed to link new products")
	assert.Contains(t, err.Error(), expectedError.Error())

	// ตรวจสอบว่า Mock Expectations ทั้งหมดถูกเรียกใช้ครบถ้วนและถูกต้องตามลำดับ
	// นี่คือการยืนยันว่า Begin, Exec, Exec, Exec(Error), และ Rollback เกิดขึ้นจริง
	assert.NoError(t, mock.ExpectationsWereMet())
}
