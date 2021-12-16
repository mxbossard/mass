package templates

import (
	"io"
	"io/fs"
	"os"
	"fmt"
	"embed"
	"strings"
	"text/template"
	"path/filepath"
)

const ConfigTemplate = "config.yaml"

////go:embed image/* template/*
//go:embed src/*
var templates embed.FS

type Renderer struct {
	templatesDir string
	templatesFs fs.FS
}

func New(templatesDir string) Renderer {
	var templatesFs fs.FS
	if templatesDir == "" {
		templatesFs = templates
	} else {
		templatesFs = os.DirFS(templatesDir)
	}
	r := Renderer{templatesDir, templatesFs}
	return r
}

// Copy embeded templates into a directory
func Init(templatesDir string) (err error) {
	if templatesDir == "" {
		// No templatesDir so nothing to init
		return
	}

	templatesDir, err = filepath.Abs(templatesDir)
	if err != nil {
		return
	}

	// Copy embeded src dir in settings dir
	dirPath := "src"
	dirEntries, err := templates.ReadDir(dirPath)
	if err != nil {
		return
	}

	targetDir := filepath.Join(templatesDir, dirPath)
	err = os.Mkdir(targetDir, 0755)
	if err != nil {
		return
	}

	for _, dirEntry := range dirEntries {
		fileName := dirEntry.Name()
		fileContent, err := templates.ReadFile(fmt.Sprintf("%s/%s", dirPath, fileName))
		if err != nil {
			return err
		}

		filePath := filepath.Join(targetDir, fileName)
		if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
			return err
		}
		//fmt.Printf("Copied template %s into dir %s ...\n", fileName, filePath)
	}
	return
}

func (r Renderer) read(name string) (data string, err error) {
	name = "src/" + name
	file, err := r.templatesFs.Open(name)
	defer file.Close()
	if err != nil {
		return
	}

	builder := strings.Builder{}
	const maxSz = 64
	// create buffer
	b := make([]byte, maxSz)

	for {
		// read content to buffer
		readTotal, err := file.Read(b)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		builder.Write(b[:readTotal])
	}

	data = builder.String()
	return
}

func (r Renderer) Make(name string) (t *template.Template, err error) {
	content, err := r.read(name)
	if err != nil {
		return
	}

	t, err = template.New(name).Parse(content)
	//t, err = template.ParseFS(templates, "src/" + name)
	return
}

func (r Renderer) Render(name string, target io.Writer, data interface{}) (err error) {
	t, err := r.Make(name)
	if err != nil {
		return
	}

	t = t.Option("missingkey=error")
	//t = t.Delims("", "")
	err = t.Execute(target, data)
	return
}

func (r Renderer) RenderToFile(name, path string, data interface{}) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	err = r.Render(name, file, data)
	if err != nil {
		os.Remove(path)
	}
	return
}
