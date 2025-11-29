package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:500;not null" json:"name"`
	Email        string    `gorm:"size:500;uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"size:500;not null" json:"-"`
	RefreshToken string    `gorm:"type:text" json:"-"`
	AccessToken  string    `gorm:"type:text" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Account struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BankName  string    `gorm:"size:255;not null" json:"bank_name"`
	Amount    float64   `gorm:"type:decimal(12,2);not null;default:0" json:"amount"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	CategoryID *uint     `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	AccountID  *uint     `gorm:"index" json:"account_id"`
	Account    *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type Budget struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Criteria   string    `gorm:"size:50;not null" json:"criteria"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type ScheduledTransaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Amount     float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Repetition string    `gorm:"size:50;not null" json:"repetition"`
	RepeatAt   time.Time `gorm:"not null" json:"repeat_at"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	CategoryID *uint     `gorm:"index" json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	AccountID  *uint     `gorm:"index" json:"account_id"`
	Account    *Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}
func NewDatabase(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting database instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "finance_manager")
	port := getEnv("PORT", "8080")
	db, err := NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	handler := NewHandler(db)
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./templates/index.html")
	})
	router.GET("/register", func(c *gin.Context) {
		c.File("./templates/register.html")
	})
	router.GET("/login", func(c *gin.Context) {
		c.File("./templates/login.html")
	})
	secret := getEnv("JWT_SECRET", "replace-with-secure-secret")

	api := router.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}

	protected := api.Group("", authMiddleware(secret))
	{
		protected.GET("/categories", handler.GetCategories)
		protected.POST("/categories", handler.CreateCategory)
		protected.PUT("/categories/:id", handler.UpdateCategory)
		protected.DELETE("/categories/:id", handler.DeleteCategory)
		protected.GET("/transactions", handler.GetTransactions)
		protected.POST("/transactions", handler.CreateTransaction)
		protected.PUT("/transactions/:id", handler.UpdateTransaction)
		protected.DELETE("/transactions/:id", handler.DeleteTransaction)
		protected.GET("/dashboard/stats", handler.GetDashboardStats)
		protected.GET("/accounts", handler.GetAccounts)
		protected.POST("/accounts", handler.CreateAccount)
		protected.PUT("/accounts/:id", handler.UpdateAccount)
		protected.DELETE("/accounts/:id", handler.DeleteAccount)
		protected.GET("/budgets", handler.GetBudgets)
		protected.POST("/budgets", handler.CreateBudget)
		protected.PUT("/budgets/:id", handler.UpdateBudget)
		protected.DELETE("/budgets/:id", handler.DeleteBudget)
		protected.POST("/budgets/check", handler.CheckBudgetExceeded)
		protected.GET("/scheduled", handler.GetScheduledTransactions)
		protected.POST("/scheduled", handler.CreateScheduledTransaction)
		protected.PUT("/scheduled/:id", handler.UpdateScheduledTransaction)
		protected.DELETE("/scheduled/:id", handler.DeleteScheduledTransaction)
		protected.POST("/scheduled/process", handler.ProcessScheduledTransactions)
	}
	log.Printf("Server starting on port %s...", port)
	log.Printf("Open http://localhost:%s in your browser", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func authMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Next()
	}
}
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
	var existing User
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
	user := User{
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
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	secret := getEnv("JWT_SECRET", "replace-with-secure-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	c.JSON(http.StatusOK, gin.H{"access_token": tokenString, "user": user})
}
func (h *Handler) GetCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	var categories []Category
	h.DB.Where("user_id = ?", userID).Order("name").Find(&categories)
	c.JSON(http.StatusOK, categories)
}
func (h *Handler) CreateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	c.ShouldBindJSON(&req)
	category := Category{Name: req.Name, UserID: userID}
	h.DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}
func (h *Handler) UpdateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var category Category
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	c.ShouldBindJSON(&req)
	category.Name = req.Name
	h.DB.Save(&category)
	c.JSON(http.StatusOK, category)
}
func (h *Handler) DeleteCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	if result := h.DB.Where("user_id = ?", userID).Delete(&Category{}, c.Param("id")); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	var transactions []Transaction
	h.DB.Where("user_id = ?", userID).Preload("Category").Preload("Account").Order("created_at DESC").Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}
