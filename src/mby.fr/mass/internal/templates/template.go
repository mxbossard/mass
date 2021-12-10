package templates

import (
	"io"
	"embed"
	"text/template"
)

const ConfigFilename = "config.yaml"

////go:embed image/* template/*
//go:embed src/*
var templates embed.FS

func read(name string) (data string, err error) {
	name = "src/" + name
	content, err := templates.ReadFile(name)
	if err != nil {
		return
	}
	data = string(content)
	return
}

func Make(name string) (t *template.Template, err error) {
	content, err := read(name)
	if err != nil {
		return
	}

	t, err = template.New(name).Parse(content)
	//t, err = template.ParseFS(templates, "src/" + name)
	return
}

func Render(name string, target io.Writer, data interface{}) (err error) {
	t, err := Make(name)
	if err != nil {
		return
	}

	t = t.Option("missingkey=error")
	//t = t.Delims("", "")
	err = t.Execute(target, data)
	return
}
