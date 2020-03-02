package strvar

import (
	"regexp"
	"strings"
)

type Variator struct {
	abbreviations map[string]bool
}

// From sqlboiler with additions
var abbreviations = map[string]bool{
	"acl":   true,
	"api":   true,
	"ascii": true,
	"cpu":   true,
	"eof":   true,
	"gpu":   true,
	"guid":  true,
	"id":    true,
	"ip":    true,
	"json":  true,
	"ram":   true,
	"sla":   true,
	"udp":   true,
	"ui":    true,
	"uid":   true,
	"uuid":  true,
	"uri":   true,
	"url":   true,
	"utf8":  true,
	"xml":   true,
}

func New(extra []string) *Variator {
	a := make(map[string]bool, len(abbreviations)+len(extra))
	for k, v := range abbreviations {
		a[k] = v
	}
	for _, k := range extra {
		a[k] = true
	}
	return &Variator{a}
}

func NewWithDefaults() *Variator {
	return New([]string{})
}

func (v *Variator) UpperCamelCase(n string) string {
	parts := v.parts(n)

	final := ""
	for _, p := range parts {
		if v.abbreviations[p] {
			final += strings.ToUpper(p)
		} else {
			final += strings.Title(strings.ToLower(p))
		}
	}

	return final
}

func (v *Variator) LowerCamelCase(n string) string {
	parts := v.parts(n)

	final := ""
	for i, p := range parts {
		if i == 0 {
			final += strings.ToLower(p)
			continue
		}

		if v.abbreviations[p] {
			final += strings.ToUpper(p)
		} else {
			final += strings.Title(strings.ToLower(p))
		}
	}

	return final
}

func (v *Variator) GoPackageName(n string) string {
	parts := v.parts(n)
	return strings.Join(parts, "")
}

var splitter = regexp.MustCompile(`_+`)

func (v *Variator) parts(n string) []string {
	parts := splitter.Split(strings.Trim(n, "_"), -1)
	if len(parts) == 0 {
		return []string{n}
	}
	return parts
}
