package main

import (
	"io"
	"log"
	"regexp"
	"text/template"
)

type Context map[string]interface{}

type TemplateParser struct {
	HTML string
}

func ValidateEmail(email string) bool {
	pattern := "^[a-zA-Z0-9._%+-^@]+@[a-zA-Z0-9.\\-^@]+\\.[a-zA-Z][a-zA-Z][a-zA-Z]?[a-zA-Z]?$"

	ok, err := regexp.Match(pattern, []byte(email))
	return err == nil && ok
}

func (tP *TemplateParser) Write(p []byte) (n int, err error) {
	tP.HTML += string(p)
	return len(p), nil
}

func ParseTemplate(templatesSet *template.Template, tmpl string, data interface{}) string {
	tp := &TemplateParser{}
	templatesSet.ExecuteTemplate(tp, tmpl, data)
	return tp.HTML
}

func RenderTemplate(templatesSet *template.Template, w io.Writer, tmpl string, data interface{}) {
	err := templatesSet.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Fatalf("Cannot render template: %s", err)
	}
}
