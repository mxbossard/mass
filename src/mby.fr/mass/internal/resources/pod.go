package resources

import (
	"mby.fr/mass/internal/settings"
)

type Container struct {
	Image string
	Entrypoint []string
	Cmd []string
	Args []string
}

type Probe struct {
}

type Pod struct {
	base     `yaml:"base,inline"`
	testable `yaml:"testable,inline"`
	//versionable `yaml:"versionable,inline"`

	Project        Project `yaml:"-"` // Ignore this field for yaml marshalling
	InitContainers []*Container
	Containers     []*Container
	StartupProbe   Probe
	ReadinessProbe Probe
	LivenessProbe  Probe
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
