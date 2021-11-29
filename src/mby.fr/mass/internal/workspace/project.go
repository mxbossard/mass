package workspace

import (

)

func InitProject(name string) {
	settingsService := GetSettingsService()

	CreateNewSubDirectory(settingsService.WorkspacePath(), name)
}

