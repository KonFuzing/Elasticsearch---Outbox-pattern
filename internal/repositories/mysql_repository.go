package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"ES/internal/domain"
	"ES/internal/ports"
)

// mySQLRepository คือ implementation ของ BranchRepository สำหรับ MySQL
type mySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository คือ factory function สำหรับสร้าง mySQLRepository
func NewMySQLRepository(db *sql.DB) *mySQLRepository {
	return &mySQLRepository{db: db}
}

// DB returns the underlying sql.DB object.
// This is useful for the service layer to start transactions.
func (r *mySQLRepository) DB() *sql.DB {
	return r.db
}

// CreateBranch เพิ่มข้อมูลสาขาใหม่ลงในตาราง `branch`
func (r *mySQLRepository) CreateBranch(ctx context.Context, dbtx ports.DBTX, name domain.BranchNameJSON) (int64, error) {
	jsonName, err := json.Marshal(name)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal branch name to JSON: %w", err)
	}

	query := "INSERT INTO branch (name) VALUES (?)"
	res, err := dbtx.ExecContext(ctx, query, string(jsonName))
	if err != nil {
		return 0, fmt.Errorf("failed to insert branch: %w", err)
	}

	return res.LastInsertId()
}

// LinkProductsToBranch เชื่อมโยงสาขากับสินค้าในตาราง `branches_products`
func (r *mySQLRepository) LinkProductsToBranch(ctx context.Context, dbtx ports.DBTX, branchID int64, productIDs []int) error {
	if len(productIDs) == 0 {
		return nil // ไม่มีสินค้าให้เชื่อมโยง
	}

	// สร้าง query แบบ bulk insert: INSERT INTO ... VALUES (?, ?), (?, ?), ...
	query := "INSERT INTO branches_products (branch_id, product_id) VALUES "
	var args []interface{}
	placeholders := []string{}

	for _, productID := range productIDs {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, branchID, productID)
	}

	query += strings.Join(placeholders, ", ")
	_, err := dbtx.ExecContext(ctx, query, args...)
	return err
}

// UpdateBranch อัปเดตข้อมูลชื่อของสาขา
func (r *mySQLRepository) UpdateBranch(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
	jsonName, err := json.Marshal(name)
	if err != nil {
		return fmt.Errorf("failed to marshal branch name to JSON: %w", err)
	}
	query := "UPDATE branch SET name = ? WHERE id = ?"
	_, err = dbtx.ExecContext(ctx, query, string(jsonName), id)
	return err
}

// DeleteBranch ลบข้อมูลสาขา
// เนื่องจากใน Schema มี ON DELETE CASCADE, ข้อมูลในตารางที่เกี่ยวข้องจะถูกลบไปด้วย
func (r *mySQLRepository) DeleteBranch(ctx context.Context, dbtx ports.DBTX, id int64) error {
	query := "DELETE FROM branch WHERE id = ?"
	_, err := dbtx.ExecContext(ctx, query, id)
	return err
}

// UnlinkAllProductsFromBranch ลบการเชื่อมโยงสินค้่าทั้งหมดของสาขา
func (r *mySQLRepository) UnlinkAllProductsFromBranch(ctx context.Context, dbtx ports.DBTX, branchID int64) error {
	query := "DELETE FROM branches_products WHERE branch_id = ?"
	_, err := dbtx.ExecContext(ctx, query, branchID)
	return err
}

// --- Interest ---
func (r *mySQLRepository) UpdateInterest(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
	jsonName, err := json.Marshal(name)
	if err != nil {
		return fmt.Errorf("failed to marshal interest name: %w", err)
	}
	query := "UPDATE interest SET name = ? WHERE id = ?"
	_, err = dbtx.ExecContext(ctx, query, string(jsonName), id)
	return err
}

func (r *mySQLRepository) DeleteInterest(ctx context.Context, dbtx ports.DBTX, id int64) error {
	query := "DELETE FROM interest WHERE id = ?"
	_, err := dbtx.ExecContext(ctx, query, id)
	return err
}

// --- Product ---
func (r *mySQLRepository) UpdateProduct(ctx context.Context, dbtx ports.DBTX, id int64, name domain.BranchNameJSON) error {
	jsonName, err := json.Marshal(name)
	if err != nil {
		return fmt.Errorf("failed to marshal product name: %w", err)
	}
	query := "UPDATE product SET name = ? WHERE id = ?"
	_, err = dbtx.ExecContext(ctx, query, string(jsonName), id)
	return err
}

func (r *mySQLRepository) DeleteProduct(ctx context.Context, dbtx ports.DBTX, id int64) error {
	query := "DELETE FROM product WHERE id = ?"
	_, err := dbtx.ExecContext(ctx, query, id)
	return err
}

// --- Product Option ---
func (r *mySQLRepository) UpdateProductOption(ctx context.Context, dbtx ports.DBTX, id int64, normalPrice, tagthaiPrice float64) error {
	query := "UPDATE product_option SET normal_price_thb = ?, tagthai_price_thb = ? WHERE id = ?"
	_, err := dbtx.ExecContext(ctx, query, normalPrice, tagthaiPrice, id)
	return err
}

func (r *mySQLRepository) DeleteProductOption(ctx context.Context, dbtx ports.DBTX, id int64) error {
	query := "DELETE FROM product_option WHERE id = ?"
	_, err := dbtx.ExecContext(ctx, query, id)
	return err
}

