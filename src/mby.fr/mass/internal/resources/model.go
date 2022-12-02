package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"mby.fr/mass/internal/config"
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
	ResourceKind               Kind `yaml:"resourceKind"`
	name, dir, backingFilename string
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

func buildBase(kind Kind, dirPath, backingFilename string) (b base, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	name := resourceName(dirPath)
	b = base{ResourceKind: kind, name: name, dir: absDirPath, backingFilename: backingFilename}
	return
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

func buildTestable(res Resourcer, path string) (t testable, err error) {
	testDir := DefaultTestDir
	t = testable{resource: res, testDirectory: testDir}
	return
}

type Env struct {
	base `yaml:"base,inline"` // Implicit composition: "golang inheritance"
}

func (e Env) init() (err error) {
	err = e.base.init()
	return
}

func buildEnv(path string) (r Env, err error) {
	base, err := buildBase(EnvKind, path, DefaultResourceFile)
	if err != nil {
		return
	}

	r = Env{base: base}
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
