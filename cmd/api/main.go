package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"library-system/internal/application/commands"
	"library-system/internal/application/queries"
	"library-system/internal/delivery/http/handlers"
	"library-system/internal/delivery/http/routes"
	catalogRepo "library-system/internal/infrastructure/adapters/catalog"
	"library-system/internal/infrastructure/external"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/library?sslmode=disable"
	}

	// Create database connection
	ctx := context.Background()
	pool, err := external.NewPostgresPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create repository
	bookRepo := catalogRepo.NewBookRepository(pool)

	// Create command handlers
	addBookHandler := commands.NewAddBookHandler(bookRepo)
	borrowBookHandler := commands.NewBorrowBookHandler(bookRepo)
	returnBookHandler := commands.NewReturnBookHandler(bookRepo)

	// Create query handlers
	getBookHandler := queries.NewGetBookHandler(bookRepo)
	listBooksHandler := queries.NewListBooksHandler(bookRepo)

	// Create HTTP handler
	bookHandler := handlers.NewBookHandler(
		addBookHandler,
		borrowBookHandler,
		returnBookHandler,
		getBookHandler,
		listBooksHandler,
	)

	// Setup router
	router := gin.Default()
	routes.Setup(router, bookHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
