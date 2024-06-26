// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

import "fmt"

// Command represents a generic SQL command interface.
type Command interface {
	Validate() error
}

// TableName represents a name of a table.
type TableName string

// ColumnName represents a name of a column.
type ColumnName string

// DataType represents a SQL data type.
type DataType string

// Constants for common SQL data types.
const (
	IntType          DataType = "INT"
	VarcharType      DataType = "VARCHAR"
	DateType         DataType = "DATE"
	SerialType       DataType = "SERIAL"
	TextType         DataType = "TEXT"
	JSONBType        DataType = "JSONB"
	TimestampType    DataType = "TIMESTAMP WITH TIME ZONE"
	VarcharArrayType DataType = "VARCHAR[]"
)

// NotNullConstraint represents a NOT NULL constraint.
type NotNullConstraint bool

// PrimaryKeyConstraint represents a PRIMARY KEY constraint.
type PrimaryKeyConstraint bool

// Condition represents a WHERE condition in DML and DQL statements.
type Condition struct {
	Column   ColumnName
	Operator Operator
	Value    interface{}
}

// NewCondition initializes a new Condition.
func NewCondition(column ColumnName, operator Operator, value interface{}) *Condition {
	return &Condition{
		Column:   column,
		Operator: operator,
		Value:    value,
	}
}

// Validate validates the Condition.
func (c *Condition) Validate() error {
	if c.Column == "" {
		return fmt.Errorf("column cannot be empty")
	}
	if c.Operator == "" {
		return fmt.Errorf("operator cannot be empty")
	}
	// Additional validation for value can be added based on the type and operator.
	return nil
}
