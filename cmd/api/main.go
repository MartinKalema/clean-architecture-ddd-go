package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
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

	// Database configuration
	primaryURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/library?sslmode=disable")
	replicaURLs := getReplicaURLs()

	// Create database cluster (external client)
	ctx := context.Background()
	cluster, err := external.NewDBCluster(ctx, primaryURL, replicaURLs)
	if err != nil {
		slog.Error("failed to connect to database cluster", "error", err)
		os.Exit(1)
	}
	defer cluster.Close()

	slog.Info("connected to database cluster",
		"primary", maskPassword(primaryURL),
		"replicas", len(replicaURLs),
	)

	// Create repository - inject pools directly (not the cluster)
	bookRepo := catalogRepo.NewBookRepository(
		cluster.Primary(), // writer pool
		cluster.Replica(), // reader pool (round-robin if multiple)
	)

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
	port := getEnv("PORT", "8080")

	slog.Info("starting server", "port", port)
	if err := router.Run(":" + port); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getReplicaURLs() []string {
	// Check for comma-separated replica URLs
	if urls := os.Getenv("DATABASE_REPLICA_URLS"); urls != "" {
		return strings.Split(urls, ",")
	}

	// Default local replicas
	replica1 := os.Getenv("DATABASE_REPLICA_1_URL")
	replica2 := os.Getenv("DATABASE_REPLICA_2_URL")

	if replica1 == "" && replica2 == "" {
		// Default to local Docker replicas
		return []string{
			"postgres://postgres:postgres@localhost:5433/library?sslmode=disable",
			"postgres://postgres:postgres@localhost:5434/library?sslmode=disable",
		}
	}

	var urls []string
	if replica1 != "" {
		urls = append(urls, replica1)
	}
	if replica2 != "" {
		urls = append(urls, replica2)
	}
	return urls
}

func maskPassword(url string) string {
	if idx := strings.Index(url, "://"); idx != -1 {
		rest := url[idx+3:]
		if atIdx := strings.Index(rest, "@"); atIdx != -1 {
			if colonIdx := strings.Index(rest[:atIdx], ":"); colonIdx != -1 {
				return url[:idx+3+colonIdx+1] + "****" + url[idx+3+atIdx:]
			}
		}
	}
	return url
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
