package workspace

import (

)

func InitProject(name string) (err error) {
	if settingsService, err := GetSettingsService(); err == nil {
		err = CreateNewSubDirectory(settingsService.WorkspacePath(), name)
	}
	return
}

