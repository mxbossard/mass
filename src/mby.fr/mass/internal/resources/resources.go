package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"reflect"

	"gopkg.in/yaml.v2"
)

const DefaultSourceDir = "src"
const DefaultTestDir = "test"
const DefaultVersionFile = "version.txt"
const DefaultInitialVersion = "0.0.1-dev"
const DefaultBuildFile = "Dockerfile"
const DefaultDeployFile = "compose.yaml"
const DefaultResourceFile = "resource.yaml"

func buildBase(kind Kind, path string) (b Base, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	b = Base{kind, name, absPath}
	return
}

func buildTestable(res Resource, path string) (t Testable, err error) {
	testDir := DefaultTestDir
	/*
	base, err := buildBase(kind, path)
	if err != nil {
		return
	}
	t = Testable{base, testDir}
	*/
	t = Testable{Resource: res, testDirectory: testDir}
	return
}

func buildVersionable(version string) (v Versionable) {
	v = Versionable{version}
	return
}

func Init(path string, kind Kind) (r Resourcer, err error) {
	switch kind {
	case EnvKind:
		var r Env
		r, err = BuildEnv(path)
		r.Init()
	case ProjectKind:
		var p Project
		r, err = BuildProject(path)
		p.Init()
	case ImageKind:
		var i Image
		r, err = BuildImage(path)
		i.Init()
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

func BuildProject(path string) (p Project, err error) {
	deployfile := DefaultDeployFile
	
	b, err := buildBase(ProjectKind, path)
	if err != nil {
		return
	}
	t, err := buildTestable(b, path)
	if err != nil {
		return
	}
	p = Project{Resourcer: t, DeployFile: deployfile}
	return
}

func BuildImage(path string) (r Image, err error) {
	version := DefaultInitialVersion
	buildfile := DefaultBuildFile
	sourceDir := DefaultSourceDir

	projectPath := filepath.Dir(path)
	project, err := BuildProject(projectPath)
	if err != nil {
		return
	}

	versionable := buildVersionable(version)

	b, err := buildBase(ImageKind, path)
	if err != nil {
		return
	}
	t, err := buildTestable(b, path)
	if err != nil {
		return
	}

	r = Image{
		Resourcer: 			t,
		Versionable:     	versionable,
		BuildFile:       	buildfile,
		SourceDirectory: 	sourceDir,
		Project:         	project,
	}

	return
}

func Undecorate[T any](o any, t T) (r T, ok bool) {
	r, ok = o.(T)
	if ok {
		return
	}
	metaValue := reflect.ValueOf(o)
	field := metaValue.FieldByName("Resourcer")
	if field != (reflect.Value{}) {
		//fmt.Printf("recursion on %T ...\n", field.Interface())
		return Undecorate(field.Interface(), t)
	}
	return r, false
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
		r = &res
	case ProjectKind:
		res, err := BuildProject(base.Dir())
		if err != nil {
			return r, err
		}
		//res := Project{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = &res
	case ImageKind:
		res, err := BuildImage(base.Dir())
		if err != nil {
			return r, err
		}
		//res := Image{Base: base}
		err = yaml.Unmarshal(content, &res)
		r = &res
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
	case Env, *Env, Project, *Project, Image, *Image:
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
