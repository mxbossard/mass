package resources

import "fmt"

type Service struct {
	base `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (s Service) init() (err error) {
	err = s.base.init()
	if err != nil {
		return
	}
	return
}

func buildService(projectPath, name string) (r Service, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}
	backingFilename := fmt.Sprintf("svc-%s.yaml", name)
	base, err := buildBase(ServiceKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Service{base: base, Project: &project}
	return
}
