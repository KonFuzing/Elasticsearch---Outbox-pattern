package services

import (
	"ES/internal/ports"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProductOption_RollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &mockProductOptionRepository{}
	service := NewProductOptionService(db, repo)

	ctx := context.Background()
	testID := int64(1)
	expectedError := errors.New("repository error")

	repo.UpdateProductOptionFunc = func(ctx context.Context, dbtx ports.DBTX, id int64, normalPrice, tagthaiPrice float64) error {
		return expectedError
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = service.UpdateProductOption(ctx, testID, 100, 90)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mock Repository สำหรับ ProductOption
type mockProductOptionRepository struct {
	UpdateProductOptionFunc func(ctx context.Context, dbtx ports.DBTX, id int64, normalPrice, tagthaiPrice float64) error
	DeleteProductOptionFunc func(ctx context.Context, dbtx ports.DBTX, id int64) error
}

func (m *mockProductOptionRepository) UpdateProductOption(ctx context.Context, dbtx ports.DBTX, id int64, normalPrice, tagthaiPrice float64) error {
	if m.UpdateProductOptionFunc != nil {
		return m.UpdateProductOptionFunc(ctx, dbtx, id, normalPrice, tagthaiPrice)
	}
	return nil
}

func (m *mockProductOptionRepository) DeleteProductOption(ctx context.Context, dbtx ports.DBTX, id int64) error {
	if m.DeleteProductOptionFunc != nil {
		return m.DeleteProductOptionFunc(ctx, dbtx, id)
	}
	return nil
}
