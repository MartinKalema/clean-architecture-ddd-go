package commands

import (
	"context"
	"testing"
	"time"

	"library-system/internal/domain/catalog"
)

func TestReturnBookHandler_Success(t *testing.T) {
	repo := NewMockBookRepository()
	
	// Add and borrow a book
	id := catalog.GenerateBookID()
	title, _ := catalog.NewTitle("Clean Code")
	author, _ := catalog.NewAuthor("Robert Martin")
	book := catalog.NewBook(id, title, author)
	_ = book.Borrow("john@example.com", time.Now())
	_ = repo.Add(context.Background(), book)

	handler := NewReturnBookHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, ReturnBookCommand{
		BookID: id.String(),
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Title != "Clean Code" {
		t.Errorf("expected 'Clean Code', got %s", result.Title)
	}
	
	// Verify book is no longer borrowed
	updatedBook, _ := repo.GetByID(ctx, id)
	if updatedBook.IsBorrowed() {
		t.Error("book should not be borrowed after return")
	}
}

func TestReturnBookHandler_BookNotFound(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewReturnBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, ReturnBookCommand{
		BookID: "550e8400-e29b-41d4-a716-446655440000",
	})

	if err != catalog.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestReturnBookHandler_NotBorrowed(t *testing.T) {
	repo := NewMockBookRepository()
	
	// Add a book (not borrowed)
	id := catalog.GenerateBookID()
	title, _ := catalog.NewTitle("Clean Code")
	author, _ := catalog.NewAuthor("Robert Martin")
	book := catalog.NewBook(id, title, author)
	_ = repo.Add(context.Background(), book)

	handler := NewReturnBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, ReturnBookCommand{
		BookID: id.String(),
	})

	if err != catalog.ErrBookNotBorrowed {
		t.Errorf("expected ErrBookNotBorrowed, got %v", err)
	}
}
