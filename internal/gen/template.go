// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package gen

import (
	"embed"
	"io"
	"log"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed *.tmpl
var files embed.FS

var (
	headerTmpl  = parse("header.tmpl")
	messageTmpl = parse("crud.tmpl")
)

func parse(file string) *template.Template {
	return template.Must(template.New(file).Funcs(sprig.TxtFuncMap()).ParseFS(files, file))
}

// Options are the options to set for rendering the template.
type Options struct{}

type headerParams struct {
	*protogen.File
}

type messageParams struct {
	*protogen.Message
	Options
}

// This function is called with a param which contains the entire definition of a method.
func ApplyTemplate(w io.Writer, f *protogen.File, opts Options) error {
	if err := headerTmpl.Execute(w, headerParams{
		File: f,
	}); err != nil {
		return err
	}

	return applyMessages(w, f.Messages, opts)
}

func applyMessages(w io.Writer, msgs []*protogen.Message, opts Options) error {
	for _, m := range msgs {
		if m.Desc.IsMapEntry() {
			log.Printf("Skipping %s, mapentry message", m.GoIdent.GoName)
			continue
		}

		log.Printf("Processing %s", m.GoIdent.GoName)
		if err := messageTmpl.Execute(w, messageParams{
			Message: m,
			Options: opts,
		}); err != nil {
			return err
		}

		if err := applyMessages(w, m.Messages, opts); err != nil {
			return err
		}
	}

	return nil
}
