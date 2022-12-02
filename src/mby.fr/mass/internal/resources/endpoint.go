package resources

import "fmt"

type Endpoint struct {
	base `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (e Endpoint) init() (err error) {
	err = e.base.init()
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
	base, err := buildBase(EndpointKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Endpoint{base: base, Project: &project}
	return
}
