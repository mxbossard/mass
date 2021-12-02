package workspace

import (

)

func InitEnv(name string) (err error) {
	if settingsService, err := GetSettingsService(); err == nil {
		err = CreateNewSubDirectory(settingsService.ConfigDirPath(), name)
	}
	return
}

