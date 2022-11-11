package resources

import (
	"fmt"
)

type Res interface {
	Env | Project | Image
}

/*
type ResPtr interface {
	*Env | *Project | *Image
}
*/

func BuildAny(kind Kind, baseDir string) (res any, err error) {
	switch kind {
	case EnvKind:
		res, err = BuildEnv(baseDir)
	case ProjectKind:
		res, err = BuildProject(baseDir)
	case ImageKind:
		res, err = BuildImage(baseDir)
	default:
		err = fmt.Errorf("Unable to build Resource with base dir: %s ! Not supported kind property: [%s].", baseDir, kind)
	}

	return
}

func BuildResourcer(kind Kind, baseDir string) (res Resourcer, err error) {
	r, err := BuildAny(kind, baseDir)
	res = r.(Resourcer)
	return
}

func Build[T Resourcer](path string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, path)
	r = (any)(res).(T)
	return
}

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
type ResourceCallFunc[T Resourcer, K any] = func(T) (K, error)

func CallFuncOnResource[K any](r T, f ResourceCallFunc[T, K]) (a K, err error) {
	switch res := r.(type) {
	case Env:
		return f(r)
	case Image:
		return f(r)
	case Project:
		return f(r)
	default:
		err = fmt.Errorf("Type %T not supported yet !", res)
		return
	}
}
*/