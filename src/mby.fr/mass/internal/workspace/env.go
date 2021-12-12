package workspace

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/resources"
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

func InitEnv(name string) (path string, err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	path, err = file.CreateSubDirectory(settingsService.EnvsDir(), name)
	if err != nil {
		return
	}

	err = resources.Init(path, resources.EnvKind)

	return
}

