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

func assertName(k Kind, name string) (err error) {
	// TODO
	return
}

func assertFullName(k Kind, fullName string) (err error) {
	parts := strings.Split(fullName, fullNameSeparator)
	switch k {
	case ProjectKind, EnvKind:
		if len(parts) != 1 {
			err = fmt.Errorf("Malformed fullName [%s] kind: [%s].", fullName, k)
		}
	case ImageKind, PodKind, ServiceKind, EndpointKind:
		if len(parts) != 2 {
			err = fmt.Errorf("Malformed fullName [%s] kind: [%s].", fullName, k)
		}
	default:
		err = fmt.Errorf("Unable to assert resource fullName for kind: [%s].", k)
	}
	return
}

func splitResourceHierarchy(k Kind, fullName string) (hierarchy []Uid, err error) {
	err = assertFullName(k, fullName)
	if err != nil {
		return
	}

	switch k {
	case ProjectKind, EnvKind:
		append(hierarchy, Uid{k, fullName})
	case ImageKind, PodKind, ServiceKind, EndpointKind:
		// Take project part of fullName
		parts := strings.Split(fullName, fullNameSeparator)
		append(hierarchy, SplitResourceHierarchy(ProjectKind, parts[0]))
		append(hierarchy, Uid{k, fullName})
	default:
		err = fmt.Errorf("Unable to split resource hierarchy for kind: [%s].", k)
	}
	return
}

/*
func forgeResourceFilepath(k Kind, parentDir, name string) (path string, err error) {
	err = assertName(k, fullName)
	if err != nil {
		return
	}

	//hierarchy, err := splitResourceHierarchy(k, fillName)

	switch k {
	case ProjectKind, EnvKind, ImageKind:
		// name not used
		path = filepath.Join(parentDir, DefaultResourceFile)
	//case ImageKind:
	//	path = filepath.Join(rootDir, fullName, DefaultResourceFile)
	case PodKind, ServiceKind, EndpointKind:
		resourceFile := fmt.Sprintf("%s-%s.yaml", k, name)
		path = filepath.Join(parentDir, resourceFile)
	default:
		err = fmt.Errorf("Unable to forge resource filepath for kind: [%s].", k)
	}
	return
}
*/