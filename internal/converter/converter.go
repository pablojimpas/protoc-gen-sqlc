// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

// Package converter transforms Protocol Buffer definitions into SQL schemas and queries.
package converter

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/pablojimpas/protoc-gen-sqlc/internal/core"
	sqlcpb "github.com/pablojimpas/protoc-gen-sqlc/internal/gen/sqlc"
	"github.com/pablojimpas/protoc-gen-sqlc/internal/sqlc/template"
)

const (
	// Number of parts expected in a "references" string. (table.column).
	referencesPartCount = 2
)

var (
	ErrNilEnum       = errors.New("nil enum provided")
	ErrNilMessage    = errors.New("nil message provided")
	ErrNilOptions    = errors.New("nil options provided")
	ErrTableNotFound = errors.New("table not found")
)

// SchemaBuilder transforms protobuf definitions into SQL schema structures.
type SchemaBuilder struct {
	Schema         core.Schema
	FilesByMessage map[string]*protogen.File
}

// NewSchemaBuilder creates a new SchemaBuilder with initialized fields.
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		Schema:         core.Schema{},
		FilesByMessage: make(map[string]*protogen.File),
	}
}

// Build processes all proto files in the plugin request to build a complete schema.
func (sb *SchemaBuilder) Build(p *protogen.Plugin) error {
	if p == nil {
		return errors.New("nil plugin provided")
	}

	for _, name := range p.Request.GetFileToGenerate() {
		f := p.FilesByPath[name]

		if len(f.Messages) == 0 {
			slog.Debug("skip generating file because it has no messages", slog.String("name", name))

			continue
		}

		slog.Debug("processing file", slog.String("name", name))

		for _, enum := range f.Enums {
			if err := sb.buildEnum(enum); err != nil {
				slog.Warn(
					"failed to build enum",
					slog.String("name", string(enum.Desc.Name())),
					slog.String("error", err.Error()),
				)

				continue
			}
		}

		for _, message := range f.Messages {
			sb.FilesByMessage[string(message.Desc.Name())] = f

			if err := sb.buildMessage(message); err != nil {
				slog.Warn(
					"failed to build message",
					slog.String("name", string(message.Desc.Name())),
					slog.String("error", err.Error()),
				)

				continue
			}
		}
	}

	return nil
}

// buildEnum converts a protobuf enum to a SQL enum.
func (sb *SchemaBuilder) buildEnum(protoEnum *protogen.Enum) error {
	if protoEnum == nil {
		return ErrNilEnum
	}

	values := make([]string, 0, len(protoEnum.Values))
	for _, v := range protoEnum.Values {
		values = append(values, string(v.Desc.Name()))
	}

	enum := core.Enum{
		Name:   protoEnum.GoIdent.GoName,
		Values: values,
	}

	sb.Schema.Enums = append(sb.Schema.Enums, enum)

	return nil
}

// buildMessage converts a protobuf message to a SQL table.
func (sb *SchemaBuilder) buildMessage(protoMessage *protogen.Message) error {
	if protoMessage == nil {
		return ErrNilMessage
	}

	name := protoMessage.Desc.Name()

	columns, err := buildColumns(protoMessage)
	if err != nil {
		return fmt.Errorf("building columns: %w", err)
	}

	constraints, err := buildConstraints(protoMessage)
	if err != nil {
		return fmt.Errorf("building constraints: %w", err)
	}

	table := core.Table{
		Name:        string(name),
		Columns:     columns,
		Constraints: constraints,
	}

	sb.Schema.Tables = append(sb.Schema.Tables, table)

	return nil
}

// buildColumns converts protobuf message fields to SQL columns.
func buildColumns(protoMessage *protogen.Message) ([]core.Column, error) {
	if protoMessage == nil {
		return nil, ErrNilMessage
	}

	columns := make([]core.Column, 0, len(protoMessage.Fields))

	for _, field := range protoMessage.Fields {
		columnType, err := mapDataType(field)
		if err != nil {
			slog.Warn("error mapping data type",
				slog.String("field", string(field.Desc.Name())),
				slog.String("error", err.Error()))

			continue
		}

		column := &core.Column{
			Name: string(field.Desc.Name()),
			Type: columnType,
		}

		if err := applyExtensions(field.Desc.Options(), column); err != nil {
			slog.Warn("error applying extensions",
				slog.String("field", string(field.Desc.Name())),
				slog.String("error", err.Error()))
		}

		// Format default values appropriately based on type
		if column.DefaultValue != "" {
			if column.Type != core.IntegerType && column.Type != core.FloatType &&
				column.Type != core.BooleanType {
				column.DefaultValue = fmt.Sprintf("'%v'", column.DefaultValue)
			}
		}

		columns = append(columns, *column)
	}

	return columns, nil
}

// applyExtensions applies proto extensions to a column definition.
func applyExtensions(opts protoreflect.ProtoMessage, column *core.Column) error {
	if opts == nil {
		return ErrNilOptions
	}

	if column == nil {
		return errors.New("nil column provided")
	}

	// Apply SQLC field extensions
	if proto.HasExtension(opts, sqlcpb.E_Field) {
		ext, ok := proto.GetExtension(opts, sqlcpb.E_Field).(*sqlcpb.FieldConstraints)
		if !ok {
			slog.Warn("failed to get sqlc field extension")
		} else {
			column.DefaultValue = ext.GetDefault()
			column.NotNull = ext.GetPrimary()
		}
	}

	// Apply validate extensions
	if proto.HasExtension(opts, validate.E_Field) {
		ext, ok := proto.GetExtension(opts, validate.E_Field).(*validate.FieldConstraints)
		if !ok {
			slog.Warn("failed to get validate field extension")
		} else {
			column.NotNull = column.NotNull || ext.GetRequired()

			// Check for UUID type
			stringRules := ext.GetString()
			if stringRules != nil && stringRules.GetUuid() {
				column.Type = core.UUIDType
			}
		}
	}

	return nil
}

