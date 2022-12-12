package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"mby.fr/mass/internal/config"
)

const (
	fullNameSeparator = '/'
)

type Uid struct {
    Kind Kind
    FullName string
}

type Resourcer interface {
	Kind() Kind
	FullName() string
	QualifiedName() string
	Dir() string
	//Config() (config.Config, error)
	Match(string, Kind) bool
	init() error
	backingFilepath() string
}

type fileBase struct {
	ResourceKind               Kind `yaml:"resourceKind"`
	name, dir string
}

func (r fileBase) Kind() Kind {
	return r.ResourceKind
}

func (r fileBase) FullName() string {
	return r.name
}

func (r fileBase) QualifiedName() string {
	return fmt.Sprintf("%s/%s", r.Kind(), r.FullName())
}

func (r fileBase) Dir() string {
	return r.dir
}

/*
func (r fileBase) Config() (config.Config, error) {
	c, err := config.Read(r.Dir())
	return c, err
}
*/
func (r fileBase) init() (err error) {
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

func (r fileBase) Match(name string, k Kind) bool {
	return name == r.FullName() && (k == AllKind || k == r.Kind())
}

func (r fileBase) backingFilepath() string {
	resourceFile := fmt.Sprintf("%s-%s.yaml", r.ResourceKind, r.name)
	return filepath.Join(r.Dir(), resourceFile)
}

func buildFileBase(kind Kind, dirPath, name string) (b base, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	name := resourceName(dirPath)
	b = fileBase{ResourceKind: kind, name: name, dir: absDirPath}
	return
}

type directoryBase struct {
	base
}

func (r directoryBase) backingFilepath() string {
	return filepath.Join(r.Dir(), DefaultResourceFile)
}

func buildDirectoryBase(kind Kind, dirPath string) (b base, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	name := resourceName(dirPath)
	b = directoryBase{ResourceKind: kind, name: name, dir: absDirPath}
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

func buildTestable(res Resourcer) (t testable, err error) {
	testDir := DefaultTestDir
	t = testable{resource: res, testDirectory: testDir}
	return
}

type Configer interface {
	Config() (config.Config, error)
}

type Env struct {
	base `yaml:"base,inline"` // Implicit composition: "golang inheritance"
}

func (e Env) init() (err error) {
	err = e.base.init()
	return
}

func buildEnv(parentDir, name string) (r Env, err error) {
	resourceDir := filepath.Join(parentDir, name)
	base, err := buildDirectoryBase(EnvKind, resourceDir)
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
