// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package converter

import (
	"fmt"
	"log/slog"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	sqlcpb "github.com/pablojimpas/protoc-gen-sqlc/internal/gen/sqlc"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var FilesByMessage = make(map[string]*protogen.File)

type SchemaBuilder struct {
	Schema core.Schema
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		Schema: core.Schema{},
	}
}

func (sb *SchemaBuilder) Build(p *protogen.Plugin) {
	for _, name := range p.Request.FileToGenerate {
		f := p.FilesByPath[name]

		if len(f.Messages) == 0 {
			slog.Debug("skip generating file because it has no messages", slog.String("name", name))
			continue
		}

		slog.Debug("processing file", slog.String("name", name))

		for _, enum := range f.Enums {
			sb.BuildEnum(enum)
		}

		for _, message := range f.Messages {
			FilesByMessage[string(message.Desc.Name())] = f
			sb.BuildMessage(message)
		}
	}
}

func (sb *SchemaBuilder) BuildEnum(m *protogen.Enum) {
	var values []string
	for _, v := range m.Values {
		values = append(values, string(v.Desc.Name()))
	}
	enum := core.Enum{
		Name:   m.GoIdent.GoName,
		Values: values,
	}

	sb.Schema.Enums = append(sb.Schema.Enums, enum)
}

func (sb *SchemaBuilder) BuildMessage(m *protogen.Message) {
	name := m.Desc.Name()
	columns := buildColumns(m)
	constraints := buildConstraints(m)
	table := core.Table{
		Name:        string(name),
		Columns:     columns,
		Constraints: constraints,
	}

	sb.Schema.Tables = append(sb.Schema.Tables, table)
}

func buildColumns(m *protogen.Message) []core.Column {
	var columns []core.Column
	for _, field := range m.Fields {
		c := &core.Column{
			Name: string(field.Desc.Name()),
			Type: mapDataType(field),
		}
		getExtensions(field.Desc.Options(), c)
		if c.DefaultValue != "" {
			if c.Type != core.IntegerType && c.Type != core.FloatType && c.Type != core.BooleanType {
				c.DefaultValue = fmt.Sprintf("'%v'", c.DefaultValue)
			}
		}

		columns = append(columns, *c)
	}
	return columns
}

func getExtensions(opts protoreflect.ProtoMessage, c *core.Column) {
	if proto.HasExtension(opts, sqlcpb.E_Field) {
		ext, ok := proto.GetExtension(opts, sqlcpb.E_Field).(*sqlcpb.FieldConstraints)
		if ok {
			c.DefaultValue = ext.Default
			c.NotNull = ext.Primary
		}
	}
	if proto.HasExtension(opts, validate.E_Field) {
		ext, ok := proto.GetExtension(opts, validate.E_Field).(*validate.FieldConstraints)
		if ok {
			c.NotNull = c.NotNull || *ext.Required
		}
		if ext.GetString_().GetUuid() {
			c.Type = core.UUIDType
		}
	}
}

func buildConstraints(m *protogen.Message) []core.Constraint {
	var constraints []core.Constraint
	for _, field := range m.Fields {
		opts := field.Desc.Options()
		if !proto.HasExtension(opts, sqlcpb.E_Field) {
			continue
		}
		ext, ok := proto.GetExtension(opts, sqlcpb.E_Field).(*sqlcpb.FieldConstraints)
		if !ok {
			continue
		}
		if ext.Unique {
			constraints = append(constraints, core.Constraint{
				Type:    core.UniqueConstraint,
				Columns: []string{string(field.Desc.Name())},
			})
		}
		if ext.Primary {
			constraints = append(constraints, core.Constraint{
				Type:    core.PrimaryKeyConstraint,
				Columns: []string{string(field.Desc.Name())},
			})
		}
		if ext.References != "" {
			parts := strings.Split(ext.References, ".")
			constraints = append(constraints, core.Constraint{
				Type:    core.ForeignKeyConstraint,
				Columns: []string{string(field.Desc.Name())},
				References: &core.Reference{
					Table:    parts[0],
					Columns:  []string{parts[1]},
					OnDelete: "CASCADE",
				},
			})
		}
	}
	return constraints
}

//nolint:cyclop // not much we can do to avoid this
func mapDataType(field *protogen.Field) core.ColumnType {
	var t core.ColumnType
	//nolint:exhaustive // not every type is relevant
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		t = core.BooleanType
	case protoreflect.BytesKind:
		t = core.BytesType
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		t = core.FloatType
	case protoreflect.StringKind:
		if field.Desc.Cardinality() == protoreflect.Repeated {
			t = core.TextArrayType
		} else {
			t = core.TextType
		}
	case protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		t = core.IntegerType
	case protoreflect.EnumKind:
		t = core.ColumnType(field.Enum.Desc.Name())
	case protoreflect.MessageKind:
		switch field.Message.Desc.FullName() {
		case "google.protobuf.Timestamp":
			t = core.TimestampType
		case "google.protobuf.Struct":
			t = core.JSONBType
		}
	default:
		t = core.BytesType
	}

	return t
}

func GenerateSchema(p *protogen.Plugin, s core.Schema, opts template.Options) {
	slog.Debug("generating schema.sql")
	gf := p.NewGeneratedFile("schema.sql", "")

	err := template.ApplySchema(gf, template.SchemaParams{
		Schema:       s,
		Options:      opts,
		HeaderParams: template.HeaderParams{},
	})
	if err != nil {
		gf.Skip()
		p.Error(err)
	}
}

func GenerateQueries(p *protogen.Plugin, s core.Schema, opts template.Options) {
	for message, f := range FilesByMessage {
		name := *f.Proto.Name
		slog.Debug("processing queries for file", slog.String("name", name))
		slog.Debug("generating queries in", slog.String("name", fmt.Sprintf("%s.sql", f.GeneratedFilenamePrefix)))
		gf := p.NewGeneratedFile(fmt.Sprintf("%s.sql", f.GeneratedFilenamePrefix), f.GoImportPath)

		table := *s.TableByName(message)
		err := template.ApplyCrud(gf, template.CrudParams{
			GoName:       message,
			PrimaryKey:   table.PrimaryKey(),
			Table:        table,
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
