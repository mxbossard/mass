package resources

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

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
	Match(string, Kind) bool
	init() error
	backingFilepath() string
}

type base struct {
	ResourceKind Kind `yaml:"resourceKind"`
	name, dir, backingFilename    string
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

func (r base) init() (err error) {
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

func (r base) backingFilepath() string {
	return filepath.Join(r.Dir(), r.backingFilename)
}

type Tester interface {
	AbsTestDir() string
}

type testable struct {
	resource      Resourcer //`yaml:"base,inline"`
	testDirectory string    `yaml:"testDirectory"`
}

func (t testable) AbsTestDir() string {
	return absResourcePath(t.resource.Dir(), t.testDirectory)
}

func (t testable) init() (err error) {
	// Create test dir
	err = os.MkdirAll(t.AbsTestDir(), 0755)
	return
}

type Env struct {
	base `yaml:"base,inline"` // Implicit composition: "golang inheritance"
}

func (e Env) init() (err error) {
	err = e.base.init()
	return
}

type Project struct {
	base     `yaml:"base,inline"`
	testable `yaml:"testable,inline"`

	images     []*Image
	DeployFile string `yaml:"deployFile"`
}

func (p Project) init() (err error) {
	err = p.base.init()
	if err != nil {
		return
	}
	err = p.testable.init()
	if err != nil {
		return
	}

	// Init Deploy file
	deployfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = file.SoftInitFile(p.DeployFile, deployfileContent)

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
	return absResourcePath(p.Dir(), p.DeployFile)
}

func (p *Project) Images() ([]*Image, error) {
	var err error = nil
	if len(p.images) == 0 {
		images, err := Scan[Image](p.Dir())
		var imagesPtrs []*Image
		for i := 0; i < len(images); i++ {
			imagesPtrs = append(imagesPtrs, &images[i])
		}
		if err != nil {
			return []*Image{}, err
		}
		p.images = imagesPtrs
	}
	return p.images, err
}

type Image struct {
	base        `yaml:"base,inline"`
	testable    `yaml:"testable,inline"`
	versionable `yaml:"versionable,inline"`

	SourceDirectory string  `yaml:"sourceDirectory"`
	BuildFile       string  `yaml:"buildFile"`
	Project         Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (i Image) init() (err error) {
	err = i.base.init()
	if err != nil {
		return
	}
	err = i.testable.init()
	if err != nil {
		return
	}
	err = i.versionable.init()
	if err != nil {
		return
	}

	// Init source directory
	err = os.MkdirAll(i.AbsSourceDir(), 0755)

	// Init Build file
	buildfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = file.SoftInitFile(i.BuildFile, buildfileContent)

	return
}

func (i Image) AbsSourceDir() string {
	return absResourcePath(i.Dir(), i.SourceDirectory)
}

func (i Image) AbsBuildFile() string {
	return absResourcePath(i.Dir(), i.BuildFile)
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


type Endpoint struct {
	base        `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (e Endpoint) init() (err error) {
	err = e.base.init()
	if err != nil {
		return
	}
	return
}

type Service struct {
	base        `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (s Service) init() (err error) {
	err = s.base.init()
	if err != nil {
		return
	}
	return
}

/*
type Ingress struct {
	base        `yaml:"base,inline"`

	Project *Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (i Ingress) init() (err error) {
	err = i.base.init()
	if err != nil {
		return
	}
	return
}
*/

