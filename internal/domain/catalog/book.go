package catalog

import (
	"library-system/internal/domain/shared"
	"time"

	"github.com/google/uuid"
)

// --- Value Objects ---
type BookID struct {
	value string
}

// ParseBookID validates and creates a BookID from a string
func ParseBookID(value string) (BookID, error) {
	if value == "" {
		return BookID{}, ErrBookIDEmpty
	}
	if _, err := uuid.Parse(value); err != nil {
		return BookID{}, ErrBookIDInvalidFormat
	}
	return BookID{value: value}, nil
}

// GenerateBookID creates a new unique BookID
func GenerateBookID() BookID {
	return BookID{value: uuid.New().String()}
}

func (id BookID) String() string {
	return id.value
}

// Title represents a book's title
// It must be non empty and less than 100 characters
type Title struct {
	value string
}

func NewTitle(value string) (Title, error) {
	if value == "" {
		return Title{}, shared.ValidationError{
			Field:   "Title",
			Message: "Title cannot be empty",
		}
	}
	if len(value) > 100 {
		return Title{}, shared.ValidationError{
			Field:   "Title",
			Message: "Title cannot exceed 100 characters",
		}
	}
	return Title{value: value}, nil
}

func (t Title) String() string {
	return t.value
}

// Author represents a books author
type Author struct {
	value string
}

func NewAuthor(value string) (Author, error) {
	if value == "" {
		return Author{}, shared.ValidationError{
			Field:   "Author",
			Message: "Author cannot be empty",
		}
	}
	return Author{value: value}, nil
}

func (a Author) String() string {
	return a.value
}

// --- Entity ---

// Book is the aggregate root for the catalog context
type Book struct {
	id            BookID
	title         Title
	author        Author
	isBorrowed    bool
	borrowedAt    *time.Time
	returnDueDate *time.Time

	events  []shared.DomainEvent
	version int
}

// NewBook creates a new book
func NewBook(id BookID, title Title, author Author) *Book {
	return &Book{
		id:     id,
		title:  title,
		author: author,
	}
}

// ReconstructBook rebuilds a Book from persistence (used by repositories only)
func ReconstructBook(
	id BookID,
	title Title,
	author Author,
	isBorrowed bool,
	borrowedAt *time.Time,
	returnDueDate *time.Time,
	version int,
) *Book {
	return &Book{
		id:            id,
		title:         title,
		author:        author,
		isBorrowed:    isBorrowed,
		borrowedAt:    borrowedAt,
		returnDueDate: returnDueDate,
		version:       version,
	}
}

// Getters
func (b *Book) ID() BookID {
	return b.id
}
func (b *Book) Title() Title {
	return b.title
}
func (b *Book) Author() Author {
	return b.author
}
func (b *Book) IsBorrowed() bool {
	return b.isBorrowed
}
func (b *Book) Version() int {
	return b.version
}

// Borrow marks the book as borrowed.
func (b *Book) Borrow(borrowerEmail string, borrowedAt time.Time) error {
	if b.isBorrowed {
		return ErrBookAlreadyBorrowed
	}
	if borrowerEmail == "" {
		return ErrBorrowerEmailRequired
	}

	b.isBorrowed = true
	b.borrowedAt = &borrowedAt
	dueDate := borrowedAt.AddDate(0, 0, 14)
	b.returnDueDate = &dueDate

	b.events = append(b.events, BookBorrowed{
		BookID:        b.id.String(),
		Title:         b.title.String(),
		BorrowedAt:    borrowedAt,
		ReturnDate:    dueDate,
		BorrowerEmail: borrowerEmail,
	})

	return nil
}

// Return marks the book as returned
func (b *Book) Return() error {
	if !b.isBorrowed {
		return ErrBookNotBorrowed
	}

	b.isBorrowed = false
	b.borrowedAt = nil
	b.returnDueDate = nil

	b.events = append(b.events, BookReturned{
		BookID: b.id.String(),
	})

	return nil
}

// Events methods
func (b *Book) GetEvents() []shared.DomainEvent {
	return b.events
}

func (b *Book) ClearEvents() {
	b.events = nil
}
