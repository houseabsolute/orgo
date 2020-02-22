package base

import (
	"database/sql"
	"fmt"
	"strings"

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
	GoName      string
	Columns     []*Column
	Keys        []*Key
	ForeignKeys []*ForeignKey
}

func (t *Table) ToCode() string {
	cols := "{\n"
	for _, c := range t.Columns {
		cols += c.ToCode() + ",\n"
	}
	cols += "}\n"

	keys := "{\n"
	for _, k := range t.Keys {
		keys += k.ToCode() + ",\n"
	}
	keys += "}\n"

	fks := "{\n"
	for _, fk := range t.ForeignKeys {
		fks += fk.ToCode() + ",\n"
	}
	fks += "}\n"

	tpl := `&base.Table{
    Name:        %q,
    Columns:     %s,
    Keys:        %s,
    ForeignKeys: %s,
}
`

	return fmt.Sprintf(tpl, t.Name, cols, keys, fks)
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
	//	SQLType          sqlType
}

func (c *Column) ToCode() string {
	tpl := `&base.Column{
    Name:             %q,
    TypeName:         %q,
    UnderlyingType:   %q,
    IsArray:          %v,
    ArrayElementType: %q,
    IsEnum,           %v,
    EnumValues,       %s,
    Nullable:         %v,
    Defaultvalue:     %q,
}`

	return fmt.Sprintf(
		tpl,
		c.Name,
		c.TypeName,
		c.UnderlyingType,
		c.IsArray,
		c.ArrayElementType,
		c.IsEnum,
		stringSliceToCode(c.EnumValues),
		c.Nullable,
		c.DefaultValue,
	)
}

type Key struct {
	Name    string
	IsPK    bool
	Columns []string
}

func (k *Key) ToCode() string {
	tpl := `&base.Key{
    Name:    %q,
    IsPK:    %v,
    Columns: %s,
}`

	return fmt.Sprintf(
		tpl,
		k.Name,
		k.IsPK,
		stringSliceToCode(k.Columns),
	)
}

type ForeignKey struct {
	Name        string
	ToSchema    string
	ToTable     string
	FromColumns []string
	ToColumns   []string
}

func (fk *ForeignKey) ToCode() string {
	tpl := `&base.ForeignKey{
    Name:        %q,
    ToSchema:    %q,
    ToTable:     %q,
    FromColumns: %s,
    ToColumns:   %s,
}`

	return fmt.Sprintf(
		tpl,
		fk.Name,
		fk.ToSchema,
		fk.ToTable,
		stringSliceToCode(fk.FromColumns),
		stringSliceToCode(fk.ToColumns),
	)
}

func stringSliceToCode(s []string) string {
	if len(s) == 0 {
		return "nil"
	}

	if len(s) == 1 {
		return "[" + fmt.Sprintf("%q", s[0]) + "]"
	}

	var quoted []string
	for _, v := range s {
		quoted = append(quoted, fmt.Sprintf("%q,\n", v))
	}

	return "[\n" + strings.Join(quoted, "") + "]"
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
