// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		opts := gen.Options{}
		generateSchema(p, opts)
		generateQueries(p, opts)
		return nil
	})
}

func generateSchema(p *protogen.Plugin, opts gen.Options) {
	log.Println("Generating schema.pb.sql")
	gf := p.NewGeneratedFile("schema.pb.sql", "")

	err := gen.ApplySchemaTemplate(gf, p.Files, opts)
	if err != nil {
		gf.Skip()
		p.Error(err)
	}
}

func generateQueries(p *protogen.Plugin, opts gen.Options) {
	for _, name := range p.Request.FileToGenerate {
		f := p.FilesByPath[name]

		if len(f.Messages) == 0 {
			log.Printf("Skipping %s, no messages", name)
			continue
		}

		log.Printf("Processing %s", name)
		log.Printf("Generating %s\n", fmt.Sprintf("%s.pb.sql", f.GeneratedFilenamePrefix))

		gf := p.NewGeneratedFile(fmt.Sprintf("%s.pb.sql", f.GeneratedFilenamePrefix), f.GoImportPath)

		err := gen.ApplyQueryTemplate(gf, f, opts)
		if err != nil {
			gf.Skip()
			p.Error(err)
			continue
		}
	}
}
