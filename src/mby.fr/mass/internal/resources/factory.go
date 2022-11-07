package resources

import (
	"fmt"
	"reflect"
)

/*
type Res interface {
	Env | Project | Image
}

type ResPtr interface {
	*Env | *Project | *Image
}
*/

func Build[T Resourcer](path string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, path)
	r = (any)(res).(T)
	return
}

func Init[T Resourcer](path string) (r T, err error) {
	r, err = Build[T](path)
	r.Init()
	return
}

func FromPath[T Resourcer](path string) (r T, err error) {
	res, err := Read(path)
	if err != nil {
		return
	}

	r, ok := res.(T)
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		// Expect resource value
		// In this case, r is a pointer and we want to return a value, but the type cast don't return ok.
		resPtrType := reflect.PointerTo(reflect.TypeOf(res))
		if reflect.TypeOf(r) == resPtrType {
			// Right type so res was rightly type cast
			return r, err
		}
	}

	if !ok {
		err = fmt.Errorf("bad resource type for path %s. Expected type %T but got %T", path, r, res)
	}
	return r, err
}

// Return pointer to resource
func BuildResourcer(kind Kind, baseDir string) (res Resourcer, err error) {
	switch kind {
	case EnvKind:
		r, err := BuildEnv(baseDir)
		if err != nil {
			return res, err
		}
		res = r
	case ProjectKind:
		r, err := BuildProject(baseDir)
		if err != nil {
			return res, err
		}
		res = r
	case ImageKind:
		r, err := BuildImage(baseDir)
		if err != nil {
			return res, err
		}
		res = r
	default:
		err = fmt.Errorf("Unable to build Resource with base dir: %s ! Not supported kind property: [%s].", baseDir, kind)
		return
	}

	return
}

func InitResourcer(kind Kind, path string) (res Resourcer, err error) {
	res, err = BuildResourcer(kind, path)
	if err != nil {
		return
	}
	res.Init()

	return
}

func FromPathResourcer(path string) (res Resourcer, err error) {
	res, err = Read(path)
	return
}
