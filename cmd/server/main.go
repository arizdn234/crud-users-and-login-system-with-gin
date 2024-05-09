package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.User{})

	// seed
	ur := repository.NewUserRepository(db)
	if err := ur.Seed(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	PORT := "8080"
	s := server.NewServer(db, PORT)

	log.Printf("Server is running on port %v\n\n`http://localhost:%v`", PORT, PORT)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
