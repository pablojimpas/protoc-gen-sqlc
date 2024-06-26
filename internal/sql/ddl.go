// SPDX-FileCopyrightText: 2024 Pablo Jiménez Pascual <pablo@jimpas.me>
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

import "fmt"

// DDL is a container for all DDL commands.
type DDL struct {
	Create   *Create
	Drop     *Drop
	Alter    *Alter
	Truncate *Truncate
	Comment  *Comment
	Rename   *Rename
}

// NewDDL returns a new DDL container.
func NewDDL() *DDL {
	return &DDL{}
}

// Create represents the CREATE TABLE command.
type Create struct {
	TableName TableName
	Columns   []ColumnDefinition
}

// NewCreate initializes a new Create command.
func NewCreate(tableName TableName, columns []ColumnDefinition) *Create {
	return &Create{TableName: tableName, Columns: columns}
}

// Validate validates the Create command.
func (c *Create) Validate() error {
	if c.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if len(c.Columns) == 0 {
		return fmt.Errorf("columns cannot be empty")
	}
	for _, col := range c.Columns {
		if err := col.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ColumnDefinition represents a column in a CREATE TABLE command.
type ColumnDefinition struct {
	Name       ColumnName
	DataType   DataType
	NotNull    NotNullConstraint
	PrimaryKey PrimaryKeyConstraint
}

// Validate validates the ColumnDefinition.
func (cd *ColumnDefinition) Validate() error {
	if cd.Name == "" {
		return fmt.Errorf("column name cannot be empty")
	}
	if cd.DataType == "" {
		return fmt.Errorf("data type cannot be empty")
	}
	return nil
}

// Drop represents the DROP TABLE command.
type Drop struct {
	TableName TableName
}

// NewDrop initializes a new Drop command.
func NewDrop(tableName TableName) *Drop {
	return &Drop{TableName: tableName}
}

// Validate validates the Drop command.
func (d *Drop) Validate() error {
	if d.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	return nil
}

// Alter represents the ALTER TABLE command.
type Alter struct {
	TableName TableName
	Actions   []AlterAction
}

// NewAlter initializes a new Alter command.
func NewAlter(tableName TableName, actions []AlterAction) *Alter {
	return &Alter{TableName: tableName, Actions: actions}
}

// Validate validates the Alter command.
func (a *Alter) Validate() error {
	if a.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if len(a.Actions) == 0 {
		return fmt.Errorf("actions cannot be empty")
	}
	for _, action := range a.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// AlterAction represents an action to be taken in an ALTER TABLE command.
type AlterAction struct {
	ActionType AlterActionType
	Column     ColumnDefinition
}

// Validate validates the AlterAction.
func (aa *AlterAction) Validate() error {
	if aa.ActionType == "" {
		return fmt.Errorf("action type cannot be empty")
	}
	if err := aa.Column.Validate(); err != nil {
		return err
	}
	return nil
}

// AlterActionType defines the type of action in an ALTER TABLE command.
type AlterActionType string

// Constants for common ALTER TABLE action types.
const (
	AddColumn    AlterActionType = "ADD"
	DropColumn   AlterActionType = "DROP"
	ModifyColumn AlterActionType = "MODIFY"
)

// Truncate represents the TRUNCATE TABLE command.
type Truncate struct {
	TableName TableName
}

// NewTruncate initializes a new Truncate command.
func NewTruncate(tableName TableName) *Truncate {
	return &Truncate{TableName: tableName}
}

// Validate validates the Truncate command.
func (t *Truncate) Validate() error {
	if t.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	return nil
}

// Comment represents the COMMENT ON TABLE command.
type Comment struct {
	TableName   TableName
	CommentText string
}

// NewComment initializes a new Comment command.
func NewComment(tableName TableName, commentText string) *Comment {
	return &Comment{TableName: tableName, CommentText: commentText}
}

// Validate validates the Comment command.
func (c *Comment) Validate() error {
	if c.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if c.CommentText == "" {
		return fmt.Errorf("comment text cannot be empty")
	}
	return nil
}

// Rename represents the RENAME TABLE command.
type Rename struct {
	OldTableName TableName
	NewTableName TableName
}

// NewRename initializes a new Rename command.
func NewRename(oldTableName, newTableName TableName) *Rename {
	return &Rename{OldTableName: oldTableName, NewTableName: newTableName}
}

// Validate validates the Rename command.
func (r *Rename) Validate() error {
	if r.OldTableName == "" || r.NewTableName == "" {
		return fmt.Errorf("old and new table names cannot be empty")
	}
	return nil
}
