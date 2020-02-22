package main

//go:generate go-bindata -nometadata -pkg templatebin -o templatebin/bindata.go templates

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/alecthomas/kingpin"
	"github.com/houseabsolute/orgo/cmd/orgo-gen/templatebin"
	"github.com/houseabsolute/orgo/pkg/base"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stoewer/go-strcase"
)

type generator struct {
	app                            *kingpin.Application
	database, host, user, password string
	schema                         string
	in                             string
	db                             *sql.DB
	enums                          map[string][]string
}

func main() {
	g, err := new()
	if err != nil {
		g.app.FatalUsage("%s\n", err)
	}

	g.connect()
	g.run()
}

func new() (*generator, error) {
	app := kingpin.New(
		"orgo-gen",
		"Generate orgo code for your schema.",
	).
		Author("Dave Rolsky <autarch@urth.org>").
		Version("0.0.1").
		UsageWriter(os.Stdout)
	app.HelpFlag.Short('h')

	database := app.Flag("database", "Name of database to connect to (or PGDATABASE env var)").
		Required().
		Envar("PGDATABASE").
		String()

	host := app.Flag("host", "Database host (or PGHOST env var)").
		Default("localhost").
		Envar("PGHOST").
		String()

	user := app.Flag("user", "User with which to connecto to database (or PGUSER env var)").
		Required().
		Envar("PGUSER").
		String()

	password := app.Flag("password", "Password with which to connect to database (or PGPASSWORD env var)").
		Envar("PGPASSWORD").
		String()

	schema := app.Flag("schema", "Name of the schema from which to generate code").
		Default("public").
		String()

	in := app.Flag("in", "Directory in which to generate code").
		Required().
		String()

	g := &generator{
		app:   app,
		enums: make(map[string][]string),
	}

	_, err := app.Parse(os.Args[1:])

	g.database = *database
	g.host = *host
	g.user = *user
	g.password = *password
	g.schema = *schema
	g.in = *in

	return g, err
}

func (g *generator) connect() {
	conn := fmt.Sprintf(
		"dbname=%s host=%s user=%s",
		g.database, g.host, g.user,
	)
	if g.password != "" {
		conn += fmt.Sprintf(" password=%s", g.password)
	}
	var err error
	g.db, err = sql.Open("postgres", conn)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
}

func (g *generator) run() {
	g.generateCodeForSchema()
}

func (g *generator) tables() []*base.Table {
	s := `
SELECT table_name
  FROM information_schema.tables
 WHERE table_schema = $1
   AND table_type = 'BASE TABLE'
ORDER BY table_name ASC
`

	rows, err := g.db.Query(s, g.schema)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	var tables []*base.Table
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			g.app.Fatalf("%s\n", err)
		}
		tables = append(tables, g.makeTable(name))
	}

	return tables
}

func (g *generator) makeTable(name string) *base.Table {
	// Borrowed from sqlboiler and then altered a fair bit.
	s := `
SELECT c.column_name,
       COALESCE( c.domain_name, '' ) AS domain_name,
       c.udt_name,
       COALESCE( e.data_type, '' ) AS element_type,
       c.data_type = 'ARRAY' AS is_array,
       t.typtype IS NOT NULL AND t.typtype = 'e' AS is_enum,  
       c.is_nullable = 'YES' AS is_nullable,
       COALESCE( c.column_default, '' ) AS column_default
  FROM information_schema.columns AS c
       JOIN pg_namespace AS n ON ( c.udt_schema = n.nspname )
       -- from sqlboiler
       LEFT OUTER JOIN pg_type AS t ON (
           c.data_type = 'USER-DEFINED'
           AND n.oid = t.typnamespace
           AND c.udt_name = t.typname
       )
       -- From https://www.postgresql.org/docs/10/infoschema-element-types.html
       LEFT OUTER JOIN information_schema.element_types AS e ON (
           ( c.table_catalog,  c.table_schema,  c.table_name, 'TABLE',        c.dtd_identifier )
		 = ( e.object_catalog, e.object_schema, e.object_name, e.object_type, e.collection_type_identifier )
       )
 WHERE c.table_schema = $1 and c.table_name = $2
`

	rows, err := g.db.Query(s, g.schema, name)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	var columns []*base.Column
	for rows.Next() {
		c := &base.Column{}
		var domain, udt string
		if err = rows.Scan(
			&c.Name, &domain, &udt, &c.ArrayElementType, &c.IsArray, &c.IsEnum, &c.Nullable, &c.DefaultValue,
		); err != nil {
			g.app.Fatalf("%s\n", err)
		}

		// For some reason the udt name for array types has a leading
		// underscore, so you get "_numeric" instead of "numeric". In those
		// cases it seems like there is no information in the
		// information_schema.element_types table for this column.
		if strings.HasPrefix(udt, "_") {
			c.ArrayElementType = udt[1:]
		}
		if domain != "" {
			c.TypeName = domain
			c.UnderlyingType = udt
		} else if udt != "" {
			c.TypeName = udt
		}

		if c.IsEnum {
			g.setEnumValues(c)
		}

		columns = append(columns, c)
	}

	table := &base.Table{
		Name:    name,
		GoName:  strcase.UpperCamelCase(name),
		Columns: columns,
	}
	g.setKeys(table)
	g.setForeignKeys(table)

	return table
}

