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

func TestUpdateProduct_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &mockProductRepository{}
	service := NewProductService(db, repo)

	ctx := context.Background()
	testID := int64(1)
	testName := domain.BranchNameJSON{EN: "Test", TH: "ทดสอบ"}
	expectedError := errors.New("repository error")

	repo.UpdateProductFunc = func(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
		return expectedError
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = service.UpdateProduct(ctx, testID, testName)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mock Repository สำหรับ Product
type mockProductRepository struct {
	UpdateProductFunc func(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error
	DeleteProductFunc func(ctx context.Context, dbtx ports.DBTX, id int64) error
}

func (m *mockProductRepository) UpdateProduct(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(ctx, dbtx, id, name)
	}
	return nil
}

func (m *mockProductRepository) DeleteProduct(ctx context.Context, dbtx ports.DBTX, id int64) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(ctx, dbtx, id)
	}
	return nil
}
