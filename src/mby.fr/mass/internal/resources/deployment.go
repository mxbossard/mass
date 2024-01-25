package resources

import (
	"os"
	"path/filepath"
	"strings"
)

var deploymentDirPrefix = "dep-"

type Deployment struct {
	directoryBase   `yaml:"base,inline"`
	configurableDir `yaml:"-"` // Ignore this field for yaml marshalling
	testable        `yaml:"testable,inline"`
	versionable     `yaml:"versionable,inline"`

	DeploymentDeps []string `yaml:"dependencies.deployments"`
	ImageDeps      []string `yaml:"dependencies.images"`

	SourceDirectory string `yaml:"sourceDirectory"`
	UpCmd           string `yaml:"upCmd"`
	DownCmd         string `yaml:"downCmd"`
	StartupCmd      string `yaml:"startupCmd"`
	LivenessCmd     string `yaml:"livenessCmd"`
	ReadinesCmd     string `yaml:"readinessCmd"`

	project Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (d Deployment) init() (err error) {
	err = d.directoryBase.init()
	if err != nil {
		return
	}
	err = d.testable.init()
	if err != nil {
		return
	}
	err = d.versionable.init()
	if err != nil {
		return
	}

	// Init source directory
	err = os.MkdirAll(d.AbsSourceDir(), 0755)

	return
}

func (d Deployment) AbsSourceDir() string {
	return absResourcePath(d.Dir(), d.SourceDirectory)
}

func (d Deployment) FullName() string {
	p, _ := d.Project()
	return p.FullName() + "/" + d.DeploymentName()
}

func (d Deployment) DeploymentName() string {
	return d.directoryBase.name
}

func (d Deployment) FullTaggedName() string {
	if d.Version() != "" {
		return strings.ToLower(d.FullName()) + ":" + d.Version()
	} else {
		return strings.ToLower(d.FullName()) + ":latest"
	}
}

func (d Deployment) Match(name string, k Kind) bool {
	return d.directoryBase.Match(name, k) || name == d.DeploymentName() && KindsMatch(k, d.Kind())
}

func (d Deployment) Project() (project Project, err error) {
	// Lazy loading
	if "" == d.project.directoryBase.base.name {
		projectDir := filepath.Dir(d.Dir())
		project, err = Read[Project](projectDir)
		if err != nil {
			return
		}
		d.project = project
	}
	return d.project, err
}

func buildDeployment(project Project, deploymentName string) (r Deployment, err error) {
	resDir := filepath.Join(project.Dir(), "dep-"+deploymentName)
	version := DefaultInitialVersion
	sourceDir := DefaultSourceDir

	b, err := buildDirectoryBase(ImageKind, resDir, deploymentName)
	if err != nil {
		return
	}

	r = Deployment{
		directoryBase:   b,
		SourceDirectory: sourceDir,
		project:         project,
	}
	r.configurableDir = buildConfigurableDir(b)

	t, err := buildTestable(r)
	if err != nil {
		return
	}
	r.testable = t

	versionable := buildVersionable(r, version)
	r.versionable = versionable

	return
}
