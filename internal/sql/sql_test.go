// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package sql_test

import (
	"testing"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/sql"
)

func TestExampleSchema(t *testing.T) {
	// Define the authors table
	authorsColumns := []sql.ColumnDefinition{
		{Name: "author_id", DataType: sql.SerialType, NotNull: true, PrimaryKey: true},
		{Name: "name", DataType: sql.TextType, NotNull: true},
		{Name: "biography", DataType: sql.JSONBType},
	}
	createAuthorsTable := sql.NewCreate("authors", authorsColumns)

	// Validate the authors table
	if err := createAuthorsTable.Validate(); err != nil {
		t.Logf("Error validating authors table: %v\n", err)
	} else {
		t.Log("Authors table validated successfully")
	}

	// Define the books table
	booksColumns := []sql.ColumnDefinition{
		{Name: "book_id", DataType: sql.SerialType, NotNull: true, PrimaryKey: true},
		{Name: "author_id", DataType: sql.IntType, NotNull: true},
		{Name: "isbn", DataType: sql.TextType, NotNull: true},
		{Name: "title", DataType: sql.TextType, NotNull: true},
		{Name: "year", DataType: sql.IntType, NotNull: true},
		{Name: "available", DataType: sql.TimestampType, NotNull: true},
		{Name: "tags", DataType: sql.VarcharArrayType, NotNull: true},
	}
	createBooksTable := sql.NewCreate("books", booksColumns)

	// Validate the books table
	if err := createBooksTable.Validate(); err != nil {
		t.Logf("Error validating books table: %v\n", err)
	} else {
		t.Log("Books table validated successfully")
	}
}
