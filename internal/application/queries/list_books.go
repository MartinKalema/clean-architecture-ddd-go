package queries

import (
	"context"

	"library-system/internal/domain/catalog"
)

const (
	DefaultLimit = 50
	MaxLimit     = 100
)

// ListBooksQuery represents a request to list books with pagination
type ListBooksQuery struct {
	Limit  int
	Offset int
}

// BookSummary is a simplified view of a book for listings
type BookSummary struct {
	ID         string
	Title      string
	Author     string
	IsBorrowed bool
}

// ListBooksResult is returned after fetching books
type ListBooksResult struct {
	Books  []BookSummary `json:"books"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

// ListBooksHandler handles the ListBooksQuery
type ListBooksHandler struct {
	repo catalog.BookRepository
}

// NewListBooksHandler creates a new handler
func NewListBooksHandler(repo catalog.BookRepository) *ListBooksHandler {
	return &ListBooksHandler{repo: repo}
}

// Handle executes the query
func (h *ListBooksHandler) Handle(ctx context.Context, query ListBooksQuery) (ListBooksResult, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	books, err := h.repo.List(ctx, limit, offset)
	if err != nil {
		return ListBooksResult{}, err
	}

	total, err := h.repo.Count(ctx)
	if err != nil {
		return ListBooksResult{}, err
	}

	summaries := make([]BookSummary, len(books))
	for i, book := range books {
		summaries[i] = BookSummary{
			ID:         book.ID().String(),
			Title:      book.Title().String(),
			Author:     book.Author().String(),
			IsBorrowed: book.IsBorrowed(),
		}
	}

	return ListBooksResult{
		Books:  summaries,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
