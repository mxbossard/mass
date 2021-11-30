package workspace

import (
	"log"
	"fmt"
	"path/filepath"
)

func Init(path string) {
	assertPathNotInExistingWorkspace(path)

	CreateNewDirectory(path)

        Chdir(path)

	workspacePath := GetWorkDirPath()
	initSettings(workspacePath)

	InitConfig()

	settingsService := GetSettingsService()
	settings := settingsService.Settings()

	for _, envName := range settings.Environments {
		InitEnv(envName)
	}

	fmt.Printf("New workspace initialized in %s\n", workspacePath)
}

func assertPathNotInExistingWorkspace(path string) {
	workPath := GetWorkDirPath()
	// Search for settings already present in target parent path
	parentPath := filepath.Dir(path)
	Chdir(parentPath)
        settingsFilePath, _ := seekSettingsFilePath()
	if settingsFilePath != "" {
		log.Fatal("Cannot init a workspace inside another workspace !")
	}
	Chdir(workPath)
}
