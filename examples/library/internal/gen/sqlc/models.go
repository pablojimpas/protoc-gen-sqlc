// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package example

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Booktype string

const (
	BooktypeBOOKTYPEUNSPECIFIED Booktype = "BOOK_TYPE_UNSPECIFIED"
	BooktypeBOOKTYPEFICTION     Booktype = "BOOK_TYPE_FICTION"
	BooktypeBOOKTYPENONFICTION  Booktype = "BOOK_TYPE_NONFICTION"
)

func (e *Booktype) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Booktype(s)
	case string:
		*e = Booktype(s)
	default:
		return fmt.Errorf("unsupported scan type for Booktype: %T", src)
	}
	return nil
}

type NullBooktype struct {
	Booktype Booktype
	Valid    bool // Valid is true if Booktype is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBooktype) Scan(value interface{}) error {
	if value == nil {
		ns.Booktype, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Booktype.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBooktype) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Booktype), nil
}

type Author struct {
	AuthorID  int32
	Name      string
	Biography []byte
}

type Book struct {
	BookID        int32
	AuthorID      int32
	Isbn          string
	BookType      Booktype
	Title         string
	Year          int32
	AvailableTime pgtype.Timestamptz
	Tags          []string
	Published     pgtype.Bool
	Price         pgtype.Float8
}
