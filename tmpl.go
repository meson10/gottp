package gottp

import (
	"html/template"

	"bytes"
	"os"
	"path/filepath"
)

var tmplPath string

const templatesPath = "templates"
const templateFile = "error.html"

func makeTmplPath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, templatesPath, templateFile)
}

func getFileTemplate() *template.Template {
	if tmplPath == "" {
		tmplPath = makeTmplPath()
	}

	return template.Must(template.ParseFiles(tmplPath))
}

func getStaticTemplate() *template.Template {
	return template.Must(template.New("error_email").Parse(templateContent))
}

func ErrorTemplate(stack *ErrorStack) string {
	var doc bytes.Buffer
	//t := getFileTemplate()

	t := getStaticTemplate()
	t.Execute(&doc, *stack)

	return doc.String()
}
