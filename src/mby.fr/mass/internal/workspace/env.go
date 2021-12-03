package workspace

import (

)

func InitEnv(name string) (err error) {
	settingsService, err := GetSettingsService()
	if err != nil {
		return
	}

	_, err = CreateNewSubDirectory(settingsService.ConfigDirPath(), name)
	return
}

