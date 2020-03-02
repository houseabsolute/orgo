package main

//go:generate go-bindata -nometadata -pkg templatebin -o templatebin/bindata.go templates

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/alecthomas/kingpin"
	"github.com/houseabsolute/orgo/cmd/orgo-gen/templatebin"
	"github.com/houseabsolute/orgo/pkg/base"
	"github.com/houseabsolute/orgo/pkg/strvar"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/tools/imports"
)

type generator struct {
	app                            *kingpin.Application
	database, host, user, password string
	schema                         string
	pkg                            string
	in                             string
	variator                       *strvar.Variator
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

	pkg := app.Flag("pkg", "The root name of the package in which code is being generated").
		Required().
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
	g.pkg = *pkg
	g.in = *in
	// Eventually this needs to be user configurable
	g.variator = strvar.NewWithDefaults()

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

type schema struct {
	SchemaName string
	GoName     string
	Tables     []*table
}

type table struct {
	base.Table
	Columns     []*column
	Keys        []*key
	ForeignKeys []*foreignKey
	GoName      string
	GoPkg       string
}

type column struct {
	base.Column
	GoName        string
	PrivateGoName string
}

type key struct {
	base.Key
}

type foreignKey struct {
	base.ForeignKey
}

func (g *generator) run() {
	s := &schema{
		SchemaName: g.schema,
		GoName:     g.variator.UpperCamelCase(g.schema),
		Tables:     g.tables(),
	}
	g.generateCodeForSchema(s)
	for _, t := range s.Tables {
		g.generateCodeForRS(t)
	}
}

func (g *generator) tables() []*table {
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

	var tables []*table
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			g.app.Fatalf("%s\n", err)
		}
		tables = append(tables, g.makeTable(name))
	}

	return tables
}

func (g *generator) makeTable(name string) *table {
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

	var columns []*column
	for rows.Next() {
		c := &column{}
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
			c.UnderlyingType = g.getUnderlyingType(domain)
		} else if udt != "" {
			c.TypeName = udt
			c.UnderlyingType = udt
		}

		if c.IsEnum {
			g.setEnumValues(c)
		}

		c.SQLType = base.TypeFor(c.UnderlyingType, c.Nullable)
		c.GoName = g.variator.UpperCamelCase(c.Name)
		c.PrivateGoName = g.variator.LowerCamelCase(c.Name)

		columns = append(columns, c)
	}

	table := &table{
		Table: base.Table{
			Name: name,
		},
		Columns: columns,
		GoName:  g.variator.UpperCamelCase(name),
		GoPkg:   g.variator.GoPackageName(name),
	}
	g.setKeys(table)
	g.setForeignKeys(table)

	return table
}

func (g *generator) getUnderlyingType(domain string) string {
	// This query recursively gets the domain -> udt name relationship. We
	// then just return the last row it finds, since that's the actual
	// underlying type of the columns.
	s := `
WITH RECURSIVE rdomains AS (
    SELECT udt_name, udt_schema, 1 AS n
      FROM information_schema.domains
     WHERE domain_name = $1
       AND domain_schema = $2
     UNION
    SELECT d.udt_name, d.udt_schema, rdomains.n + 1
      FROM information_schema.domains AS d
           JOIN rdomains ON (
               d.domain_name = rdomains.udt_name
               AND d.domain_schema = rdomains.udt_schema
           )
)
SELECT udt_name
  FROM rdomains
ORDER BY n DESC
LIMIT 1
`

	row := g.db.QueryRow(s, domain, g.schema)

	var name string
	err := row.Scan(&name)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	return name
}

func (g *generator) setEnumValues(c *column) {
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

func (g *generator) setKeys(t *table) {
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

	var keys []*key
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

func (g *generator) makeKey(name string, isPrimary bool, tableID, colNums string) *key {
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

	return &key{
		Key: base.Key{
			Name:    name,
			IsPK:    isPrimary,
			Columns: cols,
		},
	}
}

func (g *generator) setForeignKeys(t *table) {
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
		foreignKey := &foreignKey{
			base.ForeignKey{
				Name:     name,
				ToSchema: info[0][1],
				ToTable:  info[0][2],
			},
		}
		for _, i := range info {
			foreignKey.FromColumns = append(foreignKey.FromColumns, i[0])
			foreignKey.ToColumns = append(foreignKey.ToColumns, i[3])
		}

		t.ForeignKeys = append(t.ForeignKeys, foreignKey)
	}
}

func (g *generator) generateCodeForSchema(s *schema) {
	tpl := templatebin.MustAsset("templates/schema.tpl")

	parsed, err := template.New("schema").Parse(string(tpl))
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	var b bytes.Buffer
	err = parsed.Execute(&b, s)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	g.tidyAndPrint(&b)
}

var options = &imports.Options{
	TabWidth:  4,
	TabIndent: true,
	Comments:  true,
	Fragment:  false,
}

func (g *generator) generateCodeForRS(t *table) {
	tpl := templatebin.MustAsset("templates/rs.tpl")

	parsed, err := template.New("rs").Parse(string(tpl))
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	var b bytes.Buffer
	err = parsed.Execute(&b, t)
	if err != nil {
		g.app.Fatalf("%s\n", err)
	}

	g.tidyAndPrint(&b)
}

func (g *generator) tidyAndPrint(b *bytes.Buffer) {
	raw := b.Bytes()
	res, err := imports.Process("dummy.go", raw, options)
	if err != nil {
		g.app.Fatalf("\n%s\n%s\n", string(raw), err)
	}

	os.Stdout.Write(res)
}

func (t *table) ToCode() string {
	cols := "{\n"
	for _, c := range t.Columns {
		cols += c.ToCode() + ",\n"
	}
	cols += "}"

	keys := "{\n"
	for _, k := range t.Keys {
		keys += k.ToCode() + ",\n"
	}
	keys += "}"

	fks := "{\n"
	for _, fk := range t.ForeignKeys {
		fks += fk.ToCode() + ",\n"
	}
	fks += "}"

	tpl := `&base.Table{
    Name:        %q,
    Columns:     %s,
    Keys:        %s,
    ForeignKeys: %s,
}
`

	return fmt.Sprintf(tpl, t.Name, cols, keys, fks)
}

func (t *table) GoPkgShortName() string {
	parts := strings.Split(t.GoPkg, "/")
	return parts[len(parts)-1]
}

func (c *column) ToCode() string {
	tpl := `&base.Column{
    Name:             %q,
    TypeName:         %q,
    UnderlyingType:   %q,
    IsArray:          %v,
    ArrayElementType: %q,
    IsEnum:           %v,
    EnumValues:       %s,
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

func (k *key) ToCode() string {
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

func (fk *foreignKey) ToCode() string {
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
		return "{" + fmt.Sprintf("%q", s[0]) + "}"
	}

	var quoted []string
	for _, v := range s {
		quoted = append(quoted, fmt.Sprintf("%q,\n", v))
	}

	return "{\n" + strings.Join(quoted, "") + "}"
}
