package resources

import (
	"fmt"
	"path/filepath"
)

func buildBase(kind Kind, dirPath, backingFilename string) (b base, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return
	}
	name := resourceName(dirPath)
	b = base{ResourceKind: kind, name: name, dir: absDirPath, backingFilename: backingFilename}
	return
}

func buildTestable(res Resourcer, path string) (t testable, err error) {
	testDir := DefaultTestDir
	t = testable{resource: res, testDirectory: testDir}
	return
}

func buildEnv(path string) (r Env, err error) {
	base, err := buildBase(EnvKind, path, DefaultResourceFile)
	if err != nil {
		return
	}

	r = Env{base: base}
	return
}

func buildProject(path string) (p Project, err error) {
	deployfile := DefaultDeployFile

	b, err := buildBase(ProjectKind, path, DefaultResourceFile)
	if err != nil {
		return
	}
	p = Project{base: b, DeployFile: deployfile}

	t, err := buildTestable(b, path)
	if err != nil {
		return
	}
	p.testable = t

	return
}

func buildImage(path string) (r Image, err error) {
	version := DefaultInitialVersion
	buildfile := DefaultBuildFile
	sourceDir := DefaultSourceDir

	projectPath := filepath.Dir(path)
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}

	b, err := buildBase(ImageKind, path, DefaultResourceFile)
	if err != nil {
		return
	}

	r = Image{
		base:            b,
		BuildFile:       buildfile,
		SourceDirectory: sourceDir,
		Project:         project,
	}

	t, err := buildTestable(r, path)
	if err != nil {
		return
	}
	r.testable = t

	versionable := buildVersionable(r, version)
	r.versionable = versionable

	return
}

func buildPod(projectPath, name string) (r Pod, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}

	backingFilename := fmt.Sprintf("pod-%s.yaml", name)
	b, err := buildBase(PodKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Pod{
		base:            b,
		Project:         &project,
	}

	t, err := buildTestable(r, projectPath)
	if err != nil {
		return
	}
	r.testable = t

	return
}


func buildEndpoint(projectPath, name string) (r Endpoint, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}
	backingFilename := fmt.Sprintf("end-%s.yaml", name)
	base, err := buildBase(EndpointKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Endpoint{base: base, Project: &project}
	return
}


func buildService(projectPath, name string) (r Service, err error) {
	project, err := buildProject(projectPath)
	if err != nil {
		return
	}
	backingFilename := fmt.Sprintf("svc-%s.yaml", name)
	base, err := buildBase(ServiceKind, projectPath, backingFilename)
	if err != nil {
		return
	}

	r = Service{base: base, Project: &project}
	return
}

func BuildAny(kind Kind, baseDir string) (res any, err error) {
	switch kind {
	case EnvKind:
		res, err = buildEnv(baseDir)
	case ProjectKind:
		res, err = buildProject(baseDir)
	case ImageKind:
		res, err = buildImage(baseDir)
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