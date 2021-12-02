package workspace

import (
	//"os"
	"fmt"
	//"path/filepath"
)

func init() {
  	//fmt.Println("This will get called on main initialization")
}

func Init(path string) (err error) {
	fmt.Println("foo")
	ok, err := isPathInExistingWorkspace(path)
	if err != nil {
		return
	} else if ok {
		err = fmt.Errorf("Supplied path is in already existing workspace !")
		return
	}
	fmt.Println("bar")

	err = CreateDirectory(path)
	if err != nil {
		return
	}

	err = Chdir(path)
	if err != nil {
		return
	}

	workspacePath, err := WorkDirPath()
	if err != nil {
		return
	}

	err = InitSettings(workspacePath)
	if err != nil {
		return
	}

	err = InitConfig()
	if err != nil {
		return
	}

	settingsService, err := GetSettingsService()
	if err != nil {
		return
	}

	settings := settingsService.Settings()

	for _, envName := range settings.Environments {
		err = InitEnv(envName)
		if err != nil {
			return
		}
	}

	fmt.Printf("New workspace initialized in %s\n", workspacePath)
	return
}

func isPathInExistingWorkspace(path string) (ok bool, err error) {
	// Search for settings already present in target path
	settingsFilePath, err := seekSettingsFilePath(path)
	if settingsFilePath != "" {
		ok = true
	}
	return
}

