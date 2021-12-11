package workspace

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

func InitEnvs() (err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	err = file.CreateNewDirectory(settingsService.EnvsDir())
	if err != nil {
		return
	}

	settings := settingsService.Settings()

	for _, envName := range settings.Environments {
		_, err = InitEnv(envName)
		if err != nil {
			return
		}
	}

	return 
}

