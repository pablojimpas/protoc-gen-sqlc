// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package core

type Schema struct {
	Tables    []Table
	Enums     []Enum
	Sequences []Sequence
}

func (s *Schema) TableByName(name string) *Table {
	for i, t := range s.Tables {
		if t.Name == name {
			return &s.Tables[i]
		}
	}

	return nil
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

func (s *Table) PrimaryKey() string {
	for _, c := range s.Constraints {
		if c.Type == PrimaryKeyConstraint {
			return c.Columns[0]
		}
	}

	return "id"
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
	TimestampType    ColumnType = "TIMESTAMPTZ"
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
	OnDelete ForeignKeyAction
	OnUpdate ForeignKeyAction
}

type ForeignKeyAction string

const (
	ForeignKeyActionNoAction   ForeignKeyAction = "NO ACTION"
	ForeignKeyActionRestrict   ForeignKeyAction = "RESTRICT"
	ForeignKeyActionCascade    ForeignKeyAction = "CASCADE"
	ForeignKeyActionSetNull    ForeignKeyAction = "SET NULL"
	ForeignKeyActionSetDefault ForeignKeyAction = "SET DEFAULT"
)
