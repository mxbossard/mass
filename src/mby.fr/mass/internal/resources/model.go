package resources

import (
	"fmt"
	"os"
	"strings"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

type Resourcer interface {
	Kind() Kind
	Name() string
	QualifiedName() string
	Dir() string
	Config() (config.Config, error)
	Init() error
	Match(string, Kind) bool
}

type base struct {
	ResourceKind Kind `yaml:"resourceKind"`
	name, dir    string
}

func (r base) Kind() Kind {
	return r.ResourceKind
}

func (r base) Name() string {
	return r.name
}

func (r base) QualifiedName() string {
	return fmt.Sprintf("%s/%s", r.Kind(), r.Name())
}

func (r base) Dir() string {
	return r.dir
}

func (r base) Config() (config.Config, error) {
	c, err := config.Read(r.Dir())
	return c, err
}

func (r base) Init() (err error) {
	// Create resource dir
	err = os.MkdirAll(r.Dir(), 0755)
	if err != nil {
		return
	}

	// Init config
	err = config.Init(r.Dir(), r)
	if err != nil {
		return
	}

	return
}

func (r base) Match(name string, k Kind) bool {
	return name == r.Name() && (k == AllKind || k == r.Kind())
}

type Tester interface {
	AbsTestDir() string
}

type Testable struct {
	resource      Resourcer //`yaml:"base,inline"`
	testDirectory string    `yaml:"testDirectory"`
}

func (t Testable) AbsTestDir() string {
	return absResourvePath(t.resource.Dir(), t.testDirectory)
}

func (t Testable) Init() (err error) {
	// Create test dir
	err = os.MkdirAll(t.AbsTestDir(), 0755)
	return
}

type Env struct {
	base `yaml:"base,inline"` // Implicit composition: "golang inheritance"
}

func (e Env) Init() (err error) {
	err = e.base.Init()
	if err != nil {
		return
	}
	err = Write(e)
	return
}

type Project struct {
	//Resourcer //`yaml:"base,inline"`

	base     `yaml:"base,inline"`
	Testable `yaml:"testable,inline"`

	images     []*Image
	DeployFile string `yaml:"deployFile"`
}

func (p Project) Init() (err error) {
	err = p.base.Init()
	if err != nil {
		return
	}

	// Init Deploy file
	deployfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = file.SoftInitFile(p.DeployFile, deployfileContent)
	if err != nil {
		return
	}

	err = Write(p)
	return
}

func (p Project) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + p.Name()
	return
}

func (p Project) AbsDeployFile() string {
	return absResourvePath(p.Dir(), p.DeployFile)
}

func (p *Project) Images() ([]*Image, error) {
	var err error = nil
	if len(p.images) == 0 {
		images, err := ScanImages(p.Dir())
		if err != nil {
			return []*Image{}, err
		}
		p.images = images
	}
	return p.images, err
}

type Image struct {
	//Resourcer //`yaml:"base,inline"`

	base        `yaml:"base,inline"`
	Testable    `yaml:"testable,inline"`
	Versionable `yaml:"versionable,inline"`

	SourceDirectory string `yaml:"sourceDirectory"`
	BuildFile       string `yaml:"buildFile"`
	//Version         string  `yaml:"version"`
	Project Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (i Image) Init() (err error) {
	err = i.base.Init()
	if err != nil {
		return
	}

	// Init version file
	//versionFile := versionFilepath(projectPath)
	//_, err = file.SoftInitFile(versionFile, resources.DefaultInitialVersion)

	// Init source directory
	err = os.MkdirAll(i.AbsSourceDir(), 0755)

	// Init Build file
	buildfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = file.SoftInitFile(i.BuildFile, buildfileContent)

	if err != nil {
		return
	}
	err = Write(i)
	return
}

func (i Image) AbsSourceDir() string {
	return absResourvePath(i.Dir(), i.SourceDirectory)
}

func (i Image) AbsBuildFile() string {
	return absResourvePath(i.Dir(), i.BuildFile)
}

func (i Image) ImageName() string {
	return i.base.Name()
}

func (i Image) Name() string {
	return i.Project.Name() + "/" + i.ImageName()
}

func (i Image) FullName() string {
	if i.Version() != "" {
		return strings.ToLower(i.Name()) + ":" + i.Version()
	} else {
		return strings.ToLower(i.Name()) + ":latest"
	}
}

func (i Image) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + i.Name()
	return
}

func (i Image) Match(name string, k Kind) bool {
	return i.base.Match(name, k) || name == i.ImageName() && (k == AllKind || k == i.Kind())
}

func (i Image) GetVersionable() *Versionable {
	return &(i.Versionable)
}

func (i *Image) SetVersionable(v Versionable) {
	i.Versionable = v
}

func (i Image) Version() string {
	return i.Versionable.version()
}

func (i *Image) Bump(bumpMinor, bumpMajor bool) (msg string, err error) {
	msg, err = i.Versionable.bump(bumpMinor, bumpMajor)
	Write(i)
	return
}

func (i *Image) Promote() (msg string, err error) {
	msg, err = i.Versionable.promote()
	Write(i)
	return
}

func (i *Image) Release() (msg string, err error) {
	msg, err = i.Versionable.release()
	Write(i)
	return
}
