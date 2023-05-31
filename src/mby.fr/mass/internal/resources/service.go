package resources

import (
	"fmt"
	"path/filepath"
	"regexp"
	"log"
)

type ServiceType int
const(
	ComposeService = ServiceType(iota)
	K8sService
	HelmService
	KustomizeService
	serviceTypeLimit
)

func (t ServiceType) String() (s string) {
	switch t {
	case K8sService:
		return "k8s"
	case HelmService:
		return "helm"
	case KustomizeService:
		return "kustomize"
	case ComposeService:
		return "compose"
	}
	log.Fatalf("ServiceType %T not configured !", t)
	return
}

func (t ServiceType) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (st *ServiceType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	// Iterate over all kinds
	for serviceType := ServiceType(1); serviceType < serviceTypeLimit; serviceType++ {
		if s == serviceType.String() {
			*st = serviceType
			return nil
		}
	}

	return fmt.Errorf("Unable to unmarshal kind: %s", s)
}

type Service struct {
	fileBase `yaml:"base,inline"`

	project Project `yaml:"-"` // Ignore this field for yaml marshalling

	ServiceType ServiceType `yaml:"serviceType"`
	ServiceDir string `yaml:"serviceDir"`
	ServiceFile string `yaml:"serviceFile"`
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

func serviceNameFromResFilename(filename string) (string, error) {
	re := regexp.MustCompile(`^svc-(.+)\.yaml$`)
	submatch := re.FindStringSubmatch(filename)
	if len(submatch) == 2 {
		return submatch[1], nil
	}
	return "", fmt.Errorf("Bad service filename: %s !", filename)
}

func buildService(project Project, name string) (r Service, err error) {
	backingFilename := forgeServiceResFilename(name)
	base, err := buildFileBase(ServiceKind, project.Dir(), backingFilename)
	if err != nil {
		return
	}

	r = Service{fileBase: base, project: project}
	r.name = name
	return
}
