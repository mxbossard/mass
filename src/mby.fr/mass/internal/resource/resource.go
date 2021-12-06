package resource

import(
	"fmt"
	"os"
	"sync"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const defaultResourceFile = "resource.yaml"

const EnvKind = "Env"
const ProjectKind = "Project"
const ImageKind = "Image"

type Resource interface {
	Kind() string
	Name() string
	Dir() string
}

type BaseResource struct {
	ResourceKind, name, dir string
}

func (r BaseResource) Kind() string {
	return r.ResourceKind
}

func (r BaseResource) Name() string {
	return r.name
}

func (r BaseResource) Dir() string {
	return r.dir
}

type EnvResource struct {
	BaseResource // Implicit composition: "golang inheritance"
}

type ProjectResource struct {
	BaseResource
}

type ImageResource struct {
	BaseResource
}

func buildBase(kind, path string) (r BaseResource, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	r = BaseResource{kind, name, absPath}
	return
}

func BuildEnv(path string) (r EnvResource, err error) {
	base, err := buildBase(EnvKind, path)
	if err != nil {
                return
        }

	r = EnvResource{base}
	return
}

func BuildProject(path string) (r ProjectResource, err error) {
	base, err := buildBase(ProjectKind, path)
	if err != nil {
                return
        }

	r = ProjectResource{base}
	return
}

func BuildImage(path string) (r ImageResource, err error) {
	base, err := buildBase(ImageKind, path)
	if err != nil {
                return
        }

	r = ImageResource{base}
	return
}

func LoadResource(path string) (r Resource, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	resourceFilepath := filepath.Join(path, defaultResourceFile)
	content, err := os.ReadFile(resourceFilepath)
	if err != nil {
		return
	}

	base := BaseResource{}
	err = yaml.Unmarshal(content, &base)
	if err != nil {
		return
	}

	base.name = filepath.Base(path)
	base.dir = path

	kind := base.Kind()
	
	switch kind {
		case EnvKind:
		r = &EnvResource{base}
		case ProjectKind:
		r = &ProjectResource{base}
		case ImageKind:
		r = &ImageResource{base}
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

var storeResourceLock = &sync.Mutex{}
func StoreResource(r BaseResource) (err error) {
	storeResourceLock.Lock()
	defer storeResourceLock.Unlock()

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

