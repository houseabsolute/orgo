package {{ .GoPkgShortName }}

import (
	"github.com/houseabsolute/orgo/pkg/base"
	"github.com/doug-martin/goqu/v9/exp"
)

var t = {{ .ToCode }}

{{- $rs := .GoName | printf "%sRS" }}
type {{ $rs }} struct {
	schema	*Schema
	rs		*base.RS
}

func New() *{{ $rs }} {
	return {{ $rs }}{
		schema: s,
		rs:		base.NewRS(t),
	}
}

{{ range .Columns }}
type {{ .PrivateGoName }}WhereHelpers struct {
	{{- if .SQLType.CanEQ }}
	EQ func( {{ .SQLType.GoType }} ) exp.Expression
	NE func( {{ .SQLType.GoType }} ) exp.Expression
	In func( []{{ .SQLType.GoType }} ) exp.Expression
	NotIn func( []{{ .SQLType.GoType }} ) exp.Expression
		{{- if .SQLType.CanString }}
	EQStr func(string) exp.Expression
	NEStr func(string) exp.Expression
	InStr func([]string) exp.Expression
	NotInStr func([]string) exp.Expression
		{{- end }}
	{{- end }}
	{{- if .SQLType.CanLTGT }}
	LT func( {{ .SQLType.GoType }} ) exp.Expression
	GT func( {{ .SQLType.GoType }} ) exp.Expression
	LTE func( {{ .SQLType.GoType }} ) exp.Expression
	GTE func( {{ .SQLType.GoType }} ) exp.Expression
		{{- if .SQLType.CanString }}
	LTStr func(string) exp.Expression
	GTStr func(string) exp.Expression
	LTEStr func(string) exp.Expression
	GTEStr func(string) exp.Expression
		{{- end }}
	{{- end }}
	{{- if .SQLType.CanIsNull }}
	IsNull func() exp.Expression
	IsNotNull func() exp.Expression
	{{- end }}
}

{{ end -}}
type WhereHelpers struct {
{{- range .Columns }}
	{{.GoName}} {{.PrivateGoName}}WhereHelpers
{{- end }}
}

var Where = WhereHelpers{
{{- range .Columns -}}
{{- $qname := .Name | printf "%q" }}
{{- $privpl := .PrivateGoName | printf "%ss" }}
	{{ .GoName }}: {{.PrivateGoName}}WhereHelpers{
		{{- if .SQLType.CanEQ -}}
			{{- if .SQLType.CanString }}
		EQ: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, {{ .PrivateGoName }}.String()}
		},
		NE: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.NeqOp: {{ .PrivateGoName }}.String() } }
		},
		In: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			var in []string
			for _, i := range {{ $privpl }} {
				in = append(in, i.String() )
			}
			return qoqu.Ex{ {{ $qname }}: in }
		},
		NotIn: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			var in []string
			for _, i := range {{ $privpl }} {
				in = append(in, i.String() )
			}
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ exp.NotInOp: in } }
		},
		EQStr: func( {{ .PrivateGoName }} string) exp.Expression {
			return qoqu.Ex{ {{$qname}}: {{ .PrivateGoName }} }
		},
		NEStr: func( {{ .PrivateGoName }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.NeqOp: {{ .PrivateGoName }} } }
		},
		InStr: func( {{ $privpl }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: {{ $privpl }} }
		},
		NotInStr: func( {{ $privpl }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ exp.NotInOp: {{ $privpl }} } }
		},
			{{- else }}
		EQ: func( {{.PrivateGoName}} {{.SQLType.GoType}} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, {{.PrivateGoName}} }
		},
		NE: func( {{.PrivateGoName}} {{.SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.NeqOp: {{.PrivateGoName}} } }
		},
		In: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: {{ $privpl }} }
		},
		NotIn: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ exp.NotInOp: {{ $privpl }} } }
		},
			{{- end -}}
		{{- end }}
		{{- if .SQLType.CanLTGT -}}
			{{- if .SQLType.CanString }}
		LT: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LtOp: {{ .PrivateGoName }}.String() } }
		},
		GT: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GtOp: {{ .PrivateGoName }}.String() } }
		},
		LTE: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LteOp: {{ .PrivateGoName }}.String() } }
		},
		GTE: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GteOp: {{ .PrivateGoName }}.String() } }
		},
		LTStr: func( {{ .PrivateGoName }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LtOp: {{ .PrivateGoName }} } }
		},
		GTStr: func( {{ .PrivateGoName }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GtOp: {{ .PrivateGoName }} } }
		},
		LTEStr: func( {{ .PrivateGoName }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LteOp: {{ .PrivateGoName }} } }
		},
		GTEStr: func( {{ .PrivateGoName }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GteOp: {{ .PrivateGoName }} } }
		},
			{{ else }}
		LT: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LtOp: {{ .PrivateGoName }} } }
		},
		GT: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GtOp: {{ .PrivateGoName }} } }
		},
		LTE: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.LteOp: {{ .PrivateGoName }} } }
		},
		GTE: func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.GteOp: {{ .PrivateGoName }} } }
		},
			{{- end -}}
		{{- end }}
		{{- if .SQLType.CanIsNull }}
		IsNull: func() exp.expression {
			return qoqu.Ex{ {{ $qname }}, nil }
		},
		IsNull: func() exp.expression {
			return qoqu.Ex{ {{ $qname }}, goqu.Op{ exp.NeOp, nil } }
		},
	{{- end }}
	},
{{- end }}
}
