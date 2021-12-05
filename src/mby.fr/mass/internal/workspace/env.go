package workspace

import (

)

func InitEnv(name string) (path string, err error) {
	settingsService, err := GetSettingsService()
	if err != nil {
		return
	}

	path, err = CreateSubDirectory(settingsService.ConfigDir(), name)
	return
}

