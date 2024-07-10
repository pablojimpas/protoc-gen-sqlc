// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package converter

import (
	"fmt"
	"log/slog"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
	"google.golang.org/protobuf/compiler/protogen"
)

func GenerateSchema(p *protogen.Plugin, opts template.Options) {
	slog.Debug("generating schema.sql")
	gf := p.NewGeneratedFile("schema.sql", "")

	err := template.ApplySchema(gf, template.SchemaParams{Schema: schema, Options: opts, HeaderParams: template.HeaderParams{}})
	if err != nil {
		gf.Skip()
		p.Error(err)
	}
}

func GenerateQueries(p *protogen.Plugin, opts template.Options) {
	for _, name := range p.Request.FileToGenerate {
		f := p.FilesByPath[name]

		if len(f.Messages) == 0 {
			slog.Debug("skip generating file because it has no messages", slog.String("name", name))
			continue
		}

		slog.Debug("processing queries for file", slog.String("name", name))
		slog.Debug("generating queries in", slog.String("name", fmt.Sprintf("%s.sql", f.GeneratedFilenamePrefix)))

		gf := p.NewGeneratedFile(fmt.Sprintf("%s.sql", f.GeneratedFilenamePrefix), f.GoImportPath)

		err := template.ApplyCrud(gf, template.CrudParams{
			GoName:       f.Messages[0].GoIdent.GoName,
			PrimaryKey:   schema.Tables[0].Constraints[0].Columns[0],
			Table:        schema.Tables[0],
			Options:      opts,
			HeaderParams: template.HeaderParams{Sources: []string{*f.Proto.Name}},
		})
		if err != nil {
			gf.Skip()
			p.Error(err)
			continue
		}
	}
}

var schema = core.Schema{
	Enums: []core.Enum{
		{
			Name:   "book_type",
			Values: []string{"FICTION", "NONFICTION"},
		},
	},
	Tables: []core.Table{
		{
			Name: "authors",
			Columns: []core.Column{
				{Name: "author_id", Type: core.SerialType, NotNull: true},
				{Name: "name", Type: core.TextType, NotNull: true, DefaultValue: "'Anonymous'"},
				{Name: "biography", Type: core.JSONBType},
			},
			Constraints: []core.Constraint{
				{Type: core.PrimaryKeyConstraint, Columns: []string{"author_id"}},
			},
		},
		{
			Name: "books",
			Columns: []core.Column{
				{Name: "book_id", Type: core.SerialType, NotNull: true},
				{Name: "author_id", Type: core.IntegerType, NotNull: true},
				{Name: "isbn", Type: core.TextType, NotNull: true},
				{Name: "book_type", Type: "book_type", NotNull: true, DefaultValue: "'FICTION'"},
				{Name: "title", Type: core.TextType, NotNull: true},
				{Name: "year", Type: core.IntegerType, NotNull: true, DefaultValue: "2000"},
				{Name: "available", Type: core.TimestampType, NotNull: true, DefaultValue: "'NOW()'"},
				{Name: "tags", Type: core.VarcharArrayType, NotNull: true, DefaultValue: "'{}'"},
			},
			Constraints: []core.Constraint{
				{Type: core.PrimaryKeyConstraint, Columns: []string{"book_id"}},
				{Type: core.ForeignKeyConstraint, Columns: []string{"author_id"}, References: &core.Reference{Table: "authors", Columns: []string{"author_id"}}},
				{Type: core.UniqueConstraint, Columns: []string{"isbn"}},
			},
		},
	},
}
