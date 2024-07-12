// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/converter"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(
		func(p *protogen.Plugin) error {
			p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
			opts := template.Options{}
			sb := converter.NewSchemaBuilder()
			sb.Build(p)
			converter.GenerateSchema(p, sb.Schema, opts)
			converter.GenerateQueries(p, sb.Schema, opts)
			return nil
		},
	)
}
