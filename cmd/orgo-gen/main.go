package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/houseabsolute/orgo/pkg/generator"
	_ "github.com/lib/pq"
)

type app struct {
	kp        *kingpin.Application
	generator *generator.Generator
}

func main() {
	app, err := new()
	if err != nil {
		app.kp.FatalUsage("%s\n", err)
	}

	err = app.generator.Generate()
	if err != nil {
		app.kp.FatalUsage("%s\n", err)
	}
}

func new() (*app, error) {
	kp := kingpin.New(
		"orgo-gen",
		"Generate orgo code for your schema.",
	).
		Author("Dave Rolsky <autarch@urth.org>").
		Version("0.0.1").
		UsageWriter(os.Stdout)
	kp.HelpFlag.Short('h')

	database := kp.Flag("database", "Name of database to connect to (or PGDATABASE env var)").
		Required().
		Envar("PGDATABASE").
		String()

	host := kp.Flag("host", "Database host (or PGHOST env var)").
		Default("localhost").
		Envar("PGHOST").
		String()

	user := kp.Flag("user", "User with which to connecto to database (or PGUSER env var)").
		Required().
		Envar("PGUSER").
		String()

	password := kp.Flag("password", "Password with which to connect to database (or PGPASSWORD env var)").
		Envar("PGPASSWORD").
		String()

	schema := kp.Flag("schema", "Name of the schema from which to generate code").
		Default("public").
		String()

	pkg := kp.Flag("pkg", "The root name of the package in which code is being generated").
		Required().
		String()

	in := kp.Flag("in", "Directory in which to generate code").
		Required().
		String()

	app := &app{
		kp: kp,
	}

	_, err := kp.Parse(os.Args[1:])

	db := app.connect(*database, *host, *user, *password)
	app.generator = generator.New(db, *schema, *pkg, *in)

	return app, err
}

func (app *app) connect(database, host, user, password string) *sql.DB {
	conn := fmt.Sprintf(
		"dbname=%s host=%s user=%s",
		database, host, user,
	)
	if password != "" {
		conn += fmt.Sprintf(" password=%s", password)
	}

	db, err := sql.Open("postgres", conn)
	if err != nil {
		app.kp.Fatalf("%s\n", err)
	}

	return db
}
