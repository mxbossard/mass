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

func FromPath[T Resourcer](path string) (res T, err error) {
	r, err := Read(path)
	if err != nil {
		return
	}

	res, ok := r.(T)
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		// Expect resource value
		// In this case, r is a pointer and we want to return a value, but the type cast don't return ok.
		resPtrType := reflect.PointerTo(reflect.TypeOf(res))
		if reflect.TypeOf(r) == resPtrType {
			// Right type so res was rightly type cast
			return res, err
		}
	}

	if !ok {
		err = fmt.Errorf("bad resource type for path %s. Expected type %T but got %T", path, res, r)
	}
	return res, err
}

// Return pointer to resource
//func FromKind[T Resourcer](kind Kind, baseDir string) (res T, err error) {
func FromKind(kind Kind, baseDir string) (res Resourcer, err error) {
	switch kind {
	case EnvKind:
		r, err := BuildEnv(baseDir)
		if err != nil {
			return res, err
		}
		res = &r
	case ProjectKind:
		r, err := BuildProject(baseDir)
		if err != nil {
			return res, err
		}
		res = &r
	case ImageKind:
		r, err := BuildImage(baseDir)
		if err != nil {
			return res, err
		}
		res = &r
	default:
		err = fmt.Errorf("Unable to build Resource with base dir: %s ! Not supported kind property: [%s].", baseDir, kind)
		return
	}
	
	return
}

func Init(kind Kind, path string) (res Resourcer, err error) {
	res, err = FromKind(kind, path)
	if err != nil {
		return
	}
	res.Init()

	
	return
}

func Init2[T Resourcer](kind Kind, path string) (res T, err error) {
	_, err = Init(kind, path)
	if err != nil {
		return
	}
	res, err = FromPath[T](path)
	return
}

func Init3[T Resourcer](path string) (res T, err error) {
	var kind Kind
	switch (interface{})(res).(type) {
	case Env, *Env:
		kind = EnvKind
	case Project, *Project:
		kind = ProjectKind
	case Image, *Image:
		kind = ImageKind
	}
	_, err = Init(kind, path)
	if err != nil {
		return
	}
	res, err = FromPath[T](path)
	return
}