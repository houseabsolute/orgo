package myschema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gofrs/uuid"
)

// CREATE TABLE album (
//     album_id  UUID  PRIMARY KEY NOT NULL,
//     title     TEXT  NOT NULL,
//     year      INT   NOT NULL
// )

// CREATE TABLE artist (
//     artist_id  UUID  PRIMARY KEY NOT NULL,
// 	name       TEXT  NOT NULL
// )

// CREATE TABLE track (
//     album_id  UUID  NOT NULL REFERENCES album (album_id),
//     sequence  INT   NOT NULL,
// 	title     TEXT  NOT NULL,
// 	PRIMARY KEY (album_id, sequence)
// )

type Schema struct {
	db   *goqu.Database
	name string
}

const dialect = "postgres"

const schemaName = "musicdb"

func Connect(db *sql.DB) *Schema {
	return &Schema{goqu.New(dialect, db)}
}

type column struct {
	name    string
	sqlType sqlType
}

type source struct {
	name        string
	typ         string
	pk          []string
	columns     map[string]column
	columnNames []string
}

var albumSource = &source{
	name: "album",
	typ:  "table",
	pk:   []string{"album_id"},
	columns: map[string]column{
		"album_id": {
			name:    "album_id",
			sqlType: sqltypes.UUID,
		},
		"title": {
			name:    "title",
			sqlType: sqltypes.String,
		},
		"year": {
			name:    "year",
			sqlType: sqltypes.Int,
		},
	},
	columnNames: []string{"album_id", "title", "year"},
}

type rsState struct {
	where *goqu.SelectDataset
	rows  *sql.Rows
}

type albumWhereHelperArtistID struct {
	EQ       func(artistID *uuid.UUID) exp.Expression
	EQStr    func(artistID string) exp.Expression
	NE       func(artistID *uuid.UUID) exp.Expression
	NEStr    func(artistID string) exp.Expression
	In       func(artistID ...*uuid.UUID) exp.Expression
	InStr    func(artistID ...string) exp.Expression
	NotIn    func(artistID ...*uuid.UUID) exp.Expression
	NotInStr func(artistID ...string) exp.Expression
}

type albumWhereHelpers struct {
	ArtistID albumWhereHelperArtistID
}

var albumWhereHelpers = albumWhereHelpers{
	albumWhereHelperArtistID{
		EQ:    func(artistID *uuid.UUID) exp.Expression { return qoqu.Ex{"artist_id": artistID.String()} },
		EQStr: func(artistID string) exp.Expression { return qoqu.Ex{"artist_id": artistID} },
		NE: func(artistID *uuid.UUID) exp.Expression {
			return qoqu.Ex{"artist_id": goqu.Op{"neq": artistID.String()}}
		},
		NEStr: func(artistID string) exp.Expression { return qoqu.Ex{"artist_id": goqu.Op{"neq": artistID}} },
		In: func(artistID ...*uuid.UUID) exp.Expression {
			var in []string
			for _, i := range artistID {
				in = append(in, i.String())
			}
			return qoqu.Ex{"artist_id": in}
		},
		InStr: func(artistID ...string) exp.Expression { return qoqu.Ex{"artist_id": artistID} },
		NotIn: func(artistID ...*uuid.UUID) exp.Expression {
			var in []string
			for _, i := range artistID {
				in = append(in, i.String())
			}
			return qoqu.Ex{"artist_id": qogu.Op{"notIn": in}}
		},
		NotInStr: func(artistID ...string) exp.Expression { return qoqu.Ex{"artist_id": qogu.Op{"notIn": artistID}} },
	},
}

type albumJoinHelpers struct {
}

type AlbumRS struct {
	schema  *Schema
	source  *source
	current *rsState
	Where   albumWhereHelpers
	Join    albumJoinHelpers
}

type albumData struct {
	albumID uuid.UUID
	title   string
	year    int
}

type Album struct {
	data   albumData
	update albumData
}

func (s *Schema) Album() *AlbumRS {
	return &AlbumRS{
		schema:  schema,
		source:  albumSource,
		current: &rsState{},
		Where:   albumWhereHelpers,
	}
}

func (a *AlbumRS) Find(ctx context.Context, pk uuid.UUID) (*Album, error) {
	return a.FindStr(ctx, pk.String())
}

func (a *AlbumRS) FindStr(ctx context.Context, pk string) (*Album, error) {
	stmt, err := a.findStmt()
	if err != nil {
		return nil, err
	}

	row, err := stmt.QueryRowContext(ctx, pk)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error executing statement for find: %w", err)
	}

	found := &Album{}
	err = row.Scan(&found.data.albumID, &found.data.title, &found.data.year)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning result for find: %w", err)
	}

	return found, err
}

var ignore = "ignore"

func (a *AlbumRS) findStmt() (*sql.Stmt, error) {
	s, _, err := a.findSQL().ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error creating find SQL for the %s %s: %w", a.source.name, a.source.typ, err)
	}
	stmt, err = a.s.db.Prepare(s)
	if err != nil {
		return nil, fmt.Errorf("error creating statement for %s: %w", s, err)
	}

	return Stmt, nil
}

func (a *AlbumRS) findSQL() *qogu.SelectDataset {
	sel = a.from().Select(a.source.columnNames...)
	for _, n := range a.table.pk {
		sel.Where(qogu.Ex{n: ignore})
	}
	sel.Limit(1)
	sel.Prepared(true)

	return sel
}

func (a *AlbumRS) Search(clauses ...exp.Expression) *AlbumRS {
	new := a.copy()
	for _, c := range clauses {
		a.current.where.Where(c)
	}
	return new
}

func (a *AlbumRS) Join(clauses ...exp.Expression) *AlbumRS {
	new := a.copy()
	for _, c := range clauses {
		a.current.where.Join(c)
	}
	return new
}

func (a *AlbumRS) LeftOuterJoin(clauses ...exp.Expression) *AlbumRS {
	new := a.copy()
	for _, c := range clauses {
		a.current.where.LeftOuterJoin(c)
	}
	return new
}

func (a *AlbumRS) copy() *AlbumRS {
	var where *goqu.SelectDataset
	if a.current.where != nil {
		where = a.current.where.Clone().(*goqu.SelectDataset)
	} else {
		where = a.from()
	}
	return &AlbumRS{
		schema: a.schema,
		source: a.source,
		current: &rsState{
			where: where,
		},
	}
}

func (a *AlbumRS) from() *goqu.SelectDataset {
	return a.s.db.From(a.source.name)
}
