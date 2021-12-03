package workspace

import (

)

var sourcesDir = "src"

func InitProject(name string) (err error) {
	settingsService, err := GetSettingsService()
	if err != nil {
		return
	}

	projectPath, err := CreateNewSubDirectory(settingsService.WorkspacePath(), name)
	if err != nil {
		return
	}

	_, err = CreateNewSubDirectory(projectPath, sourcesDir)
	return
}

