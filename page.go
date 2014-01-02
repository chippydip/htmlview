package htmlview

import (
	"html/template"
	"net/http"
)

type Page struct {
	tmpl       *template.Template
	args       map[string]interface{}
	header     http.Header
	StatusCode int
}

func (p *Page) Set(name string, value interface{}) *Page {
	if p.args == nil {
		p.args = map[string]interface{}{}
	}
	p.args[name] = value

	return p
}

func (p *Page) Get(name string) interface{} {
	return p.args[name]
}

func (p *Page) Header() http.Header {
	if p.header == nil {
		p.header = http.Header{}
	}
	return p.header
}

func (p *Page) Render(w http.ResponseWriter) error {
	// Copy headers
	if len(p.header) > 0 {
		hdr := w.Header()
		for k, v := range p.header {
			hdr[k] = v
		}
	}

	// Set the return code
	if p.StatusCode != 0 {
		w.WriteHeader(p.StatusCode)
	}

	// Execute the template
	return p.tmpl.Execute(w, p.args)
}