func (h *Handler) CreateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req Transaction
	c.ShouldBindJSON(&req)
	req.UserID = userID
	h.DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&req)
		if req.AccountID != nil {
			tx.Model(&Account{}).Where("id = ? AND user_id = ?", *req.AccountID, userID).
				UpdateColumn("amount", gorm.Expr("amount + ?", req.Amount))
		}
		return nil
	})
	c.JSON(http.StatusCreated, req)
}
func (h *Handler) UpdateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var transaction Transaction
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	oldAmount, oldAccountID := transaction.Amount, transaction.AccountID
	c.ShouldBindJSON(&transaction)
	transaction.UserID = userID
	h.DB.Transaction(func(tx *gorm.DB) error {
		if oldAccountID != nil {
			tx.Model(&Account{}).Where("id = ? AND user_id = ?", *oldAccountID, userID).
				UpdateColumn("amount", gorm.Expr("amount - ?", oldAmount))
		}
		if transaction.AccountID != nil {
			tx.Model(&Account{}).Where("id = ? AND user_id = ?", *transaction.AccountID, userID).
				UpdateColumn("amount", gorm.Expr("amount + ?", transaction.Amount))
		}
		tx.Save(&transaction)
		return nil
	})
	c.JSON(http.StatusOK, transaction)
}
func (h *Handler) DeleteTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var transaction Transaction
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	h.DB.Transaction(func(tx *gorm.DB) error {
		if transaction.AccountID != nil {
			tx.Model(&Account{}).Where("id = ? AND user_id = ?", *transaction.AccountID, userID).
				UpdateColumn("amount", gorm.Expr("amount - ?", transaction.Amount))
		}
		tx.Delete(&transaction)
		return nil
	})
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")
	var accounts []Account
	h.DB.Where("user_id = ?", userID).Find(&accounts)
	c.JSON(http.StatusOK, accounts)
}
func (h *Handler) CreateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var account Account
	c.ShouldBindJSON(&account)
	account.UserID = userID
	h.DB.Create(&account)
	c.JSON(http.StatusCreated, account)
}
func (h *Handler) UpdateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var account Account
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.ShouldBindJSON(&account)
	account.UserID = userID
	h.DB.Save(&account)
	c.JSON(http.StatusOK, account)
}
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	if result := h.DB.Where("user_id = ?", userID).Delete(&Account{}, c.Param("id")); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	var income, expense float64
	var count int64
	h.DB.Model(&Transaction{}).Where("user_id = ? AND amount > 0 AND created_at BETWEEN ? AND ?", userID, start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&income)
	h.DB.Model(&Transaction{}).Where("user_id = ? AND amount < 0 AND created_at BETWEEN ? AND ?", userID, start, end).
		Select("COALESCE(SUM(ABS(amount)), 0)").Scan(&expense)
	h.DB.Model(&Transaction{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).Count(&count)
	var accounts []Account
	h.DB.Where("user_id = ?", userID).Find(&accounts)
	var total float64
	for _, a := range accounts {
		total += a.Amount
	}
	c.JSON(http.StatusOK, gin.H{
		"total_income":          income,
		"total_expenses":        expense,
		"transaction_count":     count,
		"balance":               income - expense,
		"accounts":              accounts,
		"total_account_balance": total,
	})
}
func (h *Handler) GetBudgets(c *gin.Context) {
	userID := c.GetUint("user_id")
	var budgets []Budget
	h.DB.Where("user_id = ?", userID).Preload("Category").Find(&budgets)
	type Response struct {
		Budget
		Spent      float64 `json:"spent"`
		Remaining  float64 `json:"remaining"`
		Percentage float64 `json:"percentage"`
	}
	var result []Response
	for _, b := range budgets {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		if b.Criteria == "annual" {
			start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		}
		var spent float64
		h.DB.Model(&Transaction{}).
			Where("user_id = ? AND category_id = ? AND amount < 0 AND created_at >= ?", userID, b.CategoryID, start).
			Select("COALESCE(SUM(ABS(amount)), 0)").Scan(&spent)
		remaining := b.Amount - spent
		percentage := 0.0
		if b.Amount > 0 {
			percentage = (remaining / b.Amount) * 100
		}
		result = append(result, Response{b, spent, remaining, percentage})
	}
	c.JSON(http.StatusOK, result)
}
func (h *Handler) CreateBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	var budget Budget
	c.ShouldBindJSON(&budget)
	budget.UserID = userID
	h.DB.Create(&budget)
	c.JSON(http.StatusCreated, budget)
}
func (h *Handler) UpdateBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	var budget Budget
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}
	c.ShouldBindJSON(&budget)
	budget.UserID = userID
	h.DB.Save(&budget)
	c.JSON(http.StatusOK, budget)
}
func (h *Handler) DeleteBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	if result := h.DB.Where("user_id = ?", userID).Delete(&Budget{}, c.Param("id")); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) CheckBudgetExceeded(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		CategoryID *uint   `json:"category_id"`
		Amount     float64 `json:"amount"`
	}
	c.ShouldBindJSON(&req)
	if req.CategoryID == nil || req.Amount >= 0 {
		c.JSON(http.StatusOK, gin.H{"exceeded": false})
		return
	}
	var budget Budget
	if err := h.DB.Where("category_id = ? AND user_id = ?", *req.CategoryID, userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"exceeded": false})
		return
	}
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if budget.Criteria == "annual" {
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	}
	var spent float64
	h.DB.Model(&Transaction{}).
		Where("user_id = ? AND category_id = ? AND amount < 0 AND created_at >= ?", userID, *req.CategoryID, start).
		Select("COALESCE(SUM(ABS(amount)), 0)").Scan(&spent)
	newTotal := spent + (-req.Amount)
	c.JSON(http.StatusOK, gin.H{
		"exceeded":  newTotal > budget.Amount,
		"budget":    budget.Amount,
		"spent":     spent,
		"new_total": newTotal,
	})
}
func (h *Handler) GetScheduledTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	var scheduled []ScheduledTransaction
	h.DB.Where("user_id = ?", userID).Preload("Category").Preload("Account").Order("repeat_at ASC").Find(&scheduled)
	c.JSON(http.StatusOK, scheduled)
}
func (h *Handler) CreateScheduledTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var st ScheduledTransaction
	c.ShouldBindJSON(&st)
	st.UserID = userID
	h.DB.Create(&st)
	c.JSON(http.StatusCreated, st)
}
func (h *Handler) UpdateScheduledTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var st ScheduledTransaction
	if err := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&st).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.ShouldBindJSON(&st)
	st.UserID = userID
	h.DB.Save(&st)
	c.JSON(http.StatusOK, st)
}
func (h *Handler) DeleteScheduledTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	if result := h.DB.Where("user_id = ?", userID).Delete(&ScheduledTransaction{}, c.Param("id")); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) ProcessScheduledTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	var scheduled []ScheduledTransaction
	h.DB.Where("user_id = ? AND repeat_at <= ?", userID, time.Now()).Find(&scheduled)
	processed := 0
	for _, st := range scheduled {
		h.DB.Transaction(func(tx *gorm.DB) error {
			tx.Create(&Transaction{
				Name:       st.Name,
				Amount:     st.Amount,
				UserID:     userID,
				CategoryID: st.CategoryID,
				AccountID:  st.AccountID,
			})
			if st.AccountID != nil {
				tx.Model(&Account{}).Where("id = ? AND user_id = ?", *st.AccountID, userID).
					UpdateColumn("amount", gorm.Expr("amount + ?", st.Amount))
			}
			var nextRepeat time.Time
			switch st.Repetition {
			case "monthly":
				nextRepeat = st.RepeatAt.AddDate(0, 1, 0)
			case "3 months":
				nextRepeat = st.RepeatAt.AddDate(0, 3, 0)
			case "6 months":
				nextRepeat = st.RepeatAt.AddDate(0, 6, 0)
			case "annually":
				nextRepeat = st.RepeatAt.AddDate(1, 0, 0)
			}
			tx.Model(&st).Update("repeat_at", nextRepeat)
			processed++
			return nil
		})
	}
	c.JSON(http.StatusOK, gin.H{"processed": processed})
}
