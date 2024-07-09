// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package template

import (
	"embed"
	"io"
	"text/template"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"

	"github.com/Masterminds/sprig/v3"
)

//go:embed *.tmpl
var files embed.FS

var (
	headerTmpl = parse("header.tmpl")
	schemaTmpl = parse("schema.tmpl")
	crudTmpl   = parse("crud.tmpl")
)

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

func ApplySchema(w io.Writer, p SchemaParams) error {
	if err := headerTmpl.Execute(w, p.HeaderParams); err != nil {
		return err
	}

	return schemaTmpl.Execute(w, p)
}

func ApplyCrud(w io.Writer, p CrudParams) error {
	if err := headerTmpl.Execute(w, p.HeaderParams); err != nil {
		return err
	}

	return crudTmpl.Execute(w, p)
}
