// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: author.sql

package example

import (
	"context"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO Author (
  author_id, name, biography
) VALUES (
  $1, $2, $3
)
RETURNING author_id, name, biography
`

type CreateAuthorParams struct {
	AuthorID  int32
	Name      string
	Biography []byte
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, arg.AuthorID, arg.Name, arg.Biography)
	var i Author
	err := row.Scan(&i.AuthorID, &i.Name, &i.Biography)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM Author
WHERE author_id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, authorID int32) error {
	_, err := q.db.Exec(ctx, deleteAuthor, authorID)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT author_id, name, biography FROM Author
WHERE author_id = $1 LIMIT 1
`

// Code generated by protoc-gen-sqlc. DO NOT EDIT.
// source:
//
//	examples/library/v1/author.proto
func (q *Queries) GetAuthor(ctx context.Context, authorID int32) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, authorID)
	var i Author
	err := row.Scan(&i.AuthorID, &i.Name, &i.Biography)
	return i, err
}

const listAuthor = `-- name: ListAuthor :many
SELECT author_id, name, biography FROM Author
ORDER BY author_id
`

func (q *Queries) ListAuthor(ctx context.Context) ([]Author, error) {
	rows, err := q.db.Query(ctx, listAuthor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.AuthorID, &i.Name, &i.Biography); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAuthor = `-- name: UpdateAuthor :one
UPDATE Author SET
  name = $2, 
  biography = $3
WHERE author_id = $1
RETURNING author_id, name, biography
`

type UpdateAuthorParams struct {
	AuthorID  int32
	Name      string
	Biography []byte
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	row := q.db.QueryRow(ctx, updateAuthor, arg.AuthorID, arg.Name, arg.Biography)
	var i Author
	err := row.Scan(&i.AuthorID, &i.Name, &i.Biography)
	return i, err
}
