package main

import (
	"log"
	"os"

	"example.com/todo-api/internal/database"
	"example.com/todo-api/internal/server"
	"example.com/todo-api/internal/todo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found; using process environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := database.Open(dsn)
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	repository := todo.NewGormRepository(db)
	service := todo.NewService(repository)
	handler := todo.NewHandler(service)
	router := server.NewRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
