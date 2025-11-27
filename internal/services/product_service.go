package services

import (
	"context"
	"database/sql"
	"fmt"

	"ES/internal/domain"
	"ES/internal/ports"
)

type productService struct {
	db   *sql.DB
	repo ports.ProductRepository
}

func NewProductService(db *sql.DB, repo ports.ProductRepository) ports.ProductService {
	return &productService{db: db, repo: repo}
}

func (s *productService) UpdateProduct(ctx context.Context, id int64, name domain.BranchNameJSON) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.UpdateProduct(ctx, tx, id, name); err != nil {
		return fmt.Errorf("failed to update product in transaction: %w", err)
	}

	return tx.Commit()
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.DeleteProduct(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete product in transaction: %w", err)
	}

	return tx.Commit()
}
