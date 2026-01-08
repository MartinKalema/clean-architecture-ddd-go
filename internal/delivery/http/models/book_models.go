package models

// AddBookRequest is the request body for adding a book
type AddBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// BorrowBookRequest is the request body for borrowing a book
type BorrowBookRequest struct {
	BorrowerEmail string `json:"borrower_email" binding:"required,email"`
}
