package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"library-system/internal/application/commands"
	"library-system/internal/application/queries"
	"library-system/internal/delivery/http/models"
)

// BookHandler handles book HTTP requests
type BookHandler struct {
	addBook    *commands.AddBookHandler
	borrowBook *commands.BorrowBookHandler
	returnBook *commands.ReturnBookHandler
	getBook    *queries.GetBookHandler
	listBooks  *queries.ListBooksHandler
}

// NewBookHandler creates a new handler
func NewBookHandler(
	addBook *commands.AddBookHandler,
	borrowBook *commands.BorrowBookHandler,
	returnBook *commands.ReturnBookHandler,
	getBook *queries.GetBookHandler,
	listBooks *queries.ListBooksHandler,
) *BookHandler {
	return &BookHandler{
		addBook:    addBook,
		borrowBook: borrowBook,
		returnBook: returnBook,
		getBook:    getBook,
		listBooks:  listBooks,
	}
}

// AddBook handles POST /books
func (h *BookHandler) AddBook(c *gin.Context) {
	var req models.AddBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.addBook.Handle(c.Request.Context(), commands.AddBookCommand{
		Title:  req.Title,
		Author: req.Author,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetBook handles GET /books/:id
func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")

	result, err := h.getBook.Handle(c.Request.Context(), queries.GetBookQuery{
		BookID: id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListBooks handles GET /books
func (h *BookHandler) ListBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result, err := h.listBooks.Handle(c.Request.Context(), queries.ListBooksQuery{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// BorrowBook handles POST /books/:id/borrow
func (h *BookHandler) BorrowBook(c *gin.Context) {
	id := c.Param("id")

	var req models.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.borrowBook.Handle(c.Request.Context(), commands.BorrowBookCommand{
		BookID:        id,
		BorrowerEmail: req.BorrowerEmail,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReturnBook handles POST /books/:id/return
func (h *BookHandler) ReturnBook(c *gin.Context) {
	id := c.Param("id")

	result, err := h.returnBook.Handle(c.Request.Context(), commands.ReturnBookCommand{
		BookID: id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
