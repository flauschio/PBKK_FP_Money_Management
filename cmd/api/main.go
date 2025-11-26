package main

import (
	"log"
	"os"

	"finance-manager/database"
	"finance-manager/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "finance_manager")
	port := getEnv("PORT", "8080")
	db, err := database.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	handler := handlers.NewHandler(db.DB)
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./templates/index.html")
	})
	api := router.Group("/api")
	{
		api.GET("/categories", handler.GetCategories)
		api.POST("/categories", handler.CreateCategory)
		api.PUT("/categories/:id", handler.UpdateCategory)
		api.DELETE("/categories/:id", handler.DeleteCategory)
		api.GET("/transactions", handler.GetTransactions)
		api.POST("/transactions", handler.CreateTransaction)
		api.PUT("/transactions/:id", handler.UpdateTransaction)
		api.DELETE("/transactions/:id", handler.DeleteTransaction)
		api.GET("/dashboard/stats", handler.GetDashboardStats)
		api.GET("/accounts", handler.GetAccounts)
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
