package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"secure-api/internal/db"
	"secure-api/internal/handlers"
	"secure-api/internal/middleware"

	_ "secure-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// ---------------- LOAD ENV ----------------
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(".env file not found or failed to load!")
	}

	// ---------------- DB ----------------
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is empty!")
	}

	db.Connect(dbURL)

	// ---------------- GIN ----------------
	r := gin.Default()

	// ---------------- PUBLIC ROUTES ----------------
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// ---------------- PROTECTED ROUTES ----------------
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/account/coins", handlers.GetCoins)
	}

	// ---------------- SWAGGER ----------------
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ---------------- DEBUG ROUTES ----------------
	log.Println("The Server running on http://localhost:8000")

	for _, route := range r.Routes() {
		log.Printf("➡ %s %s\n", route.Method, route.Path)
	}

	// ---------------- RUN SERVER ----------------
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Error: Failed to start server:", err)
	}
}
