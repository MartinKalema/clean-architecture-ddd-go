package catalog

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"library-system/internal/domain/catalog"
)

// bookRow represents a book row in the database
type bookRow struct {
	ID            string
	Title         string
	Author        string
	IsBorrowed    bool
	BorrowedAt    *time.Time
	ReturnDueDate *time.Time
	Version       int
}

// BookRepository implements catalog.BookRepository
type BookRepository struct {
	pool *pgxpool.Pool
}

// NewBookRepository creates a new repository.
func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

// Add inserts a new book
func (r *BookRepository) Add(ctx context.Context, book *catalog.Book) error {
	_, err := r.pool.Exec(
		ctx, `INSERT INTO books (id, title, author, is_borrowed, borrowed_at, return_due_date, version)
  		VALUES ($1, $2, $3, $4, $5, $6, $7)`, book.ID().String(), book.Title().String(), book.Author().String(), book.IsBorrowed(), nil, nil, book.Version(),
	)
	return err
}

// GetByID fetches a book by ID
func (r *BookRepository) GetByID(ctx context.Context, id catalog.BookID) (*catalog.Book, error) {
	var row bookRow
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, author, is_borrowed, borrowed_at, return_due_date, version
		FROM books WHERE id = $1
	`, id.String()).Scan(
		&row.ID, &row.Title, &row.Author,
		&row.IsBorrowed, &row.BorrowedAt, &row.ReturnDueDate, &row.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rowToBook(row)
}

// List fetches books with pagination
func (r *BookRepository) List(ctx context.Context, limit, offset int) ([]*catalog.Book, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, author, is_borrowed, borrowed_at, return_due_date, version
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*catalog.Book
	for rows.Next() {
		var row bookRow
		if err := rows.Scan(
			&row.ID, &row.Title, &row.Author,
			&row.IsBorrowed, &row.BorrowedAt, &row.ReturnDueDate, &row.Version,
		); err != nil {
			return nil, err
		}
		book, err := rowToBook(row)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

// Count returns total number of books
func (r *BookRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM books`).Scan(&count)
	return count, err
}

// Update updates an existing book
func (r *BookRepository) Update(ctx context.Context, book *catalog.Book) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE books
		SET title = $2, author = $3, is_borrowed = $4, borrowed_at = $5, return_due_date = $6, version = version + 1
		WHERE id = $1
	`, book.ID().String(), book.Title().String(), book.Author().String(),
		book.IsBorrowed(), nil, nil)
	return err
}

// Remove deletes a book
func (r *BookRepository) Remove(ctx context.Context, id catalog.BookID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id.String())
	return err
}

// rowToBook converts a database row to a domain entity
func rowToBook(row bookRow) (*catalog.Book, error) {
	bookID, err := catalog.ParseBookID(row.ID)
	if err != nil {
		return nil, err
	}
	title, err := catalog.NewTitle(row.Title)
	if err != nil {
		return nil, err
	}
	author, err := catalog.NewAuthor(row.Author)
	if err != nil {
		return nil, err
	}

	return catalog.ReconstructBook(
		bookID,
		title,
		author,
		row.IsBorrowed,
		row.BorrowedAt,
		row.ReturnDueDate,
		row.Version,
	), nil
}
