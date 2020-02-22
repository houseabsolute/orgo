package {{.SchemaName}}

import "github.com/houseabsolute/pkg/base"

type Schema struct {
    dialect string
    *base.Schema
}

var tables map[string]*base.Table

func init() {
    tables = map[string]*base.Table{
{{- range .Tables}}
{{- $qname := .Name | printf "%q" }}
    {{$qname}}: {{.ToCode}},
{{- end}}
    }
}
}

{{range .Tables}}
{{- $qname := .Name | printf "%q" }}
{{- $rs := .GoName | printf "%sRS" }}
type {{$rs}} struct {
    schema  *Schema
    rs      *base.RS
}

func (s *Schema) {{$rs}}() *{{$rs}} {
    return {{$rs}}{
        schema: s,
        rs:     base.NewRS(tables[{{$qname}}]),
    }
}
{{end}}
