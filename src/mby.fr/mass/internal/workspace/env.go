package workspace

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

func InitEnv(name string) (path string, err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	path, err = file.CreateSubDirectory(settingsService.ConfigDir(), name)
	return
}

