package workspace

import (

)

func InitEnv(name string) {
	settingsService := GetSettingsService()

	CreateNewSubDirectory(settingsService.ConfigDirPath(), name)
}

