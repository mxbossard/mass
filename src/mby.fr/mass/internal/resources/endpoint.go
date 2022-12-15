package resources

import "fmt"

type Endpoint struct {
	fileBase `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (e Endpoint) init() (err error) {
	err = e.fileBase.init()
	if err != nil {
		return
	}
	return
}

func buildEndpoint(projectPath, name string) (r Endpoint, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}
	backingFilename := fmt.Sprintf("end-%s.yaml", name)
	base, err := buildFileBase(EndpointKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Endpoint{fileBase: base, Project: &project}
	return
}
