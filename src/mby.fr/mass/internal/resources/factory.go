package resources

import (
	"fmt"
)

func BuildAny0(kind Kind, baseDir string) (res any, err error) {
	switch kind {
	case EnvKind:
		res, err = buildEnv(baseDir)
	case ProjectKind:
		res, err = buildProject(baseDir)
	case ImageKind:
		res, err = buildImage(baseDir)
	case PodKind:
		res, err = buildPod(baseDir, "")
	default:
		err = fmt.Errorf("Unable to build Resource with base dir: %s ! Not supported kind property: [%s].", baseDir, kind)
	}

	return
}

func BuildResourcer0(kind Kind, baseDir string) (res Resourcer, err error) {
	r, err := BuildAny(kind, baseDir)
	res = r.(Resourcer)
	return
}

func Build0[T Resourcer](path string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, path)
	r = (any)(res).(T)
	return
}

func BuildAny(kind Kind, name string, parent Resourcer) (res any, err error) {
	switch kind {
	case EnvKind:
		res, err = buildEnv(name)
	case ProjectKind:
		res, err = buildProject(name)
	case ImageKind:
		res, err = buildImage(name, parent)
	case PodKind:
		res, err = buildPod(name, parent)
	default:
		err = fmt.Errorf("Unable to build Resource with base dir: %s ! Not supported kind property: [%s].", baseDir, kind)
	}

	return
}

func BuildResourcer(kind Kind, name string, parent Resourcer) (res Resourcer, err error) {
	r, err := BuildAny(kind, name, parent)
	res = r.(Resourcer)
	return
}

func Build[T Resourcer](name string, parent Resourcer) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, name, parent)
	r = (any)(res).(T)
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
