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

func BuildAny(kind Kind, parentDir, name string) (res any, err error) {
	
	switch kind {
	case EnvKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildEnv(resDir)
	case ProjectKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildProject(resDir)
	case ImageKind:
		res, err = buildImage(parentDir, name)
	case PodKind:
		res, err = buildPod(parentDir, name)
	default:
		err = fmt.Errorf("Unable to build Resource with name: %s in parent dir: %s ! Not supported kind property: [%s].", name, parentDir, kind)
	}

	return
}

func BuildResourcer(kind Kind, parentDir, name string) (res Resourcer, err error) {
	r, err := BuildAny(kind, parentDir, name)
	res = r.(Resourcer)
	return
}

func Build[T Resourcer](parentDir, name string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildResourcer(kind, parentDir, name)
	r = (any)(res).(T)
	return
}
/*
func BuildAnyPrimary(kind Kind, parentDir, name string) (res any, err error) {
	
	switch kind {
	case EnvKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildEnv(resDir)
	case ProjectKind:
		resDir := filepath.Join(parentDir, name)
		res, err = buildProject(resDir)
	default:
		err = fmt.Errorf("Unable to build Resource with name: %s in parent dir: %s ! Not supported kind property: [%s].", name, parentDir, kind)
	}

	return
}

func BuildAnyChild(kind Kind, parent Resourcer, name string) (res any, err error) {
	switch kind {
	case ImageKind:
		if project, ok := parent.(Project); ok {
			res, err = buildImage(project, name)
		} else {
			err = fmt.Errorf("Image parent resource must be a project, not %T !", parent)	
		}
	case PodKind:
		if project, ok := parent.(Project); ok {
			res, err = buildPod(project, name)
		} else {
			err = fmt.Errorf("Image parent resource must be a project, not %T !", parent)	
		}
	default:
		err = fmt.Errorf("Unable to build Resource with name: %s in parent dir: %s ! Not supported kind property: [%s].", name, parentDir, kind)
	}

	return
}

func BuildPrimaryResourcer(kind Kind, parentDir, name string) (res Resourcer, err error) {
	r, err := BuildAnyPrimary(kind, parentDir, name)
	res = r.(Resourcer)
	return
}

func BuildChildResourcer(kind Kind, parent Resourcer, name string) (res Resourcer, err error) {
	r, err := BuildAnyChild(kind, parent, name)
	res = r.(Resourcer)
	return
}

func BuildPrimary[T Resourcer](parentDir, name string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildPrimaryResourcer(kind, parentDir, name)
	r = (any)(res).(T)
	return
}

func BuildChild[T Resourcer](parent Resourcer, name string) (r T, err error) {
	kind := KindFromResource(r)
	res, err := BuildChildResourcer(kind, parent, name)
	r = (any)(res).(T)
	return
}
*/

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
