package templates

import (
	"io"
	"io/fs"
	"os"
	//"fmt"
	"embed"
	"strings"
	"text/template"
	"path/filepath"
)

const ConfigTemplate = "config.yaml"

////go:embed image/* template/*
//go:embed templates/*
var templates embed.FS
const templatesRootDir = "templates"

type Renderer struct {
	templatesDir string
	templatesFs fs.FS
	rootDir string
}

func New(templatesDir string) Renderer {
	var templatesFs fs.FS
	var rootDir string
	if templatesDir == "" {
		templatesFs = templates
		rootDir = templatesRootDir
	} else {
		templatesFs = os.DirFS(templatesDir)
		rootDir = "."
	}
	r := Renderer{templatesDir, templatesFs, rootDir}
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

	// Copy templates FS into targetDir ignoring root dir
	targetDir := templatesDir
	rootDir := templatesRootDir // do not copy this directory

	copyFunc := func(path string, d fs.DirEntry, err error) error {
		//fmt.Printf("Walking path: %s\n", path)
		targetPath := strings.TrimPrefix(path, rootDir)
		if d == nil {
			return nil
		} else if d.IsDir() {
			// Create dir
			dest := filepath.Join(targetDir, targetPath)
			//fmt.Printf("Creating dir %s ...\n", dest)
			err = os.MkdirAll(dest, 0755)
			if err != nil {
				return err
			}
		} else {
			// Copy file from FS/path into targetDir/targetPath
			//fileName := d.Name()
			filePath := filepath.Join(targetDir, targetPath)
			//fmt.Printf("Copying template %s into dir %s ...\n", fileName, filePath)

			fileContent, err := templates.ReadFile(path)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
				return err
			}
		}
		return nil
	}

	err = fs.WalkDir(templates, rootDir, copyFunc)
	return
}

func (r Renderer) read(name string) (data string, err error) {
	name = filepath.Join(r.rootDir, name)
	file, err := r.templatesFs.Open(name)
	if err != nil {
		return
	}
	defer file.Close()

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
