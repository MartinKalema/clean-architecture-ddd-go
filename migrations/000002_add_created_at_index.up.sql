-- Index for ORDER BY created_at DESC in list queries
CREATE INDEX idx_books_created_at ON books(created_at DESC);
