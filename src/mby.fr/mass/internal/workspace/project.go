package workspace

import (

)

func InitProject(name string) {
	settingsService := GetSettingsService()
        settings := settingsService.Settings()

	CreateNewDirectory(settings.WorkspacePath, name)
}

