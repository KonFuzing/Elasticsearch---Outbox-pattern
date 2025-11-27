package services

import (
	"context"
	"database/sql"
	"fmt"

	"ES/internal/domain"
	"ES/internal/ports"
)

type interestService struct {
	db   *sql.DB
	repo ports.InterestRepository
}

func NewInterestService(db *sql.DB, repo ports.InterestRepository) ports.InterestService {
	return &interestService{db: db, repo: repo}
}

func (s *interestService) UpdateInterest(ctx context.Context, id int64, name domain.BranchNameJSON) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.UpdateInterest(ctx, tx, id, name); err != nil {
		return fmt.Errorf("failed to update interest in transaction: %w", err)
	}

	return tx.Commit()
}

func (s *interestService) DeleteInterest(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.repo.DeleteInterest(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete interest in transaction: %w", err)
	}

	return tx.Commit()
}