func (g *generator) setEnumValues(c *base.Column) {
	if len(g.enums[c.TypeName]) > 0 {
		c.EnumValues = g.enums[c.TypeName]
		return
	}

	// Borrowed from sqlboiler
	s := `
SELECT e.enumlabel
  FROM pg_enum AS e
       JOIN pg_type AS t ON ( e.enumtypid = t.typelem )
 WHERE t.typname = $1
ORDER BY enumsortorder
`

	rows, err := g.db.Query(s, "_"+c.TypeName)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var v string
		if err = rows.Scan(&v); err != nil {
			g.app.Fatalf("%s\n", err)
		}
		values = append(values, v)
	}

	g.enums[c.TypeName] = values
	c.EnumValues = values
}

func (g *generator) setKeys(t *base.Table) {
	// From
	// https://metacpan.org/release/DBIx-Class-Schema-Loader/source/lib/DBIx/Class/Schema/Loader/DBI/Pg.pm#L117
	s := `
SELECT i.relname, x.indisprimary, x.indrelid, x.indkey
  FROM pg_index AS x
	   JOIN pg_catalog.pg_class AS c ON ( c.oid = x.indrelid )
	   JOIN pg_catalog.pg_class AS i ON ( i.oid = x.indexrelid )
	   JOIN pg_catalog.pg_namespace AS n ON ( n.oid = c.relnamespace)
 WHERE x.indisunique = 't'
   AND x.indpred     IS NULL
   AND c.relkind     = 'r'
   AND i.relkind     = 'i'
   AND n.nspname     = $1
   AND c.relname     = $2
ORDER BY i.relname
`

	rows, err := g.db.Query(s, g.schema, t.Name)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	var keys []*base.Key
	for rows.Next() {
		var name, tableID, colNums string
		var isPrimary bool
		if err = rows.Scan(&name, &isPrimary, &tableID, &colNums); err != nil {
			g.app.Fatalf("%s\n", err)
		}

		keys = append(keys, g.makeKey(name, isPrimary, tableID, colNums))
	}

	t.Keys = keys
}

func (g *generator) makeKey(name string, isPrimary bool, tableID, colNums string) *base.Key {
	var nums []int
	for _, i := range strings.Split(colNums, " ") {
		n, err := strconv.Atoi(i)
		if err != nil {
			g.app.Fatalf("%s\n", err)
		}
		nums = append(nums, n)
	}

	s := `
SELECT attname
  FROM pg_attribute
 WHERE attrelid = $1
   AND attnum = ANY($2)
ORDER BY attnum
`

	rows, err := g.db.Query(s, tableID, pq.Array(nums))
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var c string
		if err = rows.Scan(&c); err != nil {
			g.app.Fatalf("%s\n", err)
		}
		cols = append(cols, c)
	}

	return &base.Key{
		Name:    name,
		IsPK:    isPrimary,
		Columns: cols,
	}
}

func (g *generator) setForeignKeys(t *base.Table) {
	// From https://stackoverflow.com/a/1152321
	s := `
SELECT tc.constraint_name AS name,
	   kcu.column_name AS from_column,
	   ccu.table_schema AS to_table_schema,
	   ccu.table_name AS to_table_name,
	   ccu.column_name AS to_column_name
  FROM information_schema.table_constraints AS tc
       JOIN information_schema.key_column_usage AS kcu
           ON ( tc.constraint_name = kcu.constraint_name
		        AND tc.table_schema = kcu.table_schema )
    JOIN information_schema.constraint_column_usage AS ccu
           ON ( ccu.constraint_name = tc.constraint_name
                AND ccu.table_schema = tc.table_schema )
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = $1
  AND tc.table_name = $2
ORDER BY kcu.ordinal_position
`

	rows, err := g.db.Query(s, g.schema, t.Name)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
	defer rows.Close()

	keys := make(map[string][][]string)
	for rows.Next() {
		var name, from_column, to_schema, to_table, to_column string
		if err = rows.Scan(&name, &from_column, &to_schema, &to_table, &to_column); err != nil {
			g.app.Fatalf("%s\n", err)
		}
		keys[name] = append(keys[name], []string{from_column, to_schema, to_table, to_column})
	}

	for name, info := range keys {
		foreignKey := &base.ForeignKey{
			Name:     name,
			ToSchema: info[0][1],
			ToTable:  info[0][2],
		}
		for _, i := range info {
			foreignKey.FromColumns = append(foreignKey.FromColumns, i[0])
			foreignKey.ToColumns = append(foreignKey.ToColumns, i[3])
		}

		t.ForeignKeys = append(t.ForeignKeys, foreignKey)
	}
}

func (g *generator) generateCodeForSchema() {
	tpl := templatebin.MustAsset("templates/schema.tpl")

	parsed, err := template.New("schema").Parse(string(tpl))
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	err = parsed.Execute(
		os.Stdout,
		struct {
			SchemaName string
			Tables     []*base.Table
		}{
			SchemaName: g.schema,
			Tables:     g.tables(),
		},
	)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}
}

func (g *generator) generateCodeForTable(t *base.Table) {

}
