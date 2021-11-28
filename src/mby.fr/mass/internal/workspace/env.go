package workspace

import (

)

func InitEnv(name string) {
	settingsService := GetSettingsService()
        settings := settingsService.Settings()

	CreateNewDirectory(settings.ConfigDirPath(), name)
}

