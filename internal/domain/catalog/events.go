package catalog

import "time"

// BookAdded is raised when a new book is added to the catalog
type BookAdded struct {
	BookID string
	Title  string
	Author string
}

func (e BookAdded) EventName() string {
	return "catalog.book_added"
}

// BookBorrowed is raised when a book is borrowed
type BookBorrowed struct {
	BookID        string
	Title         string
	BorrowedAt    time.Time
	ReturnDate    time.Time
	BorrowerEmail string
}

func (e BookBorrowed) EventName() string {
	return "catalog.book_borrowed"
}

// BookReturned is raised when a book is returned
type BookReturned struct {
	BookID string
}

func (e BookReturned) EventName() string {
	return "catalog.book_returned"
}
