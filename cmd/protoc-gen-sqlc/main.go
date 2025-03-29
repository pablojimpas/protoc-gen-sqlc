// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/converter"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
)

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(
		func(p *protogen.Plugin) error {
			p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
			tmpl := template.New()
			opts := template.Options{}
			sb := converter.NewSchemaBuilder()

			if err := sb.Build(p); err != nil {
				return err
			}

			if err := converter.GenerateSchema(p, sb.Schema, tmpl, opts); err != nil {
				return err
			}

			if err := converter.GenerateQueries(p, sb.Schema, sb.FilesByMessage, tmpl, opts); err != nil {
				return err
			}

			return nil
		},
	)
}
