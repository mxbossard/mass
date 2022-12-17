package resources

import (
	"fmt"
	"path/filepath"

	"mby.fr/mass/internal/settings"
)

type Container struct {
	Image      string   `yaml:"image"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
	Cmd        []string `yaml:"cmd,omitempty"`
	Args       []string `yaml:"args,omitempty"`
}

type ProbeKind int

const (
	Http = ProbeKind(iota)
	Tcp
	Exec
	probeKindLimit
)

type Probe struct {
	Kind   ProbeKind         `yaml:"kind"`
	Config map[string]string `yaml:",inline"`
}

type RestartPolicy int

const (
	Always = RestartPolicy(iota)
	OnFailure
	Never
	restartPolicyLimit
)

type Pod struct {
	fileBase     `yaml:"base,inline"`
	testable `yaml:"testable,inline"`
	//versionable `yaml:"versionable,inline"`

	project        Project       `yaml:"-"` // Ignore this field for yaml marshalling
	InitContainers []*Container  `yaml:"initContainers,omitempty"`
	Containers     []*Container  `yaml:"containers"`
	StartupProbe   Probe         `yaml:"startupProbe,omitempty"`
	ReadinessProbe Probe         `yaml:"readinessProbe,omitempty"`
	LivenessProbe  Probe         `yaml:"livenessProbe,omitempty"`
	RestartPolicy  RestartPolicy `yaml:"restartPolicy,omitempty"`
}

func (p Pod) init() (err error) {
	err = p.base.init()
	if err != nil {
		return
	}
	err = p.testable.init()
	if err != nil {
		return
	}
	/*
		err = p.versionable.init()
		if err != nil {
			return
		}
	*/

	return
}

func (p Pod) PodName() string {
	return p.fileBase.name
}

func (p Pod) FullName() string {
	project, _ := p.Project()
	return project.FullName() + "/" + p.PodName()
}

/*
func (p Pod) FullName() string {
	if i.Version() != "" {
		return strings.ToLower(p.FullName()) + ":" + p.Version()
	} else {
		return strings.ToLower(p.FullName()) + ":latest"
	}
}
*/

func (p Pod) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + p.FullName()
	return
}

func (p Pod) Match(name string, k Kind) bool {
	return p.fileBase.Match(name, k) || name == p.PodName() && (k == AllKind || k == p.Kind())
}

func (p Pod) Project() (project Project, err error) {
	// Lazy loading
	if ("" == p.project.directoryBase.base.name) {
		projectDir := filepath.Dir(p.Dir())
		project, err = Read[Project](projectDir)
		if err != nil {
			return
		}
		p.project = project
	}
	return
}

func buildPod(projectPath, name string) (r Pod, err error) {
	backingFilename := fmt.Sprintf("pod-%s.yaml", name)
	b, err := buildFileBase(PodKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Pod{
		fileBase:    b,
	}

	t, err := buildTestable(r)
	if err != nil {
		return
	}
	r.testable = t

	return
}
