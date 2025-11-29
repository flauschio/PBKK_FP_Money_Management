package handlers

import (
	"net/http"
	"os"
	"time"

	"finance-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	// create transaction and update account balance atomically
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		if transaction.AccountID != nil {
			// apply the transaction amount to the selected account
			if err := tx.Model(&models.Account{}).
				Where("id = ?", *transaction.AccountID).
				UpdateColumn("amount", gorm.Expr("amount + ?", transaction.Amount)).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
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
	oldAmount := transaction.Amount
	var oldAccountID *uint
	if transaction.AccountID != nil {
		v := *transaction.AccountID
		oldAccountID = &v
	}
	newAmount := req.Amount
	var newAccountID *uint
	if req.AccountID != nil {
		v := *req.AccountID
		newAccountID = &v
	}

	// perform update and adjust account balances atomically
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// adjust balances
		if oldAccountID != nil && newAccountID != nil && *oldAccountID == *newAccountID {
			// same account: apply delta
			delta := newAmount - oldAmount
			if delta != 0 {
				if err := tx.Model(&models.Account{}).
					Where("id = ?", *oldAccountID).
					UpdateColumn("amount", gorm.Expr("amount + ?", delta)).Error; err != nil {
					return err
				}
			}
		} else {
			// different accounts: revert old and apply new
			if oldAccountID != nil {
				if err := tx.Model(&models.Account{}).
					Where("id = ?", *oldAccountID).
					UpdateColumn("amount", gorm.Expr("amount - ?", oldAmount)).Error; err != nil {
					return err
				}
			}
			if newAccountID != nil {
				if err := tx.Model(&models.Account{}).
					Where("id = ?", *newAccountID).
					UpdateColumn("amount", gorm.Expr("amount + ?", newAmount)).Error; err != nil {
					return err
				}
			}
		}

		// update transaction record
		transaction.Name = req.Name
		transaction.Amount = newAmount
		transaction.CategoryID = req.CategoryID
		transaction.AccountID = req.AccountID
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
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
	// load transaction first to know amount/account
	var transaction models.Transaction
	if err := h.DB.First(&transaction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if transaction.AccountID != nil {
			if err := tx.Model(&models.Account{}).
				Where("id = ?", *transaction.AccountID).
				UpdateColumn("amount", gorm.Expr("amount - ?", transaction.Amount)).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&models.Transaction{}, transaction.ID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (h *Handler) CreateAccount(c *gin.Context) {
	var req struct {
		BankName string  `json:"bank_name" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := models.Account{
		BankName: req.BankName,
		Amount:   req.Amount,
	}
	if err := h.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	result := h.DB.Delete(&models.Account{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Register handles user registration: validates input, hashes password, and stores the user.
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login verifies credentials and returns a signed JWT access token.
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "replace-with-secure-secret"
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": tokenString, "user": user})
}
