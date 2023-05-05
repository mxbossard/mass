package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"mby.fr/mass/internal/config"
)

const (
	fullNameSeparator = "/"
)

type Uid struct {
	Kind     Kind
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
	//backingFilepath() string
}

type base struct {
	ResourceKind Kind   `yaml:"resourceKind"`
	name         string `yaml:"-"` // Ignore this field for yaml marshalling
	dir          string `yaml:"-"` // Ignore this field for yaml marshalling
}

func (r base) Kind() Kind {
	return r.ResourceKind
}

func (r base) FullName() string {
	return r.name
}

func (r base) QualifiedName() string {
	return fmt.Sprintf("%s/%s", r.Kind(), r.FullName())
}

func (r base) Dir() string {
	return r.dir
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
	return name == r.FullName() && (k == AllKind || k == r.Kind())
}

type fileBase struct {
	base `yaml:",inline"`
}

func (r fileBase) backingFilepath() string {
	resourceFile := fmt.Sprintf("%s-%s.yaml", r.ResourceKind, r.name)
	return filepath.Join(r.base.Dir(), resourceFile)
}

func buildFileBase(kind Kind, dirPath, name string) (b fileBase, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	b = fileBase{base{ResourceKind: kind, name: name, dir: absDirPath}}
	return
}

type directoryBase struct {
	base `yaml:",inline"`
}

func (r directoryBase) backingFilepath() string {
	return filepath.Join(r.base.Dir(), DefaultResourceFile)
}

func buildDirectoryBase(kind Kind, dirPath string) (b directoryBase, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	name := resourceName(dirPath)
	b = directoryBase{base{ResourceKind: kind, name: name, dir: absDirPath}}
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

type configurableDir struct {
	base directoryBase
}

func (c configurableDir) Config() (config.Config, error) {
	cfg, err := config.Read(c.base.Dir())
	return cfg, err
}

func buildConfigurableDir(res directoryBase) (c configurableDir) {
	c = configurableDir{base: res}
	return
}

type Env struct {
	directoryBase   `yaml:"base,inline"` // Implicit composition: "golang inheritance"
	configurableDir `yaml:"-"`           // Ignore this field for yaml marshalling
}

func (e Env) init() (err error) {
	err = e.directoryBase.init()
	return
}

func buildEnv(envDir string) (r Env, err error) {
	base, err := buildDirectoryBase(EnvKind, envDir)
	if err != nil {
		return
	}

	r = Env{directoryBase: base}
	r.configurableDir = buildConfigurableDir(base)
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
