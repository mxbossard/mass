package project

import (
	"os"
	"path/filepath"

	"mby.fr/mass/internal/workspace"
)

const defaultImageSourcesDir = "src"
const defaultTestDir = "test"
const defaultVersionFile = "version.txt"
const defaultInitialVersion = "0.0.1"

var forbiddenNames = []string{defaultImageSourcesDir, defaultTestDir, "config"}

func InitProject(name string) (projectPath string, err error) {
	settingsService, err := workspace.GetSettingsService()
	if err != nil {
		return
	}

	// Create project dir
	projectPath, err = workspace.CreateSubDirectory(settingsService.ProjectsDir(), name)
	if err != nil {
		return
	}

	// Create test dir
	path := testDirpath(projectPath)
	err = workspace.CreateDirectory(path)
	if err != nil {
		return
	}

	// Create version file
	path = versionFilepath(projectPath)
	_, err = os.Stat(path); 
	if os.IsNotExist(err) {
		// Do not overwrite version file if it already exists
		err = os.WriteFile(path, []byte(defaultInitialVersion), 0644)
	} else if err != nil {
		return
	}

	return
}

func InitImage(projectName, name string) (imagePath string, err error) {
	return
}

type Image struct {
	Name string
	Dir string
	TestDir string
	Version string
}

type Project struct {
	Name string
	Dir string
	TestDir string
	Version string
	Images []Image

}

func testDirpath(parentPath string) string {
	return filepath.Join(parentPath, defaultTestDir)
}

func versionFilepath(parentPath string) string {
	return filepath.Join(parentPath, defaultVersionFile)
}

func buildProject(path string) (p Project, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	name := filepath.Base(path)

	testDir := testDirpath(path)
	_, err = os.Stat(testDir); 
	if os.IsNotExist(err) {
		testDir = ""
	} else if err != nil {
		return
	}

	versionFile := versionFilepath(path)
	content, err := os.ReadFile(versionFile); 
	version := string(content)
	if os.IsNotExist(err) {
		version = ""
	} else if err != nil {
		return
	}

	p = Project{Name: name, Dir: path, TestDir: testDir, Version: version}
	return
}

func ListProjects() (projects []Project, err error) {
	// List directories with a version file to build project list
        settingsService, err := workspace.GetSettingsService()
        if err != nil {
                return
        }

	projectsDir := settingsService.ProjectsDir()
	dirEntries, err := os.ReadDir(projectsDir)
	if err != nil {
		return
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			dirpath := filepath.Join(projectsDir, dirEntry.Name())
			versionFile := versionFilepath(dirpath)
			_, err = os.Stat(versionFile);
        		if os.IsNotExist(err) {
				err = nil
				// Version file does not exists => not a project
				continue
			} else if err != nil {
				return
			}
			// Found a project
			p, err := buildProject(dirpath)
			if err != nil {
				return projects, err
			}
			projects = append(projects, p)
		}
	}

	return
}

func GetProject(name string) (p Project, ok bool, err error) {
	projects, err := ListProjects()
	for _, p = range projects {
		if p.Name == name {
			ok = true
			return
		}
	}
	return
	
}

