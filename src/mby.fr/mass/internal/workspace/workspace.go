package workspace

import (
	"log"
	"os"
	"fmt"
)

var defaultEnvs = []string{"dev", "stage", "prod"}

func Init(path string) {
	CreateNewDirectory(path)

        err := os.Chdir(path)
        if (err != nil) {
                log.Fatal(err)
        }

	InitSettings()
	InitConfig()

	for _, envName := range defaultEnvs {
		InitEnv(envName)
	}

	settingsService := GetSettingsService()
        settings := settingsService.Settings()
	worspacePath := settings.WorkspacePath

	fmt.Printf("New workspace initialized in %s\n", worspacePath)
}

