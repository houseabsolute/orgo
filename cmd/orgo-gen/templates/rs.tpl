package {{ .GoPkgShortName }}

import "github.com/houseabsolute/orgo/pkg/base"

var t = {{ .ToCode }}

{{- $rs := .GoName | printf "%sRS" }}
type {{ $rs }} struct {
    schema  *Schema
    rs      *base.RS
}

func New() *{{ $rs }} {
    return {{ $rs }}{
        schema: s,
        rs:     base.NewRS(t),
    }
}

{{- range .Columns }}
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
        {{- if .SQLType.CanString }}
    LTStr func(string) exp.Expression
        {{- end }}
    GT func( {{ .SQLType.GoType }} ) exp.Expression
        {{- if .SQLType.CanString }}
    GTStr func(string) exp.Expression
        {{- end }}
    {{- end }}
    {{- if .SQLType.CanIsNull }}
    IsNull func() exp.Expression
    IsNotNull func() exp.Expression
    {{- end }}
}
{{- end }}

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
		EQ:    func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
            return qoqu.Ex{ {{ $qname }}, {{ .PrivateGoName }}.String()}
        },
		NE:    func( {{ .PrivateGoName }} {{ .SQLType.GoType }} ) exp.Expression {
            return qoqu.Ex{ {{ $qname }}, goqu.Op{"neq": {{ .PrivateGoName }}.String() } }
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
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ "notIn": in } }
		},
		EQStr: func( {{ .PrivateGoName }} string) exp.Expression {
            return qoqu.Ex{ {{$qname}}: {{ .PrivateGoName }} }
        },
		NEStr: func( {{ .PrivateGoName }} string ) exp.Expression {
            return qoqu.Ex{ {{ $qname }}, goqu.Op{"neq": {{ .PrivateGoName }} } }
        },
		InStr: func( {{ $privpl }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: {{ $privpl }} }
		},
		NotInStr: func( {{ $privpl }} string ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ "notIn": {{ $privpl }} } }
		},
            {{ else }}
		EQ:    func( {{.PrivateGoName}} {{.SQLType.GoType}} ) exp.Expression {
            return qoqu.Ex{ {{ $qname }}, {{.PrivateGoName}} }
        },
		NE:    func( {{.PrivateGoName}} {{.SQLType.GoType }} ) exp.Expression {
            return qoqu.Ex{ {{ $qname }}, goqu.Op{"neq": {{.PrivateGoName}} } }
        },
		In: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: {{ $privpl }} }
		},
		NotIn: func( {{ $privpl }} {{ .SQLType.GoType }} ) exp.Expression {
			return qoqu.Ex{ {{ $qname }}: qoqu.Op{ "notIn": {{ $privpl }} } }
		},
            {{- end -}}
        {{- end }}
	},
{{- end }}
}
