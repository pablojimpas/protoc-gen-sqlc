// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package core

type Schema struct {
	Tables    []Table
	Enums     []Enum
	Sequences []Sequence
}

type Sequence struct {
	Name      string
	Start     int
	Increment int
	MinValue  int
	MaxValue  int
}

type Enum struct {
	Name   string
	Values []string
}

type Table struct {
	Name        string
	Columns     []Column
	Constraints []Constraint
	Indexes     []Index
}

type Index struct {
	Name    string
	Columns []string
}

type Column struct {
	Name         string
	Type         ColumnType
	NotNull      bool
	DefaultValue string
}

type ColumnType string

const (
	IntegerType      ColumnType = "INTEGER"
	TextType         ColumnType = "TEXT"
	SerialType       ColumnType = "SERIAL"
	DateType         ColumnType = "DATE"
	TimestampType    ColumnType = "TIMESTAMP"
	VarcharType      ColumnType = "VARCHAR"
	VarcharArrayType ColumnType = "VARCHAR[]"
	TextArrayType    ColumnType = "TEXT[]"
	JSONBType        ColumnType = "JSONB"
	UUIDType         ColumnType = "UUID"
	BytesType        ColumnType = "BYTES"
	FloatType        ColumnType = "FLOAT"
	BooleanType      ColumnType = "BOOLEAN"
)

type Constraint struct {
	Type       ConstraintType
	Columns    []string
	References *Reference
}

type ConstraintType string

const (
	PrimaryKeyConstraint ConstraintType = "PRIMARY KEY"
	ForeignKeyConstraint ConstraintType = "FOREIGN KEY"
	UniqueConstraint     ConstraintType = "UNIQUE"
)

type Reference struct {
	Table    string
	Columns  []string
	OnDelete string
}
