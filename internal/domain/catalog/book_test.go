package catalog

import (
	"testing"
	"time"
)

// --- Value Object Tests ---

func TestParseBookID_ValidUUID(t *testing.T) {
	id, err := ParseBookID("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if id.String() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected ID to match input")
	}
}

func TestParseBookID_EmptyString(t *testing.T) {
	_, err := ParseBookID("")
	if err != ErrBookIDEmpty {
		t.Errorf("expected ErrBookIDEmpty, got %v", err)
	}
}

func TestParseBookID_InvalidFormat(t *testing.T) {
	_, err := ParseBookID("not-a-uuid")
	if err != ErrBookIDInvalidFormat {
		t.Errorf("expected ErrBookIDInvalidFormat, got %v", err)
	}
}

func TestGenerateBookID(t *testing.T) {
	id1 := GenerateBookID()
	id2 := GenerateBookID()
	
	if id1.String() == "" {
		t.Error("expected non-empty ID")
	}
	if id1.String() == id2.String() {
		t.Error("expected unique IDs")
	}
}

func TestNewTitle_Valid(t *testing.T) {
	title, err := NewTitle("Clean Code")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if title.String() != "Clean Code" {
		t.Errorf("expected 'Clean Code', got %s", title.String())
	}
}

func TestNewTitle_Empty(t *testing.T) {
	_, err := NewTitle("")
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestNewTitle_TooLong(t *testing.T) {
	longTitle := make([]byte, 101)
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	_, err := NewTitle(string(longTitle))
	if err == nil {
		t.Error("expected error for title > 100 chars")
	}
}

func TestNewAuthor_Valid(t *testing.T) {
	author, err := NewAuthor("Robert Martin")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if author.String() != "Robert Martin" {
		t.Errorf("expected 'Robert Martin', got %s", author.String())
	}
}

func TestNewAuthor_Empty(t *testing.T) {
	_, err := NewAuthor("")
	if err == nil {
		t.Error("expected error for empty author")
	}
}

// --- Entity Tests ---

func TestNewBook(t *testing.T) {
	id := GenerateBookID()
	title, _ := NewTitle("Clean Code")
	author, _ := NewAuthor("Robert Martin")
	
	book := NewBook(id, title, author)
	
	if book.ID() != id {
		t.Error("ID mismatch")
	}
	if book.Title() != title {
		t.Error("Title mismatch")
	}
	if book.Author() != author {
		t.Error("Author mismatch")
	}
	if book.IsBorrowed() {
		t.Error("new book should not be borrowed")
	}
}

func TestBook_Borrow_Success(t *testing.T) {
	book := createTestBook()
	borrowedAt := time.Now()
	
	err := book.Borrow("john@example.com", borrowedAt)
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !book.IsBorrowed() {
		t.Error("book should be borrowed")
	}
	if len(book.GetEvents()) != 1 {
		t.Errorf("expected 1 event, got %d", len(book.GetEvents()))
	}
}

func TestBook_Borrow_AlreadyBorrowed(t *testing.T) {
	book := createTestBook()
	book.Borrow("john@example.com", time.Now())
	
	err := book.Borrow("jane@example.com", time.Now())
	
	if err != ErrBookAlreadyBorrowed {
		t.Errorf("expected ErrBookAlreadyBorrowed, got %v", err)
	}
}

func TestBook_Borrow_EmptyEmail(t *testing.T) {
	book := createTestBook()
	
	err := book.Borrow("", time.Now())
	
	if err != ErrBorrowerEmailRequired {
		t.Errorf("expected ErrBorrowerEmailRequired, got %v", err)
	}
}

func TestBook_Return_Success(t *testing.T) {
	book := createTestBook()
	book.Borrow("john@example.com", time.Now())
	book.ClearEvents() // Clear borrow event
	
	err := book.Return()
	
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if book.IsBorrowed() {
		t.Error("book should not be borrowed")
	}
	if len(book.GetEvents()) != 1 {
		t.Errorf("expected 1 event, got %d", len(book.GetEvents()))
	}
}

func TestBook_Return_NotBorrowed(t *testing.T) {
	book := createTestBook()
	
	err := book.Return()
	
	if err != ErrBookNotBorrowed {
		t.Errorf("expected ErrBookNotBorrowed, got %v", err)
	}
}

// --- Test Helpers ---

func createTestBook() *Book {
	id := GenerateBookID()
	title, _ := NewTitle("Test Book")
	author, _ := NewAuthor("Test Author")
	return NewBook(id, title, author)
}
