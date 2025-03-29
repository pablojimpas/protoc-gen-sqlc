// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package template_test

import (
	"bytes"
	"testing"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
)

func TestApplySchemaTemplate(t *testing.T) {
	t.Parallel()

	schema := core.Schema{
		Tables: []core.Table{
			{
				Name: "authors",
				Columns: []core.Column{
					{Name: "id", Type: core.SerialType, NotNull: true},
					{Name: "name", Type: core.TextType, NotNull: true},
					{Name: "bio", Type: core.TextType},
				},
				Constraints: []core.Constraint{
					{Type: core.PrimaryKeyConstraint, Columns: []string{"id"}},
				},
			},
			{
				Name: "books",
				Columns: []core.Column{
					{Name: "id", Type: core.SerialType, NotNull: true},
					{Name: "title", Type: core.TextType, NotNull: true},
					{Name: "author_id", Type: core.IntegerType, NotNull: true},
				},
				Constraints: []core.Constraint{
					{Type: core.PrimaryKeyConstraint, Columns: []string{"id"}},
					{
						Type:       core.ForeignKeyConstraint,
						Columns:    []string{"author_id"},
						References: &core.Reference{Table: "authors", Columns: []string{"id"}},
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	tmpl := template.New()

	err := tmpl.ApplySchema(
		&buf,
		&template.SchemaParams{schema, template.Options{}, template.HeaderParams{}},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestApplyCrudTemplate(t *testing.T) {
	t.Parallel()

	table := core.Table{
		Name: "books",
		Columns: []core.Column{
			{Name: "id", Type: core.SerialType, NotNull: true},
			{Name: "title", Type: core.TextType, NotNull: true},
			{Name: "author_id", Type: core.IntegerType, NotNull: true},
		},
		Constraints: []core.Constraint{
			{Type: core.PrimaryKeyConstraint, Columns: []string{"id"}},
			{
				Type:       core.ForeignKeyConstraint,
				Columns:    []string{"author_id"},
				References: &core.Reference{Table: "authors", Columns: []string{"id"}},
			},
		},
	}

	var buf bytes.Buffer

	tmpl := template.New()

	err := tmpl.ApplyCrud(
		&buf,
		&template.CrudParams{
			GoName:       "Book",
			PrimaryKey:   "id",
			Table:        table,
			Options:      template.Options{},
			HeaderParams: template.HeaderParams{},
		},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestHeaderTemplate(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	tmpl := template.New()

	err := tmpl.ApplySchema(
		&buf,
		&template.SchemaParams{
			HeaderParams: template.HeaderParams{
				Sources: []string{"source1.proto", "source2.proto"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
}
