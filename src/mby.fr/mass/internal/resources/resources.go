package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

const DefaultSourceDir = "src"
const DefaultTestDir = "test"
const DefaultVersionFile = "version.txt"
const DefaultInitialVersion = "0.0.1"
const DefaultBuildFile = "Dockerfile"
const DefaultDeployFile = "compose.yaml"
const DefaultResourceFile = "resource.yaml"

type Resource interface {
	Kind() Kind
	Name() string
	QualifiedName() string
	Dir() string
	Config() (config.Config, error)
	Init() error
	Match(string, Kind) bool
}

type Base struct {
	//Resource
	ResourceKind Kind `yaml:"resourceKind"`
	name, dir    string
}

func (r Base) Kind() Kind {
	return r.ResourceKind
}

func (r Base) Name() string {
	return r.name
}

func (r Base) QualifiedName() string {
	return fmt.Sprintf("%s/%s", r.ResourceKind, r.Name())
}

func (r Base) Dir() string {
	return r.dir
}

func (r Base) Config() (config.Config, error) {
	c, err := config.Read(r.Dir())
	return c, err
}

func (r Base) Init() (err error) {
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

func (r Base) Match(name string, k Kind) bool {
	return name == r.Name() && (k == AllKind || k == r.Kind())
}

type Tester interface {
	Resource
	AbsTestDir() string
}

type Testable struct {
	Base          `yaml:"base,inline"`
	TestDirectory string `yaml:"testDirectory"`
}

func (t Testable) AbsTestDir() string {
	return absResourvePath(t.Dir(), t.TestDirectory)
}

func (t Testable) Init() (err error) {
	// Create test dir
	err = os.MkdirAll(t.AbsTestDir(), 0755)
	return
}

type Env struct {
	Base `yaml:"base,inline"` // Implicit composition: "golang inheritance"
}

func (e Env) Init() (err error) {
	err = e.Base.Init()
	if err != nil {
		return
	}
	err = Write(e)
	return
}

type Project struct {
	//Base       `yaml:"base,inline"`
	Testable   `yaml:"testable,inline"`
	images     []Image
	DeployFile string `yaml:"deployFile"`
}

func (p Project) Init() (err error) {
	err = p.Base.Init()
	if err != nil {
		return
	}
	err = p.Testable.Init()
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

func (p *Project) Images() ([]Image, error) {
	var err error = nil
	if len(p.images) == 0 {
		images, err := ScanImages(p.Dir())
		if err != nil {
			return []Image{}, err
		}
		p.images = images
	}
	return p.images, err
}

type Image struct {
	//Base            `yaml:"base,inline"`
	Testable        `yaml:"testable,inline"`
	SourceDirectory string  `yaml:"sourceDirectory"`
	BuildFile       string  `yaml:"buildFile"`
	Version         string  `yaml:"version"`
	Project         Project `yaml:"-"` // Ignore this field for yaml marshalling
}

func (i Image) Init() (err error) {
	err = i.Base.Init()
	if err != nil {
		return
	}
	err = i.Testable.Init()
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
	return i.name
}

func (i Image) Name() string {
	return i.Project.Name() + "/" + i.ImageName()
}

func (i Image) FullName() string {
	if i.Version != "" {
		return strings.ToLower(i.Name()) + ":" + i.Version
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
	return i.Base.Match(name, k) || name == i.ImageName() && (k == AllKind || k == i.Kind())
}

func buildBase(kind Kind, path string) (b Base, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	b = Base{kind, name, absPath}
	return
}

func buildTestable(kind Kind, path string) (t Testable, err error) {
	testDir := DefaultTestDir
	base, err := buildBase(kind, path)
	if err != nil {
		return
	}
	t = Testable{base, testDir}
	return
}

func Init(path string, kind Kind) (b Base, err error) {
	switch kind {
	case EnvKind:
		var r Env
		r, err = BuildEnv(path)
		r.Init()
		b = r.Base
	case ProjectKind:
		var p Project
		p, err = BuildProject(path)
		p.Init()
		b = p.Base
	case ImageKind:
		var i Image
		i, err = BuildImage(path)
		i.Init()
		b = i.Base
	default:
		err = fmt.Errorf("Unable to load Resource from path: %s ! Not supported kind property: [%s].", path, kind)
	}

	return
}

func BuildEnv(path string) (r Env, err error) {
	base, err := buildBase(EnvKind, path)
	if err != nil {
		return
	}

	r = Env{Base: base}
	return
}

func BuildProject(path string) (r Project, err error) {
	testable, err := buildTestable(ProjectKind, path)
	if err != nil {
		return
	}

	deployfile := DefaultDeployFile

	r = Project{Testable: testable, DeployFile: deployfile}
	return
}

func BuildImage(path string) (r Image, err error) {
	testable, err := buildTestable(ImageKind, path)
	if err != nil {
		return
	}

	version := DefaultInitialVersion
	buildfile := DefaultBuildFile
	sourceDir := DefaultSourceDir

	projectPath := filepath.Dir(path)
	project, err := BuildProject(projectPath)
	if err != nil {
		return
	}

	r = Image{
		Testable:        testable,
		Version:         version,
		BuildFile:       buildfile,
		SourceDirectory: sourceDir,
		Project:         project,
	}
	return
}

func Read(path string) (r Resource, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	resourceFilepath := filepath.Join(path, DefaultResourceFile)
	content, err := os.ReadFile(resourceFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = ResourceNotFound{path, NewKindSet(AllKind)}
		}
		return
	}
	//fmt.Println("ResourceFile content:", string(content))

	base := Base{}
	err = yaml.Unmarshal(content, &base)
	if err != nil {
		return
	}

	base.name = filepath.Base(path)
	base.dir = path

	kind := base.Kind()
	switch kind {
	case EnvKind:
		res, err := BuildEnv(base.Dir())
		if err != nil {
			return r, err
		}
		//res.Base = base
		err = yaml.Unmarshal(content, &res)
		r = res
	case ProjectKind:
		res, err := BuildProject(base.Dir())
		if err != nil {
			return r, err
		}
		//res := Project{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = res
	case ImageKind:
		res, err := BuildImage(base.Dir())
		if err != nil {
			return r, err
		}
		//res := Image{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = res
	default:
		err = fmt.Errorf("Unable to load Resource from path: %s ! Not supported kind property: [%s].", resourceFilepath, kind)
		return
	}

	return
}

var writeLock = &sync.Mutex{}

func Write(r Resource) (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()

	var content []byte
	switch r := r.(type) {
	case Env, Project, Image:
		content, err = yaml.Marshal(r)
	//case Project:
	//content, err = yaml.Marshal(r)
	//case Image:
	//content, err = yaml.Marshal(r)
	default:
		err = fmt.Errorf("Unable to write Resource ! Not supported kind property: [%T].", r)
		return
	}

	if err != nil {
		return
	}

	resourceFilepath := filepath.Join(r.Dir(), DefaultResourceFile)
	err = os.WriteFile(resourceFilepath, content, 0644)

	return
}

func resourceName(path string) string {
	return filepath.Base(path)
}

// Return a resource relative path from an absolute path
func relResourcePath(resRootPath string, resPath string) (path string, err error) {
	resPath, err = filepath.Abs(resPath)
	if err != nil {
		return
	}
	resRootPath, err = filepath.Abs(resRootPath)
	if err != nil {
		return
	}
	path = strings.TrimPrefix(resPath, resRootPath)
	return
}

// Return a absolute path from a relative resource path
func absResourvePath(relRootPath string, resPath string) (path string) {
	path = filepath.Join(relRootPath, resPath)
	return
}
