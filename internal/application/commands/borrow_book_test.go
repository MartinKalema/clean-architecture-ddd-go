package commands

import (
	"context"
	"testing"
	"time"

	"library-system/internal/domain/catalog"
)

func TestBorrowBookHandler_Success(t *testing.T) {
	repo := NewMockBookRepository()
	
	// Add a book first
	id := catalog.GenerateBookID()
	title, _ := catalog.NewTitle("Clean Code")
	author, _ := catalog.NewAuthor("Robert Martin")
	book := catalog.NewBook(id, title, author)
	repo.Add(context.Background(), book)

	handler := NewBorrowBookHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, BorrowBookCommand{
		BookID:        id.String(),
		BorrowerEmail: "john@example.com",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Title != "Clean Code" {
		t.Errorf("expected 'Clean Code', got %s", result.Title)
	}
	if result.BorrowedAt.IsZero() {
		t.Error("expected BorrowedAt to be set")
	}
	if result.ReturnDueDate.Before(time.Now()) {
		t.Error("expected ReturnDueDate to be in the future")
	}
}

func TestBorrowBookHandler_BookNotFound(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewBorrowBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, BorrowBookCommand{
		BookID:        "550e8400-e29b-41d4-a716-446655440000",
		BorrowerEmail: "john@example.com",
	})

	if err != catalog.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound, got %v", err)
	}
}

func TestBorrowBookHandler_InvalidBookID(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewBorrowBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, BorrowBookCommand{
		BookID:        "invalid-id",
		BorrowerEmail: "john@example.com",
	})

	if err != catalog.ErrBookIDInvalidFormat {
		t.Errorf("expected ErrBookIDInvalidFormat, got %v", err)
	}
}

func TestBorrowBookHandler_AlreadyBorrowed(t *testing.T) {
	repo := NewMockBookRepository()
	
	// Add and borrow a book
	id := catalog.GenerateBookID()
	title, _ := catalog.NewTitle("Clean Code")
	author, _ := catalog.NewAuthor("Robert Martin")
	book := catalog.NewBook(id, title, author)
	book.Borrow("first@example.com", time.Now())
	repo.Add(context.Background(), book)

	handler := NewBorrowBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, BorrowBookCommand{
		BookID:        id.String(),
		BorrowerEmail: "second@example.com",
	})

	if err != catalog.ErrBookAlreadyBorrowed {
		t.Errorf("expected ErrBookAlreadyBorrowed, got %v", err)
	}
}
