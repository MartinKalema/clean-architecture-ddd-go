package queries

import (
	"context"
	"time"

	"library-system/internal/domain/catalog"
)

// GetBookQuery represents a request to get a book by ID
type GetBookQuery struct {
	BookID string
}

// GetBookResult is returned after fetching a book
type GetBookResult struct {
	ID            string
	Title         string
	Author        string
	IsBorrowed    bool
	BorrowedAt    *time.Time
	ReturnDueDate *time.Time
}

// GetBookHandler handles the GetBookQuery
type GetBookHandler struct {
	repo catalog.BookRepository
}

// NewGetBookHandler creates a new handler
func NewGetBookHandler(repo catalog.BookRepository) *GetBookHandler {
	return &GetBookHandler{repo: repo}
}

// Handle executes the query
func (h *GetBookHandler) Handle(ctx context.Context, query GetBookQuery) (GetBookResult, error) {
	bookID, err := catalog.ParseBookID(query.BookID)
	if err != nil {
		return GetBookResult{}, err
	}

	book, err := h.repo.GetByID(ctx, bookID)
	if err != nil {
		return GetBookResult{}, err
	}
	if book == nil {
		return GetBookResult{}, catalog.ErrBookNotFound
	}

	return GetBookResult{
		ID:         book.ID().String(),
		Title:      book.Title().String(),
		Author:     book.Author().String(),
		IsBorrowed: book.IsBorrowed(),
	}, nil
}
