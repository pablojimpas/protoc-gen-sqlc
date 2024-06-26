// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

import "fmt"

// OrderDirection represents the direction of an ORDER BY clause.
type OrderDirection string

// Constants for common ORDER BY directions.
const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

// Operator represents a comparison operator in a WHERE condition.
type Operator string

// Constants for common comparison operators.
const (
	Equal        Operator = "="
	NotEqual     Operator = "!="
	GreaterThan  Operator = ">"
	LessThan     Operator = "<"
	GreaterEqual Operator = ">="
	LessEqual    Operator = "<="
)

// DQL is a container for all DQL commands.
type DQL struct {
	Select *Select
}

// NewDQL returns a new DQL container.
func NewDQL() *DQL {
	return &DQL{}
}

// Select represents the SELECT command.
type Select struct {
	Columns   []ColumnName
	TableName TableName
	Where     *Condition
	OrderBy   []Order
	Limit     int
}

// NewSelect initializes a new Select command.
func NewSelect(columns []ColumnName, tableName TableName, where *Condition, orderBy []Order, limit int) *Select {
	return &Select{
		Columns:   columns,
		TableName: tableName,
		Where:     where,
		OrderBy:   orderBy,
		Limit:     limit,
	}
}

// Validate validates the Select command.
func (s *Select) Validate() error {
	if len(s.Columns) == 0 {
		return fmt.Errorf("columns cannot be empty")
	}
	if s.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if s.Where != nil {
		if err := s.Where.Validate(); err != nil {
			return err
		}
	}
	for _, order := range s.OrderBy {
		if err := order.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Order represents an ORDER BY clause in a SELECT statement.
type Order struct {
	Column    ColumnName
	Direction OrderDirection
}

// Validate validates the Order clause.
func (o *Order) Validate() error {
	if o.Column == "" {
		return fmt.Errorf("column cannot be empty")
	}
	if o.Direction != Asc && o.Direction != Desc {
		return fmt.Errorf("direction must be ASC or DESC")
	}
	return nil
}
