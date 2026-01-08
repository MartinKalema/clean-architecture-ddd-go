package routes

import (
	"github.com/gin-gonic/gin"

	"library-system/internal/delivery/http/handlers"
)

// Setup configures all routes
func Setup(router *gin.Engine, bookHandler *handlers.BookHandler) {
	api := router.Group("/api/v1")
	{
		books := api.Group("/books")
		{
			books.POST("", bookHandler.AddBook)
			books.GET("", bookHandler.ListBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.POST("/:id/borrow", bookHandler.BorrowBook)
			books.POST("/:id/return", bookHandler.ReturnBook)
		}
	}
}
