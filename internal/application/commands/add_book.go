package commands

import (
	"context"
	"library-system/internal/domain/catalog"
)

// AddBookCommand represents intent to add a book
type AddBookCommand struct {
	Title  string
	Author string
}

// AddBookResult is returned after adding a book
type AddBookResult struct {
	ID         string
	Title      string
	Author     string
	IsBorrowed bool
}

// AddBookHandler handles the AddBookCommand
type AddBookHandler struct {
	repo catalog.BookRepository
}

// NewAddBookHandler creates a new handler
func NewAddBookHandler(repo catalog.BookRepository) *AddBookHandler {
	return &AddBookHandler{repo: repo}
}

// Handle executes the command
func (h *AddBookHandler) Handle(ctx context.Context, cmd AddBookCommand) (AddBookResult, error) {
	// Create value object
	title, err := catalog.NewTitle(cmd.Title)
	if err != nil {
		return AddBookResult{}, err
	}

	author, err := catalog.NewAuthor(cmd.Author)
	if err != nil {
		return AddBookResult{}, err
	}

	// Create entity
	book := catalog.NewBook(catalog.GenerateBookID(), title, author)

	// Persist
	if err := h.repo.Add(ctx, book); err != nil {
		return AddBookResult{}, err
	}

	// Return result
	return AddBookResult{
		ID:         book.ID().String(),
		Title:      book.Title().String(),
		Author:     book.Author().String(),
		IsBorrowed: book.IsBorrowed(),
	}, nil

}
