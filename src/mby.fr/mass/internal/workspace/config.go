package workspace

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

func InitConfig() (err error) {
	if settingsService, err := settings.GetSettingsService(); err == nil {
		err = file.CreateNewDirectory(settingsService.ConfigDir())
	}
	return 
}

