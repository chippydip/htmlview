package htmlview

import (
	"html/template"
	"path/filepath"
)

type Template struct {
	tmpl  *template.Template
	funcs template.FuncMap
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Funcs(funcMap template.FuncMap) *Template {
	if t.tmpl != nil {
		t.tmpl.Funcs(funcMap)
	} else {
		if t.funcs == nil {
			t.funcs = make(template.FuncMap)
		}
		for k, v := range funcMap {
			t.funcs[k] = v
		}
	}
	return t
}

func (t *Template) ParseFiles(filenames ...string) *Template {
	if len(filenames) == 0 {
		return t
	}

	tmpl := t.tmpl
	if tmpl == nil {
		name := filepath.Base(filenames[0])
		tmpl = template.New(name).Funcs(t.funcs)
	}

	t.tmpl = template.Must(tmpl.ParseFiles(filenames...))
	t.funcs = nil

	return t
}

func (t *Template) Parse(text string) *Template {
	tmpl := t.tmpl
	if tmpl == nil {
		tmpl = template.New("").Funcs(t.funcs)
	}

	t.tmpl = template.Must(tmpl.Parse(text))
	t.funcs = nil

	return t
}

func (t *Template) Clone() *Template {
	return &Template{
		tmpl: template.Must(t.tmpl.Clone()),
	}
}

func (t *Template) NewPage() *Page {
	return &Page{
		tmpl: t.tmpl,
	}
}
