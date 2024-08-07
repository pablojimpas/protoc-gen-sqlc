// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: book.sql

package example

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBook = `-- name: CreateBook :one
INSERT INTO Book (
  book_id, author_id, isbn, book_type, title, year, available_time, tags, published, price
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING book_id, author_id, isbn, book_type, title, year, available_time, tags, published, price
`

type CreateBookParams struct {
	BookID        int32
	AuthorID      int32
	Isbn          string
	BookType      Booktype
	Title         string
	Year          int32
	AvailableTime pgtype.Timestamptz
	Tags          []string
	Published     pgtype.Bool
	Price         pgtype.Float8
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRow(ctx, createBook,
		arg.BookID,
		arg.AuthorID,
		arg.Isbn,
		arg.BookType,
		arg.Title,
		arg.Year,
		arg.AvailableTime,
		arg.Tags,
		arg.Published,
		arg.Price,
	)
	var i Book
	err := row.Scan(
		&i.BookID,
		&i.AuthorID,
		&i.Isbn,
		&i.BookType,
		&i.Title,
		&i.Year,
		&i.AvailableTime,
		&i.Tags,
		&i.Published,
		&i.Price,
	)
	return i, err
}

const deleteBook = `-- name: DeleteBook :exec
DELETE FROM Book
WHERE book_id = $1
`

func (q *Queries) DeleteBook(ctx context.Context, bookID int32) error {
	_, err := q.db.Exec(ctx, deleteBook, bookID)
	return err
}

const getBook = `-- name: GetBook :one
SELECT book_id, author_id, isbn, book_type, title, year, available_time, tags, published, price FROM Book
WHERE book_id = $1 LIMIT 1
`

// Code generated by protoc-gen-sqlc. DO NOT EDIT.
// source:
//
//	examples/library/v1/book.proto
func (q *Queries) GetBook(ctx context.Context, bookID int32) (Book, error) {
	row := q.db.QueryRow(ctx, getBook, bookID)
	var i Book
	err := row.Scan(
		&i.BookID,
		&i.AuthorID,
		&i.Isbn,
		&i.BookType,
		&i.Title,
		&i.Year,
		&i.AvailableTime,
		&i.Tags,
		&i.Published,
		&i.Price,
	)
	return i, err
}

const listBook = `-- name: ListBook :many
SELECT book_id, author_id, isbn, book_type, title, year, available_time, tags, published, price FROM Book
ORDER BY book_id
`

func (q *Queries) ListBook(ctx context.Context) ([]Book, error) {
	rows, err := q.db.Query(ctx, listBook)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.BookID,
			&i.AuthorID,
			&i.Isbn,
			&i.BookType,
			&i.Title,
			&i.Year,
			&i.AvailableTime,
			&i.Tags,
			&i.Published,
			&i.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBook = `-- name: UpdateBook :one
UPDATE Book SET
  author_id = $2, 
  isbn = $3, 
  book_type = $4, 
  title = $5, 
  year = $6, 
  available_time = $7, 
  tags = $8, 
  published = $9, 
  price = $10
WHERE book_id = $1
RETURNING book_id, author_id, isbn, book_type, title, year, available_time, tags, published, price
`

type UpdateBookParams struct {
	BookID        int32
	AuthorID      int32
	Isbn          string
	BookType      Booktype
	Title         string
	Year          int32
	AvailableTime pgtype.Timestamptz
	Tags          []string
	Published     pgtype.Bool
	Price         pgtype.Float8
}

func (q *Queries) UpdateBook(ctx context.Context, arg UpdateBookParams) (Book, error) {
	row := q.db.QueryRow(ctx, updateBook,
		arg.BookID,
		arg.AuthorID,
		arg.Isbn,
		arg.BookType,
		arg.Title,
		arg.Year,
		arg.AvailableTime,
		arg.Tags,
		arg.Published,
		arg.Price,
	)
	var i Book
	err := row.Scan(
		&i.BookID,
		&i.AuthorID,
		&i.Isbn,
		&i.BookType,
		&i.Title,
		&i.Year,
		&i.AvailableTime,
		&i.Tags,
		&i.Published,
		&i.Price,
	)
	return i, err
}
