package resources

import (
	"fmt"
	"path/filepath"
)

type Pod struct {
	fileBase `yaml:"base,inline"`
	//testable `yaml:"testable,inline"`

	project Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (p Pod) FullName() string {
	project, _ := p.Project()
	return project.FullName() + "/" + p.Name()
}

func (p Pod) Match(name string, k Kind) bool {
	return p.fileBase.Match(name, k) || name == p.Name() && (k == AllKind || k == p.Kind())
}

func (p Pod) Project() (project Project, err error) {
	// Lazy loading
	if "" == p.project.directoryBase.base.name {
		projectDir := filepath.Dir(p.Dir())
		project, err = Read[Project](projectDir)
		if err != nil {
			return
		}
		p.project = project
	}
	return
}

func forgePodResFilename(name string) string {
	return fmt.Sprintf("pod-%s.yaml", name)
}

func buildPod(projectPath, name string) (r Pod, err error) {
	backingFilename := forgePodResFilename(name)
	b, err := buildFileBase(PodKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Pod{
		fileBase: b,
	}

	return
}

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

type PodSpec struct {
	fileBase `yaml:"base,inline"`
	testable `yaml:"testable,inline"`
	//versionable `yaml:"versionable,inline"`

	name           string        `yaml:"name"`
	InitContainers []*Container  `yaml:"initContainers,omitempty"`
	Containers     []*Container  `yaml:"containers"`
	StartupProbe   Probe         `yaml:"startupProbe,omitempty"`
	ReadinessProbe Probe         `yaml:"readinessProbe,omitempty"`
	LivenessProbe  Probe         `yaml:"livenessProbe,omitempty"`
	RestartPolicy  RestartPolicy `yaml:"restartPolicy,omitempty"`
	//volumeMounts   []VolumeMount `yaml:"volumeMounts,omitempty"`

	project Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (p PodSpec) init() (err error) {
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
