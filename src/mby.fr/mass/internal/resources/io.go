package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
	//"github.com/ghodss/yaml"
)

var writeLock = &sync.Mutex{}

func Write(r Resourcer) (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()

	var content []byte
	switch res := r.(type) {
	case *Env, *Project, *Image:
		fmt.Printf("Debug: resource pointer [%T] content: [%s] ...\n", res, res)
		content, err = yaml.Marshal(res)
	case Env, Project, Image, base:
		fmt.Printf("Debug: resource [%T] content: [%v] ...\n", res, res)
		content, err = yaml.Marshal(&res)
	default:
		err = fmt.Errorf("Unable to write Resource ! Not supported kind property: [%T].", r)
		return
	}

	if err != nil {
		return
	}

	resourceFilepath := filepath.Join(r.Dir(), DefaultResourceFile)
	fmt.Printf("Debug: WRITING content: [%s] in file: [%s] ...\n", content, resourceFilepath)
	err = os.WriteFile(resourceFilepath, content, 0644)

	return
}

func Read(path string) (r Resourcer, err error) {
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
	fmt.Println("Debug: READING ResourceFile content:", string(content))

	base := base{}
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