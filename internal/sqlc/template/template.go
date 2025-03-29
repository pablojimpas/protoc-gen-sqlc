// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package template

import (
	"embed"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
)

//go:embed *.tmpl
var files embed.FS

// Templates holds the compiled templates for code generation.
type Templates struct {
	header *template.Template
	schema *template.Template
	crud   *template.Template
}

// New creates a new set of initialized templates.
func New() *Templates {
	return &Templates{
		header: parse("header.tmpl"),
		schema: parse("schema.tmpl"),
		crud:   parse("crud.tmpl"),
	}
}

func parse(file string) *template.Template {
	return template.Must(template.New(file).Funcs(sprig.TxtFuncMap()).ParseFS(files, file))
}

type Options struct{}

type HeaderParams struct {
	Sources []string
}

type SchemaParams struct {
	core.Schema
	Options
	HeaderParams
}

type CrudParams struct {
	GoName     string
	PrimaryKey string
	core.Table
	Options
	HeaderParams
}

// ApplySchema applies the schema template with the provided parameters.
func (t *Templates) ApplySchema(w io.Writer, p *SchemaParams) error {
	if err := t.header.Execute(w, p.HeaderParams); err != nil {
		return err
	}

	return t.schema.Execute(w, p)
}

// ApplyCrud applies the CRUD template with the provided parameters.
func (t *Templates) ApplyCrud(w io.Writer, p *CrudParams) error {
	if err := t.header.Execute(w, p.HeaderParams); err != nil {
		return err
	}

	return t.crud.Execute(w, p)
}
