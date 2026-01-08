package commands

import (
	"context"
	"library-system/internal/domain/catalog"
	"time"
)

// BorrowBookCommand represents intent to borrow a book
type BorrowBookCommand struct {
	BookID        string
	BorrowerEmail string
}

// BorrowBookResult is returned after borrowing a book.
type BorrowBookResult struct {
	BookID        string
	Title         string
	BorrowedAt    time.Time
	ReturnDueDate time.Time
}

// BorrowBookHandler handles the BorrowBookCommand
type BorrowBookHandler struct {
	repo catalog.BookRepository
}

// NewBorrowBookHandler creates a new handler
func NewBorrowBookHandler(repo catalog.BookRepository) *BorrowBookHandler {
	return &BorrowBookHandler{repo: repo}
}

// Handle executes the command
func (h *BorrowBookHandler) Handle(ctx context.Context, cmd BorrowBookCommand) (BorrowBookResult, error) {
	// Parse BookID
	bookID, err := catalog.ParseBookID(cmd.BookID)
	if err != nil {
		return BorrowBookResult{}, err
	}

	// Get book from repository
	book, err := h.repo.GetByID(ctx, bookID)
	if err != nil {
		return BorrowBookResult{}, err
	}
	if book == nil {
		return BorrowBookResult{}, catalog.ErrBookNotFound
	}

	// Execute domain logic
	borrowedAt := time.Now()
	if err := book.Borrow(cmd.BorrowerEmail, borrowedAt); err != nil {
		return BorrowBookResult{}, err
	}

	// Persist changes
	if err := h.repo.Update(ctx, book); err != nil {
		return BorrowBookResult{}, err
	}

	return BorrowBookResult{
		BookID:        book.ID().String(),
		Title:         book.Title().String(),
		BorrowedAt:    borrowedAt,
		ReturnDueDate: borrowedAt.AddDate(0, 0, 14),
	}, nil
}
