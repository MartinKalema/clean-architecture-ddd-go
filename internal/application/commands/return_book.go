package commands

import (
	"context"

	"library-system/internal/domain/catalog"
)

// ReturnBookCommand represents intent to return a book
type ReturnBookCommand struct {
	BookID string
}

// ReturnBookResult is returned after returning a book
type ReturnBookResult struct {
	BookID string
	Title  string
}

// ReturnBookHandler handles the ReturnBookCommand
type ReturnBookHandler struct {
	repo catalog.BookRepository
}

// NewReturnBookHandler creates a new handler
func NewReturnBookHandler(repo catalog.BookRepository) *ReturnBookHandler {
	return &ReturnBookHandler{repo: repo}
}

// Handle executes the command
func (h *ReturnBookHandler) Handle(ctx context.Context, cmd ReturnBookCommand) (ReturnBookResult, error) {
	bookID, err := catalog.ParseBookID(cmd.BookID)
	if err != nil {
		return ReturnBookResult{}, err
	}

	book, err := h.repo.GetByID(ctx, bookID)
	if err != nil {
		return ReturnBookResult{}, err
	}
	if book == nil {
		return ReturnBookResult{}, catalog.ErrBookNotFound
	}

	if err := book.Return(); err != nil {
		return ReturnBookResult{}, err
	}

	if err := h.repo.Update(ctx, book); err != nil {
		return ReturnBookResult{}, err
	}

	return ReturnBookResult{
		BookID: book.ID().String(),
		Title:  book.Title().String(),
	}, nil
}
