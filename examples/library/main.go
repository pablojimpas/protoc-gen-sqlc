// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	example "github.com/pablojimpas/protoc-gen-sqlc/examples/library/internal/gen/sqlc"
)

//go:embed internal/gen/pb/sqlc/schema.sql
var schemaFS embed.FS

func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgresql://exampleuser:examplepassword@localhost/exampledb")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := example.New(conn)

	buf, err := fs.ReadFile(schemaFS, "internal/gen/pb/sqlc/schema.sql")
	if err != nil {
		return err
	}

	if _, err := conn.Exec(ctx, string(buf)); err != nil {
		return fmt.Errorf("failed to execute initial schema: %w", err)
	}

	// list all authors
	authors, err := queries.ListAuthor(ctx)
	if err != nil {
		return err
	}

	log.Println("Authors list:")
	log.Println(authors)

	bob, err := queries.CreateAuthor(ctx, example.CreateAuthorParams{
		Name:     "Uncle Bob",
		AuthorID: 5,
	})
	if err != nil {
		return err
	}

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, example.CreateAuthorParams{
		Name:     "Brian Kernighan",
		AuthorID: 11,
	})
	if err != nil {
		return err
	}

	log.Println("Inserted Author:")
	log.Println(insertedAuthor)

	// list all authors
	authors, err = queries.ListAuthor(ctx)
	if err != nil {
		return err
	}

	log.Println("Authors list:")
	log.Println(authors)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.AuthorID)
	if err != nil {
		return err
	}

	log.Println("Inserted Author == Fetched Author:")
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))

	book, err := queries.CreateBook(ctx, example.CreateBookParams{
		BookID:        1,
		AuthorID:      fetchedAuthor.AuthorID,
		Isbn:          "ABC123",
		BookType:      example.BooktypeBOOKTYPENONFICTION,
		Title:         "The C Programming Language",
		Year:          1983,
		AvailableTime: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Tags:          []string{"programming", "c"},
		Published:     pgtype.Bool{Bool: true, Valid: true},
		Price:         pgtype.Float8{Float64: 19.99, Valid: true},
	})

	if err != nil {
		return fmt.Errorf("create book: %w", err)
	}

	log.Println("Inserted book")
	log.Println(book)

	if err = queries.DeleteAuthor(ctx, bob.AuthorID); err != nil {
		return err
	}

	err = queries.DeleteAuthor(ctx, insertedAuthor.AuthorID)
	if err == nil {
		return errors.New("deleted an author still referenced by a book")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23503" {
		return fmt.Errorf("expected a foreign key violation deleting the author: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
