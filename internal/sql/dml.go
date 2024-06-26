// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

import "fmt"

// LockMode represents the mode of a LOCK TABLE command.
type LockMode string

// Constants for common lock modes.
const (
	LockModeRead  LockMode = "READ"
	LockModeWrite LockMode = "WRITE"
)

// DML is a container for all DML commands.
type DML struct {
	Insert      *Insert
	Update      *Update
	Delete      *Delete
	Lock        *Lock
	Call        *Call
	ExplainPlan *ExplainPlan
}

// NewDML returns a new DML container.
func NewDML() *DML {
	return &DML{}
}

// Insert represents the INSERT INTO command.
type Insert struct {
	TableName TableName
	Columns   []ColumnName
	Values    []interface{}
}

// NewInsert initializes a new Insert command.
func NewInsert(tableName TableName, columns []ColumnName, values []interface{}) *Insert {
	return &Insert{TableName: tableName, Columns: columns, Values: values}
}

// Validate validates the Insert command.
func (i *Insert) Validate() error {
	if i.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if len(i.Columns) == 0 {
		return fmt.Errorf("columns cannot be empty")
	}
	if len(i.Values) == 0 {
		return fmt.Errorf("values cannot be empty")
	}
	if len(i.Columns) != len(i.Values) {
		return fmt.Errorf("columns and values must have the same length")
	}
	return nil
}

// Update represents the UPDATE command.
type Update struct {
	TableName TableName
	Set       map[ColumnName]interface{}
	Where     *Condition
}

// NewUpdate initializes a new Update command.
func NewUpdate(tableName TableName, set map[ColumnName]interface{}, where *Condition) *Update {
	return &Update{TableName: tableName, Set: set, Where: where}
}

// Validate validates the Update command.
func (u *Update) Validate() error {
	if u.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if len(u.Set) == 0 {
		return fmt.Errorf("set clauses cannot be empty")
	}
	for col, val := range u.Set {
		if col == "" {
			return fmt.Errorf("column name in set clause cannot be empty")
		}
		if val == nil {
			return fmt.Errorf("value in set clause cannot be nil")
		}
	}
	if u.Where != nil {
		if err := u.Where.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Delete represents the DELETE command.
type Delete struct {
	TableName TableName
	Where     *Condition
}

// NewDelete initializes a new Delete command.
func NewDelete(tableName TableName, where *Condition) *Delete {
	return &Delete{TableName: tableName, Where: where}
}

// Validate validates the Delete command.
func (d *Delete) Validate() error {
	if d.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if d.Where != nil {
		if err := d.Where.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Lock represents the LOCK TABLE command.
type Lock struct {
	TableName TableName
	LockMode  LockMode
}

// NewLock initializes a new Lock command.
func NewLock(tableName TableName, lockMode LockMode) *Lock {
	return &Lock{TableName: tableName, LockMode: lockMode}
}

// Validate validates the Lock command.
func (l *Lock) Validate() error {
	if l.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if l.LockMode != LockModeRead && l.LockMode != LockModeWrite {
		return fmt.Errorf("lock mode must be READ or WRITE")
	}
	return nil
}

// Call represents the CALL command.
type Call struct {
	ProcedureName string
	Arguments     []interface{}
}

// NewCall initializes a new Call command.
func NewCall(procedureName string, arguments []interface{}) *Call {
	return &Call{ProcedureName: procedureName, Arguments: arguments}
}

// Validate validates the Call command.
func (c *Call) Validate() error {
	if c.ProcedureName == "" {
		return fmt.Errorf("procedure name cannot be empty")
	}
	return nil
}

// ExplainPlan represents the EXPLAIN PLAN command.
type ExplainPlan struct {
	Query string
}

// NewExplainPlan initializes a new ExplainPlan command.
func NewExplainPlan(query string) *ExplainPlan {
	return &ExplainPlan{Query: query}
}

// Validate validates the ExplainPlan command.
func (e *ExplainPlan) Validate() error {
	if e.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	return nil
}
