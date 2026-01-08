-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    is_borrowed BOOLEAN NOT NULL DEFAULT FALSE,
    borrowed_at TIMESTAMP,
    return_due_date TIMESTAMP,
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for listing available books
CREATE INDEX idx_books_is_borrowed ON books(is_borrowed);

-- Index for due date queries
CREATE INDEX idx_books_return_due_date ON books(return_due_date) WHERE return_due_date IS NOT NULL;
