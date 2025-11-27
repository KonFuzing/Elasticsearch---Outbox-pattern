package handlers

import (
	"log"
	"net/http"
	"strconv"

	"ES/internal/domain"
	"ES/internal/ports"

	"github.com/gin-gonic/gin"
)

// CreateBranchRequest คือ struct สำหรับรับข้อมูล JSON จาก request body
type CreateBranchRequest struct {
	Name       domain.BranchNameJSON `json:"name" binding:"required"`
	ProductIDs []int                 `json:"product_ids"`
}

// UpdateNameRequest เป็น struct กลางสำหรับรับข้อมูล JSON ที่มีแค่ name
type UpdateNameRequest struct {
	Name domain.BranchNameJSON `json:"name" binding:"required"`
}

// UpdateProductOptionRequest คือ struct สำหรับรับข้อมูล JSON ของ product_option
type UpdateProductOptionRequest struct {
	NormalPrice  float64 `json:"normal_price_thb" binding:"required"`
	TagthaiPrice float64 `json:"tagthai_price_thb" binding:"required"`
}

// HTTPHandler เก็บ dependency ที่จำเป็นสำหรับ handler ซึ่งก็คือ BranchService
type HTTPHandler struct {
	branchService        ports.BranchService
	interestService      ports.InterestService
	productService       ports.ProductService
	productOptionService ports.ProductOptionService
}

// NewHTTPHandler คือ factory function สำหรับสร้าง HTTPHandler
func NewHTTPHandler(branchSvc ports.BranchService, interestSvc ports.InterestService, productSvc ports.ProductService, productOptionSvc ports.ProductOptionService) *HTTPHandler {
	return &HTTPHandler{
		branchService:        branchSvc,
		interestService:      interestSvc,
		productService:       productSvc,
		productOptionService: productOptionSvc,
	}
}

// CreateBranch คือ handler function สำหรับ endpoint สร้างสาขา
func (h *HTTPHandler) CreateBranch(c *gin.Context) {
	var req CreateBranchRequest

	// 1. แปลง JSON body ที่ส่งมาให้เป็น struct และตรวจสอบความถูกต้อง
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. เรียกใช้ service เพื่อทำงานตาม business logic
	branch, err := h.branchService.CreateBranchWithProducts(c.Request.Context(), req.Name, req.ProductIDs)
	if err != nil {
		// Log error ฉบับเต็มไว้สำหรับนักพัฒนา
		log.Printf("Error creating branch: %v", err)
		// ส่งข้อความที่เป็นกลางกลับไปให้ client
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	// 3. ส่งผลลัพธ์กลับไปเป็น JSON พร้อม status 201 Created
	c.JSON(http.StatusCreated, gin.H{"data": branch})
}

// GetBranch คือ handler สำหรับดึงข้อมูลสาขา
func (h *HTTPHandler) GetBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	branch, err := h.branchService.GetBranch(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error getting branch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": branch})
}

// UpdateBranch คือ handler สำหรับอัปเดตข้อมูลสาขา
func (h *HTTPHandler) UpdateBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var req CreateBranchRequest // ใช้ struct เดียวกับ Create
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branch, err := h.branchService.UpdateBranchWithProducts(c.Request.Context(), id, req.Name, req.ProductIDs)
	if err != nil {
		log.Printf("Error updating branch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": branch})
}

// DeleteBranch คือ handler สำหรับลบข้อมูลสาขา
func (h *HTTPHandler) DeleteBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	err = h.branchService.DeleteBranch(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error deleting branch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
}

// --- Interest Handlers ---

func (h *HTTPHandler) UpdateInterest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interest ID"})
		return
	}

	var req UpdateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.interestService.UpdateInterest(c.Request.Context(), id, req.Name); err != nil {
		log.Printf("Error updating interest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interest updated successfully"})
}

func (h *HTTPHandler) DeleteInterest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interest ID"})
		return
	}

	if err := h.interestService.DeleteInterest(c.Request.Context(), id); err != nil {
		log.Printf("Error deleting interest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interest deleted successfully"})
}

// --- Product Handlers ---

func (h *HTTPHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req UpdateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.UpdateProduct(c.Request.Context(), id, req.Name); err != nil {
		log.Printf("Error updating product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (h *HTTPHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		log.Printf("Error deleting product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// --- Product Option Handlers ---

func (h *HTTPHandler) UpdateProductOption(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product option ID"})
		return
	}

	var req UpdateProductOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productOptionService.UpdateProductOption(c.Request.Context(), id, req.NormalPrice, req.TagthaiPrice); err != nil {
		log.Printf("Error updating product option: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product option updated successfully"})
}

func (h *HTTPHandler) DeleteProductOption(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product option ID"})
		return
	}

	if err := h.productOptionService.DeleteProductOption(c.Request.Context(), id); err != nil {
		log.Printf("Error deleting product option: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product option deleted successfully"})
}
