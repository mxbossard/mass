package resources

import (
	"io"

	"mby.fr/mass/internal/settings"
)

type Event struct {
	Timestamp   int64
	Name        string
	Description string
}

type Eventer interface {
	Events() []*Event
}

type Executioner interface {
	WaitCompletion() error
	ExitCode() (int, error)
	StdOut() io.Reader
	StdErr() io.Reader
}

type Execution struct {
	exitCode int
	stdOut   io.Reader
	stdErr   io.Reader
}

type Container interface {
	Eventer
	Run(args []string) (Executioner, error)
	Stop() error
}

type Prober interface {
	Eventer
	Probe() (bool, error)
}

type Pod struct {
	Eventer

	base     `yaml:"base,inline"`
	testable `yaml:"testable,inline"`
	//versionable `yaml:"versionable,inline"`

	Project        *Project `yaml:"-"` // Ignore this field for yaml marshalling
	events         []Event `yaml:"-"` // Ignore this field for yaml marshalling
	InitContainers []*Container
	Containers     []*Container
	StartupProbe   Prober
	ReadinessProbe Prober
	LivenessProbe  Prober
	RestartPolicy  string
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
	return p.base.Name()
}

func (p Pod) Name() string {
	return p.Project.Name() + "/" + p.PodName()
}

/*
func (p Pod) FullName() string {
	if i.Version() != "" {
		return strings.ToLower(p.Name()) + ":" + p.Version()
	} else {
		return strings.ToLower(p.Name()) + ":latest"
	}
}
*/

func (p Pod) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + p.Name()
	return
}

func (p Pod) Match(name string, k Kind) bool {
	return p.base.Match(name, k) || name == p.PodName() && (k == AllKind || k == p.Kind())
}

func (p Pod) Start() (err error) {
	return
}

func (p Pod) Events() (events []*Event) {
	return
}