// --- Outbox ---
func (r *mySQLRepository) CreateEvent(ctx context.Context, dbtx ports.DBTX, aggregateID string, aggregateType string, eventType string, payload []byte) error {
	query := "INSERT INTO outbox_events (aggregate_id, aggregate_type, event_type, payload) VALUES (?, ?, ?, ?)"
	_, err := dbtx.ExecContext(ctx, query, aggregateID, aggregateType, eventType, payload)
	return err
}

// GetRichBranchData ดึงข้อมูลสาขาที่สมบูรณ์จากหลายตาราง
func (r *mySQLRepository) GetRichBranchData(ctx context.Context, dbtx ports.DBTX, id int64) (*domain.Branch, error) {
	// หมายเหตุ: Query นี้ยังขาดข้อมูล product_ids และมีการ join ที่อาจไม่ตรงกับ schema ปัจจุบัน
	// เราจะปรับปรุงให้ถูกต้อง
	const query = `
		SELECT
			branch.id,
			branch.name,
			ANY_VALUE(branch_location.province_id) AS province_id,
			(SELECT GROUP_CONCAT(DISTINCT p.product_id) FROM branches_products p WHERE p.branch_id = branch.id) AS product_ids,
			(SELECT GROUP_CONCAT(DISTINCT i.interest_id) FROM branches_interests i WHERE i.branch_id = branch.id) AS interest_ids
		FROM
			branch
		LEFT JOIN
			branch_location ON branch.id = branch_location.branch_id
		WHERE
			branch.id = ?
		GROUP BY
			branch.id;
	`

	row := dbtx.QueryRowContext(ctx, query, id)

	var branch domain.Branch
	var nameJSON, productIDsStr, interestIDsStr sql.NullString
	var provinceID sql.NullInt64

	err := row.Scan(
		&branch.ID,
		&nameJSON,
		&provinceID,
		&productIDsStr,
		&interestIDsStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("branch with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to scan rich branch data: %w", err)
	}

	// แปลงข้อมูล JSON และ String ที่ได้จาก DB
	if nameJSON.Valid {
		if err := json.Unmarshal([]byte(nameJSON.String), &branch.Name); err != nil {
			log.Printf("WARNING: could not unmarshal branch name for id %d: %v", id, err)
		}
	}

	if provinceID.Valid {
		branch.Location = &domain.BranchLocation{ProvinceID: int(provinceID.Int64)}
	}

	if productIDsStr.Valid && productIDsStr.String != "" {
		ids := strings.Split(productIDsStr.String, ",")
		branch.ProductIDs = make([]int, 0, len(ids))
		for _, idStr := range ids {
			id, _ := strconv.Atoi(idStr)
			branch.ProductIDs = append(branch.ProductIDs, id)
		}
	}

	if interestIDsStr.Valid && interestIDsStr.String != "" {
		ids := strings.Split(interestIDsStr.String, ",")
		branch.InterestIDs = make([]int, 0, len(ids))
		for _, idStr := range ids {
			id, _ := strconv.Atoi(idStr)
			branch.InterestIDs = append(branch.InterestIDs, id)
		}
	}

	return &branch, nil
}

// GetAllRichBranchData ดึงข้อมูลสาขาที่สมบูรณ์ทั้งหมดใน query เดียวเพื่อทำ backfill
func (r *mySQLRepository) GetAllRichBranchData(ctx context.Context, dbtx ports.DBTX) ([]*domain.Branch, error) {
	const query = `
		SELECT
			branch.id,
			branch.name,
			ANY_VALUE(branch_location.province_id) AS province_id,
			(SELECT GROUP_CONCAT(DISTINCT p.product_id) FROM branches_products p WHERE p.branch_id = branch.id) AS product_ids,
			(SELECT GROUP_CONCAT(DISTINCT i.interest_id) FROM branches_interests i WHERE i.branch_id = branch.id) AS interest_ids
		FROM
			branch
		LEFT JOIN
			branch_location ON branch.id = branch_location.branch_id
		GROUP BY
			branch.id
		ORDER BY
			branch.id ASC;
	`

	rows, err := dbtx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all rich branch data: %w", err)
	}
	defer rows.Close()

	var branches []*domain.Branch
	for rows.Next() {
		var branch domain.Branch
		var nameJSON, productIDsStr, interestIDsStr sql.NullString
		var provinceID sql.NullInt64

		if err := rows.Scan(&branch.ID, &nameJSON, &provinceID, &productIDsStr, &interestIDsStr); err != nil {
			log.Printf("WARNING: could not scan row for rich branch data: %v", err)
			continue // ข้ามแถวที่มีปัญหา
		}

		if nameJSON.Valid {
			_ = json.Unmarshal([]byte(nameJSON.String), &branch.Name)
		}
		if provinceID.Valid {
			branch.Location = &domain.BranchLocation{ProvinceID: int(provinceID.Int64)}
		}
		if productIDsStr.Valid && productIDsStr.String != "" {
			ids := strings.Split(productIDsStr.String, ",")
			branch.ProductIDs = make([]int, len(ids))
			for i, idStr := range ids {
				branch.ProductIDs[i], _ = strconv.Atoi(idStr)
			}
		}
		if interestIDsStr.Valid && interestIDsStr.String != "" {
			ids := strings.Split(interestIDsStr.String, ",")
			branch.InterestIDs = make([]int, len(ids))
			for i, idStr := range ids {
				branch.InterestIDs[i], _ = strconv.Atoi(idStr)
			}
		}
		branches = append(branches, &branch)
	}
	return branches, nil
}
