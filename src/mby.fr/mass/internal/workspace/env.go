package workspace

import (
	"path/filepath"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/filez"
)

func InitEnvs() (err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}

	err = filez.CreateNewDirectory(settingsService.EnvsDir())
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

func InitEnv(name string) (envPath string, err error) {
	settingsService, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	envPath = filepath.Join(settingsService.EnvsDir(), name)
	_, err = resources.Init[resources.Env](envPath)
	return
}
