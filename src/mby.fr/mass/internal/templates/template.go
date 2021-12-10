package templates

import (
	"io"
	"os"
	"embed"
	"strings"
	"text/template"
)

const ConfigTemplate = "config.yaml"

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
	// Remove extra new ling added for no reason
	data = strings.TrimSuffix(data, "\n")
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

func RenderToFile(name, path string, data interface{}) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	err = Render(name, file, data)
	if err != nil {
		os.Remove(path)
	}
	return
}
