package services

import (
	"context"
	"database/sql"
	"fmt"

	"ES/internal/ports"
)

type productOptionService struct {
	db   *sql.DB
	repo ports.ProductOptionRepository
}

func NewProductOptionService(db *sql.DB, repo ports.ProductOptionRepository) ports.ProductOptionService {
	return &productOptionService{db: db, repo: repo}
}

func (s *productOptionService) UpdateProductOption(ctx context.Context, id int64, normalPrice, tagthaiPrice float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.UpdateProductOption(ctx, tx, id, normalPrice, tagthaiPrice); err != nil {
		return fmt.Errorf("failed to update product option in transaction: %w", err)
	}

	return tx.Commit()
}

func (s *productOptionService) DeleteProductOption(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.DeleteProductOption(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product option in transaction: %w", err)
	}

	return tx.Commit()
}
