package workspace

import (
	//"os"
	//"io/fs"
	"strings"
	"path/filepath"

	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/resources"
	"mby.fr/utils/file"
)

var forbiddenNames = []string{resources.DefaultSourceDir, resources.DefaultTestDir, "envs"}

func InitProject(name string) (projectPath string, err error) {
	err = resources.AssertResourceName(resources.ProjectKind, name)
	if err != nil {
		return
	}

	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	projectPath = filepath.Join(settingsService.ProjectsDir(), name)

	_, err = resources.Init[resources.Project](projectPath)
	return
}

func InitImage(name string) (imagePath string, err error) {
	err = resources.AssertResourceName(resources.ImageKind, name)
	if err != nil {
		return
	}

	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	var projectName, imageName string

	splittedName := strings.Split(name, "/")
	if len(splittedName) == 1 {
		// Work dir must be project dir
		workDir, err := file.WorkDirPath()
		if err != nil {
			return "", err
		}
		projectName = filepath.Base(workDir)
		imageName = splittedName[0]
	} else if len(splittedName) == 2 {
		projectName = splittedName[0]
		imageName = splittedName[1]
	}

	err = AssertResourceExists(resources.ProjectKind, projectName)
	if err != nil {
		return
	}

	projectDir := filepath.Join(settingsService.ProjectsDir(), projectName)
	imagePath = filepath.Join(projectDir, imageName)
	_, err = resources.Init[resources.Image](imagePath)
	return
}

