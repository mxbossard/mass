package resources

import(
	"fmt"
	"os"
	"sync"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"mby.fr/mass/internal/config"
	"mby.fr/utils/file"
)

const DefaultSourceDir = "src"
const DefaultTestDir = "test"
const DefaultVersionFile = "version.txt"
const DefaultInitialVersion = "0.0.1"
const DefaultBuildFile = "Dockerfile"
const DefaultResourceFile = "resource.yaml"

const EnvKind = "Env"
const ProjectKind = "Project"
const ImageKind = "Image"

var ResourceNotFound error = fmt.Errorf("Resource not found")
var InconsistentResourceKind error = fmt.Errorf("Resource kind is not consistent")

type Resource interface {
	Kind() string
	Name() string
	Dir() string
	Config() (config.Config, error)
	Init() error
}

type Base struct {
	//Resource
	ResourceKind string `yaml:"resourceKind"`
	name, dir string
}

func (r Base) Kind() string {
	return r.ResourceKind
}

func (r Base) Name() string {
	return r.name
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

type Testable struct {
	TestDirectory string `yaml:"testDirectory"`
}

func (t Testable) TestDir() (string) {
	//var err error = nil
	return t.TestDirectory
}

func (t Testable) Init() (err error) {
	// Create test dir
	err = os.MkdirAll(t.TestDir(), 0755)
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
	Base `yaml:"base,inline"`
	Testable `yaml:"testable,inline"`
	images []Image
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
	err = Write(p)
	return
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
	Base `yaml:"base,inline"`
	Testable `yaml:"testable,inline"`
	SourceDirectory string `yaml:"sourceDirectory"`
	BuildFile string `yaml:"buildFile"`
	Version string `yaml:"version"`
	p Project
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
	err = os.MkdirAll(i.SourceDir(), 0755)

	// Init Build file
        _, err = file.SoftInitFile(i.BuildFile, "")

	if err != nil {
		return
	}
	err = Write(i)
	return
}

func (i Image) SourceDir() (string) {
	return i.SourceDirectory
}

func buildBase(kind, path string) (b Base, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	b = Base{kind, name, absPath}
	return
}

func buildTestable(path string) (t Testable, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	testDir := filepath.Join(absPath, DefaultTestDir)
	t = Testable{testDir}
	return
}

func Init(path, kind string) (b Base, err error) {
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
	base, err := buildBase(ProjectKind, path)
	if err != nil {
                return
        }

	testable, err := buildTestable(path)
	if err != nil {
                return
        }

	r = Project{Base: base, Testable: testable}
	return
}

func BuildImage(path string) (r Image, err error) {
	base, err := buildBase(ImageKind, path)
	if err != nil {
                return
        }

	testable, err := buildTestable(path)
	if err != nil {
                return
        }

	version := DefaultInitialVersion
        buildfile := filepath.Join(base.Dir(), DefaultBuildFile)
	sourceDir := filepath.Join(base.Dir(), DefaultSourceDir)

	r = Image{Base: base, Testable: testable, Version: version, BuildFile: buildfile, SourceDirectory: sourceDir}
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
			err = ResourceNotFound
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
		res := Env{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = res
		case ProjectKind:
		res := Project{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = res
		case ImageKind:
		res := Image{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = res
		default:
		err = fmt.Errorf("Unable to load Resource from path: %s ! Not supported kind property: [%s].", resourceFilepath, kind)
		return
	}

	if err != nil {
		return
	}

	return
}

var writeLock = &sync.Mutex{}
func Write(r Resource) (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()

	var content []byte
	switch r.(type) {
		case Env:
		content, err = yaml.Marshal(r)
		case Project:
		content, err = yaml.Marshal(r)
		case Image:
		content, err = yaml.Marshal(r)
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

