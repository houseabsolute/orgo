package {{ .GoPkg }}

import (
    {{- range .Imports }}
    {{ . | printf "%q" }}
    {{- end }}
)

type Schema struct {
    dialect string
    *base.Schema
}

{{- range .Tables }}
{{- $qname := .Name | printf "%q" }}
{{- $rs := .GoName | printf "%sRS" }}

func (s *Schema) {{ $rs }}() *{{ .GoPkg }}.{{ $rs }} {
    return {{ .GoPkgShortName }}.NewRS(s)
}
{{- end }}
