package workspace

import (

)

func InitProject(name string) {
	settingsService := GetSettingsService()
        settings := settingsService.Settings()

	CreateNewSubDirectory(settings.WorkspacePath, name)
}

