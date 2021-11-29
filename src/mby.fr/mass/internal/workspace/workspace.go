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

	settingsService := GetSettingsService()
	settingsService.InitSettings()

	InitConfig()

	for _, envName := range defaultEnvs {
		InitEnv(envName)
	}

	worspacePath := settingsService.WorkspacePath()

	fmt.Printf("New workspace initialized in %s\n", worspacePath)
}

