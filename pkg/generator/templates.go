package generator

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/houseabsolute/orgo/pkg/generator/templatebin"
	"golang.org/x/tools/imports"
)

func (g *Generator) generateCodeForSchema(s *schema) error {
	tpl := templatebin.MustAsset("../../templates/schema.tpl")

	parsed, err := template.New("schema").Parse(string(tpl))
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = parsed.Execute(&b, s)
	if err != nil {
		return err
	}

	path := filepath.Join(g.in, "schema", s.GoPkg, s.GoPkg+".go")
	return g.tidyAndWrite(&b, path)
}

func (g *Generator) generateCodeForRS(t *table) error {
	tpl := templatebin.MustAsset("../../templates/rs.tpl")

	parsed, err := template.New("rs").Parse(string(tpl))
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = parsed.Execute(&b, t)
	if err != nil {
		return err
	}

	path := filepath.Join(g.in, "rs", t.GoPkg, t.GoPkg+".go")
	return g.tidyAndWrite(&b, path)
}

var options = &imports.Options{
	TabWidth:  4,
	TabIndent: true,
	Comments:  true,
	Fragment:  false,
}

func (g *Generator) tidyAndWrite(b *bytes.Buffer, path string) error {
	raw := b.Bytes()
	res, err := imports.Process(path, raw, options)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	log.Print(path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(res)
	if err != nil {
		return err
	}

	return nil
}
