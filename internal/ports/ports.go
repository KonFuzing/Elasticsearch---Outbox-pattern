package ports

import (
	"context"
	"database/sql"

	"ES/internal/domain"
)

// DBTX เป็น interface ที่ครอบคลุมทั้ง *sql.DB และ *sql.Tx
// ทำให้ Repository method สามารถทำงานได้ทั้งในและนอก transaction
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// BranchRepository คือ port สำหรับการติดต่อกับฐานข้อมูลของ Branch
type BranchRepository interface {
	CreateBranch(ctx context.Context, dbtx DBTX, name domain.BranchNameJSON) (int64, error)
	LinkProductsToBranch(ctx context.Context, dbtx DBTX, branchID int64, productIDs []int) error
	UpdateBranch(ctx context.Context, dbtx DBTX, id int64, name domain.BranchNameJSON) error
	DeleteBranch(ctx context.Context, dbtx DBTX, id int64) error
	UnlinkAllProductsFromBranch(ctx context.Context, dbtx DBTX, branchID int64) error
	GetRichBranchData(ctx context.Context, dbtx DBTX, id int64) (*domain.Branch, error)
}

// OutboxRepository คือ port สำหรับการเขียน event
type OutboxRepository interface {
	CreateEvent(ctx context.Context, dbtx DBTX, aggregateID string, aggregateType string, eventType string, payload []byte) error
}

// BranchService คือ port สำหรับ business logic ของ Branch
type BranchService interface {
	CreateBranchWithProducts(ctx context.Context, name domain.BranchNameJSON, productIDs []int) (*domain.Branch, error)
	UpdateBranchWithProducts(ctx context.Context, id int64, name domain.BranchNameJSON, productIDs []int) (*domain.Branch, error)
	DeleteBranch(ctx context.Context, id int64) error
	GetBranch(ctx context.Context, id int64) (*domain.Branch, error)
}

// InterestRepository คือ port สำหรับ Interest
type InterestRepository interface {
	UpdateInterest(ctx context.Context, dbtx DBTX, id int64, name domain.BranchNameJSON) error
	DeleteInterest(ctx context.Context, dbtx DBTX, id int64) error
}

// ProductRepository คือ port สำหรับ Product
type ProductRepository interface {
	UpdateProduct(ctx context.Context, dbtx DBTX, id int64, name domain.BranchNameJSON) error
	DeleteProduct(ctx context.Context, dbtx DBTX, id int64) error
}

// ProductOptionRepository คือ port สำหรับ ProductOption
type ProductOptionRepository interface {
	UpdateProductOption(ctx context.Context, dbtx DBTX, id int64, normalPrice, tagthaiPrice float64) error
	DeleteProductOption(ctx context.Context, dbtx DBTX, id int64) error
}

// InterestService คือ port สำหรับ business logic ของ Interest
type InterestService interface {
	UpdateInterest(ctx context.Context, id int64, name domain.BranchNameJSON) error
	DeleteInterest(ctx context.Context, id int64) error
}

// ProductService คือ port สำหรับ business logic ของ Product
type ProductService interface {
	UpdateProduct(ctx context.Context, id int64, name domain.BranchNameJSON) error
	DeleteProduct(ctx context.Context, id int64) error
}

// ProductOptionService คือ port สำหรับ business logic ของ ProductOption
type ProductOptionService interface {
	UpdateProductOption(ctx context.Context, id int64, normalPrice, tagthaiPrice float64) error
	DeleteProductOption(ctx context.Context, id int64) error
}
