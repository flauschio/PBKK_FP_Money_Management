package handlers

import (
	"net/http"
	"time"

	"finance-manager/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}
func (h *Handler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.DB.Order("name").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
func (h *Handler) CreateCategory(c *gin.Context) {
	var req models.CategoryCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category := models.Category{
		Name: req.Name,
	}
	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}
func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req models.CategoryCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var category models.Category
	if err := h.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category.Name = req.Name
	if err := h.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}
func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	result := h.DB.Delete(&models.Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	if err := h.DB.Preload("Category").Preload("Account").Order("created_at DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var response []models.TransactionResponse
	for _, t := range transactions {
		tr := models.TransactionResponse{
			ID:         t.ID,
			Name:       t.Name,
			Amount:     t.Amount,
			CategoryID: t.CategoryID,
			AccountID:  t.AccountID,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
		if t.Category != nil {
			tr.CategoryName = t.Category.Name
		}
		response = append(response, tr)
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req models.TransactionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction := models.Transaction{
		Name:       req.Name,
		Amount:     req.Amount,
		CategoryID: req.CategoryID,
		AccountID:  req.AccountID,
	}
	if err := h.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction.CategoryID != nil {
		h.DB.Preload("Category").First(&transaction, transaction.ID)
	}
	response := models.TransactionResponse{
		ID:         transaction.ID,
		Name:       transaction.Name,
		Amount:     transaction.Amount,
		CategoryID: transaction.CategoryID,
		AccountID:  transaction.AccountID,
		CreatedAt:  transaction.CreatedAt,
		UpdatedAt:  transaction.UpdatedAt,
	}
	if transaction.Category != nil {
		response.CategoryName = transaction.Category.Name
	}
	c.JSON(http.StatusCreated, response)
}
func (h *Handler) UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	var req models.TransactionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var transaction models.Transaction
	if err := h.DB.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	transaction.Name = req.Name
	transaction.Amount = req.Amount
	transaction.CategoryID = req.CategoryID
	transaction.AccountID = req.AccountID
	if err := h.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction.CategoryID != nil {
		h.DB.Preload("Category").First(&transaction, transaction.ID)
	}
	response := models.TransactionResponse{
		ID:         transaction.ID,
		Name:       transaction.Name,
		Amount:     transaction.Amount,
		CategoryID: transaction.CategoryID,
		AccountID:  transaction.AccountID,
		CreatedAt:  transaction.CreatedAt,
		UpdatedAt:  transaction.UpdatedAt,
	}
	if transaction.Category != nil {
		response.CategoryName = transaction.Category.Name
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	result := h.DB.Delete(&models.Transaction{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetDashboardStats(c *gin.Context) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	var stats struct {
		TotalIncome  float64
		TotalExpense float64
		Count        int64
	}
	h.DB.Model(&models.Transaction{}).
		Where("amount > 0 AND created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0) as total_income").
		Scan(&stats.TotalIncome)
	h.DB.Model(&models.Transaction{}).
		Where("amount < 0 AND created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(ABS(amount)), 0) as total_expense").
		Scan(&stats.TotalExpense)
	h.DB.Model(&models.Transaction{}).
		Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Count(&stats.Count)
	var accounts []models.Account
	h.DB.Find(&accounts)
	var totalAccountBalance float64
	for _, acc := range accounts {
		totalAccountBalance += acc.Amount
	}
	response := gin.H{
		"total_income":          stats.TotalIncome,
		"total_expenses":        stats.TotalExpense,
		"transaction_count":     stats.Count,
		"balance":               stats.TotalIncome - stats.TotalExpense,
		"accounts":              accounts,
		"total_account_balance": totalAccountBalance,
	}
	c.JSON(http.StatusOK, response)
}
func (h *Handler) GetAccounts(c *gin.Context) {
	var accounts []models.Account
	if err := h.DB.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