// buildConstraints extracts SQL constraints from protobuf message fields.
func buildConstraints(protoMessage *protogen.Message) ([]core.Constraint, error) {
	if protoMessage == nil {
		return nil, ErrNilMessage
	}

	var constraints []core.Constraint

	for _, field := range protoMessage.Fields {
		opts := field.Desc.Options()
		if !proto.HasExtension(opts, sqlcpb.E_Field) {
			continue
		}

		ext, ok := proto.GetExtension(opts, sqlcpb.E_Field).(*sqlcpb.FieldConstraints)
		if !ok {
			slog.Warn(
				"invalid extension type for field",
				slog.String("field", string(field.Desc.Name())),
			)

			continue
		}

		fieldName := string(field.Desc.Name())

		// Handle unique constraint
		if ext.GetUnique() {
			constraints = append(constraints, core.Constraint{
				Type:    core.UniqueConstraint,
				Columns: []string{fieldName},
			})
		}

		// Handle primary key constraint
		if ext.GetPrimary() {
			constraints = append(constraints, core.Constraint{
				Type:    core.PrimaryKeyConstraint,
				Columns: []string{fieldName},
			})
		}

		// Handle foreign key constraint
		if ref := ext.GetReferences(); ref != "" {
			parts := strings.Split(ref, ".")
			if len(parts) != referencesPartCount {
				slog.Warn("invalid references format", slog.String("value", ref))

				continue
			}

			constraints = append(constraints, core.Constraint{
				Type:    core.ForeignKeyConstraint,
				Columns: []string{fieldName},
				References: &core.Reference{
					Table:    parts[0],
					Columns:  []string{parts[1]},
					OnDelete: core.ForeignKeyActionNoAction,
					OnUpdate: core.ForeignKeyActionNoAction,
				},
			})
		}
	}

	return constraints, nil
}

// mapDataType converts protobuf field types to SQL column types.
func mapDataType(field *protogen.Field) (core.ColumnType, error) {
	if field == nil {
		return "", errors.New("nil field provided")
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return core.BooleanType, nil
	case protoreflect.BytesKind:
		return core.BytesType, nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return core.FloatType, nil
	case protoreflect.StringKind:
		return mapStringType(field)
	case protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return core.IntegerType, nil
	case protoreflect.EnumKind:
		return core.ColumnType(field.Enum.Desc.Name()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return mapMessageType(field)
	default:
		return core.BytesType, nil
	}
}

func mapStringType(field *protogen.Field) (core.ColumnType, error) {
	if field.Desc.Cardinality() == protoreflect.Repeated {
		return core.TextArrayType, nil
	}

	return core.TextType, nil
}

func mapMessageType(field *protogen.Field) (core.ColumnType, error) {
	if field.Message == nil || field.Message.Desc == nil {
		return "", errors.New("message field has nil descriptor")
	}

	switch field.Message.Desc.FullName() {
	case "google.protobuf.Timestamp":
		return core.TimestampType, nil
	case "google.protobuf.Struct":
		return core.JSONBType, nil
	default:
		return core.BytesType, nil
	}
}

// GenerateSchema creates a SQL schema file from the accumulated schema definition.
func GenerateSchema(
	p *protogen.Plugin,
	schema core.Schema,
	tmpl *template.Templates,
	opts template.Options,
) error {
	if p == nil {
		return errors.New("nil plugin provided")
	}

	slog.Debug("generating schema.sql")

	gf := p.NewGeneratedFile("schema.sql", "")

	err := tmpl.ApplySchema(gf, &template.SchemaParams{
		Schema:       schema,
		Options:      opts,
		HeaderParams: template.HeaderParams{},
	})
	if err != nil {
		gf.Skip()
		p.Error(err)

		return fmt.Errorf("applying schema template: %w", err)
	}

	return nil
}

// GenerateQueries creates SQL query files for each message in the schema.
func GenerateQueries(
	p *protogen.Plugin,
	schema core.Schema,
	filesByMessage map[string]*protogen.File,
	tmpl *template.Templates,
	opts template.Options,
) error {
	if p == nil {
		return errors.New("nil plugin provided")
	}

	for message, protoFile := range filesByMessage {
		if protoFile == nil {
			slog.Warn("nil proto file for message", slog.String("message", message))

			continue
		}

		name := protoFile.Proto.GetName()
		slog.Debug("processing queries for file", slog.String("name", name))
		slog.Debug(
			"generating queries in",
			slog.String("name", protoFile.GeneratedFilenamePrefix+".sql"),
		)

		gf := p.NewGeneratedFile(protoFile.GeneratedFilenamePrefix+".sql", protoFile.GoImportPath)

		table := schema.TableByName(message)
		if table == nil {
			slog.Warn("table not found for message", slog.String("message", message))

			continue
		}

		err := tmpl.ApplyCrud(gf, &template.CrudParams{
			GoName:       message,
			PrimaryKey:   table.PrimaryKey(),
			Table:        *table,
			Options:      opts,
			HeaderParams: template.HeaderParams{Sources: []string{protoFile.Proto.GetName()}},
		})
		if err != nil {
			gf.Skip()
			p.Error(err)

			continue
		}
	}

	return nil
}
