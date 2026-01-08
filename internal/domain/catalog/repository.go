package catalog

import "context"

// BookRepository defines persistence operations for books
type BookRepository interface {
	Add(ctx context.Context, book *Book) error
	GetByID(ctx context.Context, id BookID) (*Book, error)
	List(ctx context.Context, limit, offset int) ([]*Book, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, book *Book) error
	Remove(ctx context.Context, id BookID) error
}
