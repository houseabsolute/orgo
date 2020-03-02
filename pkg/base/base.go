package base

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type Schema struct {
	db   *goqu.Database
	name string
}

func New(name string, db *goqu.Database) *Schema {
	return &Schema{
		name: name,
		db:   db,
	}
}

type Table struct {
	Name        string
	Columns     []*Column
	Keys        []*Key
	ForeignKeys []*ForeignKey
}

type Column struct {
	Name             string
	TypeName         string
	UnderlyingType   string
	IsArray          bool
	ArrayElementType string
	IsEnum           bool
	EnumValues       []string
	Nullable         bool
	DefaultValue     string
	SQLType          SQLType
}

type Key struct {
	Name    string
	IsPK    bool
	Columns []string
}

type ForeignKey struct {
	Name        string
	ToSchema    string
	ToTable     string
	FromColumns []string
	ToColumns   []string
}

type RsState struct {
	where *goqu.SelectDataset
	rows  *sql.Rows
}

type RS struct {
	table   *Table
	current *RsState
}

func NewRS(table *Table) *RS {
	return &RS{
		table:   table,
		current: &RsState{},
	}
}
