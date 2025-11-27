package services

import (
	"context"
	"errors"
	"testing"

	"ES/internal/domain"
	"ES/internal/ports"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateInterest_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &mockInterestRepository{}
	service := NewInterestService(db, repo)

	ctx := context.Background()
	testID := int64(1)
	testName := domain.BranchNameJSON{EN: "Test", TH: "ทดสอบ"}
	expectedError := errors.New("repository error")

	// ตั้งค่า mock repository ให้คืนค่า error
	repo.UpdateInterestFunc = func(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
		return expectedError
	}

	// ตั้งค่า mock database transaction
	mock.ExpectBegin()
	mock.ExpectRollback() // คาดหวังว่าจะมีการ Rollback

	// เรียกใช้ service
	err = service.UpdateInterest(ctx, testID, testName)

	// ตรวจสอบผลลัพธ์
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mock Repository สำหรับ Interest
type mockInterestRepository struct {
	UpdateInterestFunc func(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error
	DeleteInterestFunc func(ctx context.Context, dbtx ports.DBTX, id int64) error
}

func (m *mockInterestRepository) UpdateInterest(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
	if m.UpdateInterestFunc != nil {
		return m.UpdateInterestFunc(ctx, dbtx, id, name)
	}
	return nil
}

func (m *mockInterestRepository) DeleteInterest(ctx context.Context, dbtx ports.DBTX, id int64) error {
	if m.DeleteInterestFunc != nil {
		return m.DeleteInterestFunc(ctx, dbtx, id)
	}
	return nil
}
