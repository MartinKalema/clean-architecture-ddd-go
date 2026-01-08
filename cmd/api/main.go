package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"library-system/internal/application/commands"
	"library-system/internal/application/queries"
	"library-system/internal/delivery/http/handlers"
	"library-system/internal/delivery/http/routes"
	catalogRepo "library-system/internal/infrastructure/adapters/catalog"
	"library-system/internal/infrastructure/external"
)

func main() {
	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/library?sslmode=disable"
	}

	// Create database connection
	ctx := context.Background()
	pool, err := external.NewPostgresPool(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("connected to database")

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

	// Setup router with structured logging
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(structuredLogger())

	routes.Setup(router, bookHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	if err := router.Run(":" + port); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func structuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		slog.Info("request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
