// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package converter

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func ConvertFrom(rd io.Reader) (*pluginpb.CodeGeneratorResponse, error) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	input, err := io.ReadAll(rd)
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %w", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(input, req)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal input: %w", err)
	}

	return Convert(req)
}

func Convert(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	files := []*pluginpb.CodeGeneratorResponse_File{}
	genFiles := make(map[string]struct{}, len(req.FileToGenerate))
	for _, file := range req.FileToGenerate {
		genFiles[file] = struct{}{}
	}

	// We need this to resolve dependencies when making protodesc versions of the files
	resolver, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: req.GetProtoFile(),
	})
	if err != nil {
		return nil, err
	}

	for _, fileDesc := range req.GetProtoFile() {
		if _, ok := genFiles[fileDesc.GetName()]; !ok {
			slog.Debug("skip generating file because it wasn't requested", slog.String("name", fileDesc.GetName()))
			continue
		}

		slog.Debug("generating file", slog.String("name", fileDesc.GetName()))

		_, err := protodesc.NewFile(fileDesc, resolver)
		if err != nil {
			slog.Error("error loading file", slog.Any("error", err))
			return nil, err
		}
	}

	features := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	return &pluginpb.CodeGeneratorResponse{
		File:              files,
		SupportedFeatures: &features,
	}, nil
}

func generateSchema(p *protogen.Plugin, opts template.Options) {
	slog.Debug("generating schema.pb.sql")
	gf := p.NewGeneratedFile("schema.pb.sql", "")

	err := template.ApplySchema(gf, template.SchemaParams{Schema: schema, Options: opts, HeaderParams: template.HeaderParams{}})
	if err != nil {
		gf.Skip()
		p.Error(err)
	}
}

func generateQueries(p *protogen.Plugin, opts template.Options) {
	for _, name := range p.Request.FileToGenerate {
		f := p.FilesByPath[name]

		if len(f.Messages) == 0 {
			slog.Debug("skip generating file because it has no messages", slog.String("name", name))
			continue
		}

		slog.Debug("processing queries for file", slog.String("name", name))
		slog.Debug("generating queries in", slog.String("name", fmt.Sprintf("%s.pb.sql", f.GeneratedFilenamePrefix)))

		gf := p.NewGeneratedFile(fmt.Sprintf("%s.pb.sql", f.GeneratedFilenamePrefix), f.GoImportPath)

		err := template.ApplyCrud(gf, template.CrudParams{
			GoName:       f.Messages[0].GoIdent.GoName,
			PrimaryKey:   "id",
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
