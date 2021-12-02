package workspace

import (
)

func InitConfig() (err error) {
	if settingsService, err := GetSettingsService(); err == nil {
		err = CreateNewDirectory(settingsService.ConfigDirPath())
	}
	return 
}

