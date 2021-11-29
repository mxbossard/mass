package workspace

import (

)

func InitEnv(name string) {
	settingsService := GetSettingsService()
        settings := settingsService.Settings()

	CreateNewSubDirectory(settings.ConfigDirPath(), name)
}

