package htmlview

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type PageFactory func(*http.Request) (*Page, error)

func NewPageFactory(t *Template) PageFactory {
	return func(*http.Request) (*Page, error) {
		return t.NewPage(), nil
	}
}

func NewErrorPageFactory(t *Template, statusCode int) PageFactory {
	return func(*http.Request) (*Page, error) {
		p := t.NewPage()
		p.StatusCode = statusCode
		return p, nil
	}
}

func (pf PageFactory) New(r *http.Request) (*Page, error) {
	if pf == nil {
		return nil, nil
	}
	return pf(r)
}

func (pf PageFactory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DefaultServe(pf, w, r)
}

func (pf PageFactory) Render(w http.ResponseWriter, r *http.Request) (*Page, error) {
	p, err := pf.New(r)
	if err != nil || p == nil {
		return p, err
	}

	return p, p.Render(w)
}

var DefaultServe = func(pf PageFactory, w http.ResponseWriter, r *http.Request) {
	// Catch and report any panics
	defer func() {
		if err := recover(); err != nil {
			code := http.StatusInternalServerError
			msg := fmt.Sprintf("%v: %v\n%v", code, err, string(debug.Stack()))
			http.Error(w, msg, code)
		}
	}()

	// Generate the page
	p, err := pf.Render(w, r)
	if err != nil {
		panic(err)
	}

	// Page not generated?
	if p == nil {
		http.NotFound(w, r)
	}
}
