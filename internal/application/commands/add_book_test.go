package commands

import (
	"context"
	"testing"

	"library-system/internal/domain/catalog"
)

// MockBookRepository is a test double for catalog.BookRepository
type MockBookRepository struct {
	books     map[string]*catalog.Book
	addError  error
	getError  error
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books: make(map[string]*catalog.Book),
	}
}

func (m *MockBookRepository) Add(ctx context.Context, book *catalog.Book) error {
	if m.addError != nil {
		return m.addError
	}
	m.books[book.ID().String()] = book
	return nil
}

func (m *MockBookRepository) GetByID(ctx context.Context, id catalog.BookID) (*catalog.Book, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	book, exists := m.books[id.String()]
	if !exists {
		return nil, nil
	}
	return book, nil
}

func (m *MockBookRepository) List(ctx context.Context, limit, offset int) ([]*catalog.Book, error) {
	books := make([]*catalog.Book, 0, len(m.books))
	for _, book := range m.books {
		books = append(books, book)
	}
	// Simple pagination for tests
	if offset >= len(books) {
		return []*catalog.Book{}, nil
	}
	end := offset + limit
	if end > len(books) {
		end = len(books)
	}
	return books[offset:end], nil
}

func (m *MockBookRepository) Count(ctx context.Context) (int, error) {
	return len(m.books), nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *catalog.Book) error {
	m.books[book.ID().String()] = book
	return nil
}

func (m *MockBookRepository) Remove(ctx context.Context, id catalog.BookID) error {
	delete(m.books, id.String())
	return nil
}

// --- Tests ---

func TestAddBookHandler_Success(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewAddBookHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, AddBookCommand{
		Title:  "Clean Code",
		Author: "Robert Martin",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Title != "Clean Code" {
		t.Errorf("expected 'Clean Code', got %s", result.Title)
	}
	if result.Author != "Robert Martin" {
		t.Errorf("expected 'Robert Martin', got %s", result.Author)
	}
	if result.IsBorrowed {
		t.Error("new book should not be borrowed")
	}
	if len(repo.books) != 1 {
		t.Errorf("expected 1 book in repo, got %d", len(repo.books))
	}
}

func TestAddBookHandler_EmptyTitle(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewAddBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, AddBookCommand{
		Title:  "",
		Author: "Robert Martin",
	})

	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestAddBookHandler_EmptyAuthor(t *testing.T) {
	repo := NewMockBookRepository()
	handler := NewAddBookHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, AddBookCommand{
		Title:  "Clean Code",
		Author: "",
	})

	if err == nil {
		t.Error("expected error for empty author")
	}
}
