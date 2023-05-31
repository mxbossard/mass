package resources

import (
	"fmt"
	"path/filepath"
)

type ServiceType int
const(
	ComposeService = ServiceType(iota)
	K8sService
	HelmService
	KustomizeService
)

type Service struct {
	fileBase `yaml:"base,inline"`

	project Project `yaml:"-"` // Ignore this field for yaml marshalling

	ServiceDir string
	ServiceFile string
	ServiceKind ServiceType
}

func (s Service) init() (err error) {
	err = s.fileBase.init()
	if err != nil {
		return
	}
	return
}

func (s Service) Project() (project Project, err error) {
	// Lazy loading
	if "" == s.project.directoryBase.base.name {
		projectDir := filepath.Dir(s.Dir())
		project, err = Read[Project](projectDir)
		if err != nil {
			return
		}
		s.project = project
	}
	return s.project, err
}

func (s Service) FullName() string {
	project, _ := s.Project()
	return project.FullName() + "/" + s.Name()
}

func (s Service) Match(name string, k Kind) bool {
	return s.fileBase.Match(name, k) || name == s.Name() && (k == AllKind || k == s.Kind())
}

func forgeServiceResFilename(name string) string {
	return fmt.Sprintf("svc-%s.yaml", name)
}

func buildService(projectPath, name string) (r Service, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}
	backingFilename := forgeServiceResFilename(name)
	base, err := buildFileBase(ServiceKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Service{fileBase: base, project: project}
	return
}
