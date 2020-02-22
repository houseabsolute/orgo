package core

type Schema struct {
	name   string
	tables []Table
}

type Table struct {
	name        string
	pk          []string
	columnNames []string
	columns     map[string]Column
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) PK() []string {
	return t.pk
}

func (t *Table) ColumnNames() []Column {
	return t.columnNames
}

func (t *Table) Column(name string) Column {
	

func (t *Table) Columns() []Column {
	return t.columns
}

type Column struct {
	name    string
	sqlType SQLType
}
