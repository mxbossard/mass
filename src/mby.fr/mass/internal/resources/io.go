package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"mby.fr/utils/filez"
)

// TODO: move Marshal and Unmarshall to private methods in model.go

var writeLock = &sync.Mutex{}

func Write(r Resourcer) (err error) {
	writeLock.Lock()
	defer writeLock.Unlock()

	var content []byte
	switch res := r.(type) {
	case *Env, *Project, *Image, *Pod, *Service, *Endpoint:
		//fmt.Printf("Debug: resource pointer [%T] content: [%s] ...\n", res, res)
		content, err = yaml.Marshal(res)
	case base, Env, Project, Image, Pod, Service, Endpoint:
		//fmt.Printf("Debug: resource [%T] content: [%v] ...\n", res, res)
		content, err = yaml.Marshal(&res)
	default:
		err = fmt.Errorf("Unable to write Resource ! Not supported kind property: [%T].", r)
		return
	}

	if err != nil {
		return
	}

	err = os.MkdirAll(r.Dir(), 0755)
	if err != nil {
		return
	}
	resourceFilepath := filepath.Join(r.Dir(), DefaultResourceFile)
	//fmt.Printf("Debug: WRITING content: [%s] in file: [%s] ...\n", content, resourceFilepath)
	err = os.WriteFile(resourceFilepath, content, 0644)

	return
}

func ReadAny(path string) (r any, err error) {
	// If path point a dir => 	parent dir is the parent dir of the path dir.
	//							path is updated to point DefaultResourceFile in path dir
	// If path point a file => 	parent dir is the parent dir of the dir containing the file.
	var parentDir string
	if ok, _ := filez.IsDirectory(path); ok {
		parentDir = filepath.Dir(path)
		path = filepath.Join(path, DefaultResourceFile)
	} else {
		parentDir = filepath.Dir(filepath.Dir(path))
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			//log.Fatal(err)
			err = ResourceNotFound{path, NewKindSet(AllKind), err}
		}
		return
	}
	//fmt.Println("Debug: READING ResourceFile content:", string(content))

	base := base{}
	err = yaml.Unmarshal(content, &base)
	if err != nil {
		return
	}

	resName := filepath.Base(filepath.Dir(path))
	var parentResOrDir any = parentDir
	var parentProject Project
	kind := base.Kind()
	if kind == ImageKind {
		parentProjectResPath := filepath.Join(parentDir, DefaultResourceFile)
		parentProject, err = Read[Project](parentProjectResPath)
		if err != nil {
			return
		}
		parentResOrDir = parentProject
	} else if kind == ServiceKind {
		projectDir := filepath.Dir(path)
		parentProject, err = Read[Project](projectDir)
		if err != nil {
			return
		}
		parentResOrDir = parentProject
		resFilename := filepath.Base(path)
		resName, err = serviceNameFromResFilename(resFilename)
		if err != nil {
			return r, err
		}
	}

	
	res, err := BuildAny(kind, resName, parentResOrDir)
	if err != nil {
		return
	}
	//fmt.Printf("Build any: %T for kind: %s\n", res, kind)

	switch re := res.(type) {
	case Env:
		err = yaml.Unmarshal(content, &re)
		return re, err
	case Image:
		err = yaml.Unmarshal(content, &re)
		//FIXME: set Full Name ?
		return re, err
	case Project:
		err = yaml.Unmarshal(content, &re)
		return re, err
	//case Pod:
	//	//err = yaml.Unmarshal(content, &re)
	//	re, err = buildPod(parentDir, "badName")
	//	return re, err
	case Service:
		err = yaml.Unmarshal(content, &re)
		return re, err
	case Endpoint:
		err = yaml.Unmarshal(content, &re)
		return re, err
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
		if _, ok := err.(ResourceNotFound); ok {
			// If ResourceNotFound error add expected type in error
			kind := KindFromResource(r)
			err = ResourceNotFound{path, NewKindSet(kind), err}
		}
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
