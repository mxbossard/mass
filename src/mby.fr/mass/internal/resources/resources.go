package resources

import(
	"fmt"
	"os"
	"sync"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"mby.fr/mass/internal/config"
)

const defaultResourceFile = "resource.yaml"

const EnvKind = "Env"
const ProjectKind = "Project"
const ImageKind = "Image"

type Resource interface {
	Kind() string
	Name() string
	Dir() string
	Config() (config.Config, error)
}

type Base struct {
	ResourceKind, name, dir string
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

type Env struct {
	Base // Implicit composition: "golang inheritance"
}

type Project struct {
	Base
	images []Image
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

func (p *Project) TestDir() (string) {
	//var err error = nil
	return ""
}

type Image struct {
	Base
	project Project
}

func (i *Image) TestDir() (string) {
	//var err error = nil
	return ""
}

func (i *Image) Version() (string) {
	//var err error = nil
	return ""
}

func buildBase(kind, path string) (r Base, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	r = Base{kind, name, absPath}
	return
}

func Init(path, kind string) (err error) {
	var b Base
	switch kind {
		case EnvKind:
		var r Env
		r, err = buildEnv(path)
		b = r.Base
		case ProjectKind:
		var r Project
		r, err = buildProject(path)
		b = r.Base
		case ImageKind:
		var r Image
		r, err = buildImage(path)
		b = r.Base
		default:
		err = fmt.Errorf("Unable to load Resource from path: %s ! Not supported kind property: [%s].", path, kind)
	}

	if err != nil {
                return
        }

	err = Store(b)

	config.Init(path, b)

	return
}

func buildEnv(path string) (r Env, err error) {
	base, err := buildBase(EnvKind, path)
	if err != nil {
                return
        }

	r = Env{Base: base}
	return
}

func buildProject(path string) (r Project, err error) {
	base, err := buildBase(ProjectKind, path)
	if err != nil {
                return
        }

	r = Project{Base: base}
	return
}

func buildImage(path string) (r Image, err error) {
	base, err := buildBase(ImageKind, path)
	if err != nil {
                return
        }

	r = Image{Base: base}
	return
}

func Load(path string) (r Resource, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	resourceFilepath := filepath.Join(path, defaultResourceFile)
	content, err := os.ReadFile(resourceFilepath)
	if err != nil {
		return
	}

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
		r = Env{Base: base}
		case ProjectKind:
		r = Project{Base: base}
		case ImageKind:
		r = Image{Base: base}
		default:
		err = fmt.Errorf("Unable to load Resource from path: %s ! Not supported kind property: [%s].", path, kind)
		return
	}

	err = yaml.Unmarshal(content, r)
	if err != nil {
		return
	}

	return
}

var storeLock = &sync.Mutex{}
func Store(r Base) (err error) {
	storeLock.Lock()
	defer storeLock.Unlock()

	content, err := yaml.Marshal(r)
        if err != nil {
                return
        }

	resourceFilepath := filepath.Join(r.Dir(), defaultResourceFile)
	err = os.WriteFile(resourceFilepath, content, 0644)

	return
}

func resourceName(path string) string {
	return filepath.Base(path)
}

