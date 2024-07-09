package gen_test

import (
	"bytes"
	"testing"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/gen"
)

func TestApplySchemaTemplate(t *testing.T) {
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
					{Type: core.ForeignKeyConstraint, Columns: []string{"author_id"}, References: &core.Reference{Table: "authors", Columns: []string{"id"}}},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := gen.ApplySchemaTemplate(&buf, gen.SchemaParams{schema, gen.Options{}, gen.HeaderParams{}})
	if err != nil {
		t.Error(err)
	}
}
