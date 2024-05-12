package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/models"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/repository"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	dbName := os.Getenv("DB_NAME")
	dbURI := dbName

	db, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// migrate
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate the table: %v", err)
	}

	// seed
	if err := repository.NewUserRepository(db).Seed(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	s := server.RunServer(gin.Default(), db, port)

	log.Printf("Server is running on port %v\n\n`http://localhost:%v`", port, port)
	if err := s.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
