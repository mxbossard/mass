package workspace

import (
	//"os"
	"fmt"
	//"path/filepath"

	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

func init() {
  	//fmt.Println("This will get called on main initialization")
}

func Init(path string) (err error) {
	ok, err := isPathInExistingWorkspace(path)
	if err != nil {
		return
	} else if ok {
		err = fmt.Errorf("Supplied path is in already existing workspace !")
		return
	}

	err = file.CreateDirectory(path)
	if err != nil {
		return
	}

	err = file.Chdir(path)
	if err != nil {
		return
	}

	workspacePath, err := file.WorkDirPath()
	if err != nil {
		return
	}

	err = settings.InitSettings(workspacePath)
	if err != nil {
		return
	}

	err = InitConfig()
	if err != nil {
		return
	}

	settingsService, err := settings.GetSettingsService()
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

	fmt.Printf("New workspace initialized in %s\n", workspacePath)
	return
}

func isPathInExistingWorkspace(path string) (ok bool, err error) {
	// Search for settings already present in target path
	settingsFilePath, err := settings.SeekSettingsFilePath(path)
	if settingsFilePath != "" {
		ok = true
	}
	return
}

