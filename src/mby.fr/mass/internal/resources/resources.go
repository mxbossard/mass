package resources

import (
	//"fmt"
	"path/filepath"
	"strings"
)

const DefaultSourceDir = "src"
const DefaultTestDir = "test"
const DefaultVersionFile = "version.txt"
const DefaultInitialVersion = "0.0.1-dev"
const DefaultBuildFile = "Dockerfile"
const DefaultDeployFile = "compose.yaml"
const DefaultResourceFile = "resource.yaml"

func InitResourcer(kind Kind, path string) (res Resourcer, err error) {
	res, err = BuildResourcer(kind, path)
	if err != nil {
		return
	}
	err = res.init()
	if err != nil {
		return
	}
	err = Write(res)
	return
}

func Init[T Resourcer](path string) (r T, err error) {
	r, err = Build[T](path)
	if err != nil {
		return
	}
	err = r.init()
	if err != nil {
		return
	}
	err = Write(r)
	return
}

/*
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
*/

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
func absResourcePath(relRootPath string, resPath string) (path string) {
	path = filepath.Join(relRootPath, resPath)
	return
}
