package base

import (
	"strings"
)

type SQLType struct {
	GoType    string
	Imports   []string
	CanIsNull bool
	CanEQ     bool
	CanLTGT   bool
	CanString bool
	Ops       map[string]string
}

var typeMap map[string]SQLType
var nullTypeMap map[string]SQLType
var arrayTypeMap map[string]SQLType

func init() {
	typeMap = map[string]SQLType{
		"bigint": SQLType{
			GoType:    "Int64",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"integer": SQLType{
			GoType:    "Int",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"int64": SQLType{
			GoType:    "Int16",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"decimal": SQLType{
			GoType:    "types.Decimal",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"double precision": SQLType{
			GoType:    "float64",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"float32": SQLType{
			GoType:    "float32",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"text": SQLType{
			GoType:    "string",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"char": {
			GoType:    "types.Byte",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"bytea": {
			GoType:    "[]byte",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"boolean": {
			GoType:    "bool",
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		// It'd be nice to separate out date, time, & timestamp types. The
		// builtin time.Time has too much information for Pg's date & time
		// types, which could lead to coding mistakes. For example, someone
		// might thank that the date associated with a bare time has a
		// meaning (or vice versa with a date type).
		"date": {
			GoType:    "time.Time",
			Imports:   []string{"time"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: true,
		},
		"json": {
			GoType:    "types.JSON",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: true,
		},
		"timestamptz": {
			GoType:    "time.Time",
			Imports:   []string{"time"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"uuid": {
			GoType:    "uuid.UUID",
			Imports:   []string{"github.com/gofrs/uuid"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: true,
		},
		"hstore": {
			GoType:    "types.HStore",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   true,
			CanString: false,
		},
		"point": {
			GoType:    "pgeo.Point",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"line": {
			GoType:    "pgeo.Line",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"lseg": {
			GoType:    "pgeo.Lseg",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"box": {
			GoType:    "pgeo.Box",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"path": {
			GoType:    "pgeo.Path",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"polygon": {
			GoType:    "pgeo.Polygon",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		"circle": {
			GoType:    "pgeo.Circle",
			Imports:   []string{"github.com/volatiletech/sqlboiler/types/pgeo"},
			CanIsNull: false,
			CanEQ:     true,
			CanLTGT:   false,
			CanString: false,
		},
		// Need good mappings for these!
		// "interval":    {},
		// "bit":         {},
		// "bit varying": {},
		// "cidr":        {},
		// "inet":        {},
		// "macaddr":     {},
		// "macaddr8":    {},
	}

	aliases := map[string][]string{
		"bigint":  []string{"bigserial"},
		"decimal": []string{"numeric"},
		"int":     []string{"serial"},
		"json":    []string{"jsonb"},
		"numeric": []string{"money"},
		// Should XML map to something else? If there was a type that _lazily_
		// parsed the XML that could be nice, although what if the XML fails
		// to parse? Then every method needs to return an error.
		"text": []string{
			"character",
			"character varying",
			"citext",
			"xml",
		},
		"time": []string{
			"time with time zone",
			"time without time zone",
			"timestamp with time zone",
			"timestamp without time zone",
		},
	}

	for k, v := range aliases {
		for _, t := range v {
			typeMap[t] = typeMap[k]
		}
	}

	nullTypeMap = make(map[string]SQLType, len(typeMap))
	for n, t := range typeMap {
		c := t.clone()
		c.CanIsNull = true
		if c.GoType == "types.Decimal" {
			c.GoType = "types.NullDecimal"
		} else if c.GoType == "types.HStore" {
			//			panic("No null hstore type!")
		} else if strings.HasPrefix(c.GoType, "types.") {
			c.GoType = strings.Replace(c.GoType, "types.", "null.", 1)
			c.Imports = []string{"github.com/volatiletech/null"}
		} else if strings.HasPrefix(c.GoType, "pgeo.") {
			c.GoType = strings.Replace(c.GoType, "pgeo.", "pgeo.Null", 1)
		} else if c.GoType == "[]byte" {
			c.GoType = "null.Bytes"
			c.Imports = []string{"github.com/volatiletech/null"}
		} else if c.GoType == "uuid.UUID" {
			c.GoType = "uuid.NullUUID"
		} else {
			c.GoType = "null." + strings.Title(strings.ToLower(c.GoType))
		}

		nullTypeMap[n] = c
	}

	// XXX - need to also handle arrays, composite types, and range types
}

// TypeFor takes a database type name and returns the right SQLType for that
// data type.
func TypeFor(t string, null bool) SQLType {
	if null {
		return nullTypeMap[t]
	}

	return typeMap[t]
}

func (t SQLType) clone() SQLType {
	return SQLType{
		GoType:    t.GoType,
		Imports:   t.Imports,
		CanIsNull: t.CanIsNull,
		CanEQ:     t.CanEQ,
		CanLTGT:   t.CanLTGT,
		CanString: t.CanString,
		Ops:       t.Ops,
	}
}
