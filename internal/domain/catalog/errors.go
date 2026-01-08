package catalog

import "errors"

var (
	ErrBookNotFound          = errors.New("book not found")
	ErrBookAlreadyBorrowed   = errors.New("book is already borrowed")
	ErrBookNotBorrowed       = errors.New("book is not borrowed")
	ErrBorrowerEmailRequired = errors.New("borrower email is required")
	ErrBookIDEmpty           = errors.New("book ID cannot be empty")
	ErrBookIDInvalidFormat   = errors.New("book ID must be a valid UUID")
)
