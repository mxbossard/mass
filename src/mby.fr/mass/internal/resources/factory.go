package resources

import (
	"fmt"
	"path/filepath"
)

/*
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
*/

func BuildAny(kind Kind, name string, parentResOrDir any) (res any, err error) {
	var parentDir string
	ok := false
	switch kind {
	case ImageKind:
		_, ok = parentResOrDir.(Resourcer)
		if !ok {
			err = fmt.Errorf("Expect a Resourcer parentRes to build kind: %v, received %T !", kind, parentResOrDir)
			return
		}
	default:
		parentDir, ok = parentResOrDir.(string)
		if !ok {
			err = fmt.Errorf("Expect a string parentDir to build kind: %v, received %T !", kind, parentResOrDir)
			return
		}
	}

	switch kind {
	case EnvKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildEnv(resDir)
	case ProjectKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildProject(resDir)
	case ImageKind:
		project, ok := parentResOrDir.(Project)
		if !ok {
			err = fmt.Errorf("Expect a P roject parentRes to build kind: %v, received %T !", kind, parentResOrDir)
			return
		}
		res, err = buildImage(project, name)
	case PodKind:
		res, err = buildPod(parentDir, name)
	default:
		err = fmt.Errorf("Unable to build Resource with name: %s in parent dir: %s ! Not supported kind property: [%s].", name, parentDir, kind)
	}

	return
}

func BuildResourcer(kind Kind, name string, parentResOrDir any) (res Resourcer, err error) {
	r, err := BuildAny(kind, name, parentResOrDir)
	if err != nil {
		return
	}
	res = r.(Resourcer)
	return
}

func Build[T Resourcer](name string, parentResOrDir any) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, name, parentResOrDir)
	if err != nil {
		return
	}
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
