package {{ .SchemaName }}

import (
    "github.com/houseabsolute/pkg/base"
    {{- range .Tables }}
    {{ .GoPkg | printf "%q" }}
    {{- end }}
)

type Schema struct {
    dialect string
    *base.Schema
}

{{- range .Tables }}
{{- $qname := .Name | printf "%q" }}
{{- $rs := .GoName | printf "%sRS" }}

func (s *Schema) {{ $rs }}() *{{ $rs }} {
    return {{ .GoPkgShortName }}.NewRS(s)
}
{{- end }}
