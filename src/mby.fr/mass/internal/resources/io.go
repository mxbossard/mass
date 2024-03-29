package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// TODO: move Marshal and Unmarshall to private methods in model.go

var writeLock = &sync.Mutex{}

func Write(r Resourcer) (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()

	var content []byte
	switch res := r.(type) {
	case *Env, *Project, *Image:
		//fmt.Printf("Debug: resource pointer [%T] content: [%s] ...\n", res, res)
		content, err = yaml.Marshal(res)
	case Env, Project, Image, base:
		//fmt.Printf("Debug: resource [%T] content: [%v] ...\n", res, res)
		content, err = yaml.Marshal(&res)
	default:
		err = fmt.Errorf("Unable to write Resource ! Not supported kind property: [%T].", r)
		return
	}

	if err != nil {
		return
	}

	resourceFilepath := filepath.Join(r.Dir(), DefaultResourceFile)
	//fmt.Printf("Debug: WRITING content: [%s] in file: [%s] ...\n", content, resourceFilepath)
	err = os.WriteFile(resourceFilepath, content, 0644)

	return
}

func ReadAny(path string) (r any, err error) {
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
	//fmt.Println("Debug: READING ResourceFile content:", string(content))

	base := base{}
	err = yaml.Unmarshal(content, &base)
	if err != nil {
		return
	}

	kind := base.Kind()
	res, err := BuildAny(kind, path)
	if err != nil {
		return
	}
	//fmt.Printf("Build any: %T for kind: %s\n", res, kind)

	switch re := res.(type) {
	case Env:
		err = yaml.Unmarshal(content, &re)
		return re, nil
	case Image:
		err = yaml.Unmarshal(content, &re)
		return re, nil
	case Project:
		err = yaml.Unmarshal(content, &re)
		return re, nil
	}

	err = fmt.Errorf("Unable to read Resource in file [%s] ! Not supported kind property: [%T].", path, res)
	return
}

func ReadResourcer(path string) (res Resourcer, err error) {
	r, err := ReadAny(path)
	if r != nil {
		res = r.(Resourcer)
	}
	return
}

func Read[T Resourcer](path string) (r T, err error) {
	res, err := ReadResourcer(path)
	if err != nil {
		return
	}

	r, ok := res.(T)
	if !ok {
		//err = fmt.Errorf("bad resource type for path %s. Expected type %T but got %T", path, r, res)
		err = BadResourceType{path, r, res}
		return
	}

	return
}
