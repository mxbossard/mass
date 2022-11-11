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

func buildBase(kind Kind, path string) (b base, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	name := resourceName(path)
	b = base{ResourceKind: kind, name: name, dir: absPath}
	return
}

func buildTestable(res Resourcer, path string) (t testable, err error) {
	testDir := DefaultTestDir
	/*
		base, err := buildBase(kind, path)
		if err != nil {
			return
		}
		t = Testable{base, testDir}
	*/
	t = testable{resource: res, testDirectory: testDir}
	return
}

func BuildEnv(path string) (r Env, err error) {
	base, err := buildBase(EnvKind, path)
	if err != nil {
		return
	}

	r = Env{base: base}
	return
}

func BuildProject(path string) (p Project, err error) {
	deployfile := DefaultDeployFile

	b, err := buildBase(ProjectKind, path)
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

func BuildImage(path string) (r Image, err error) {
	version := DefaultInitialVersion
	buildfile := DefaultBuildFile
	sourceDir := DefaultSourceDir

	projectPath := filepath.Dir(path)
	project, err := BuildProject(projectPath)
	if err != nil {
		return
	}

	b, err := buildBase(ImageKind, path)
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
