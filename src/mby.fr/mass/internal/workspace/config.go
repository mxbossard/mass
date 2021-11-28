package workspace

import (
	"fmt"
	"sync"
	"path/filepath"
	"io/ioutil"
	"log"
	"os"
	"errors"

	"gopkg.in/yaml.v2"
)

const settingsDir = ".mass"
const settingsFilePath = settingsDir + "/mass.yaml"

const configDir = "config"

// --- Settings ---

type Settings struct {
	Name string
	WorkspacePath string
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings name: %s ; workspacePath: %s.", s.Name, s.WorkspacePath)
}

func (s Settings) SettingsDirPath() string {
	return filepath.Join(s.WorkspacePath, settingsDir)
}

func (s Settings) SettingsFilePath() string {
	return filepath.Join(s.WorkspacePath, settingsFilePath)
}

func (s Settings) ConfigDirPath() string {
	return filepath.Join(s.WorkspacePath, configDir)
}

func (s Settings) store() {
	settingsFilePath := s.SettingsFilePath()
	data, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(s.SettingsDirPath(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(settingsFilePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getWorkDirPath() string {
	workDirPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return workDirPath
}

func InitSettings() {
	workspacePath := getWorkDirPath()
	name := filepath.Base(workspacePath)
	settings := Settings{name, workspacePath}

	_, err := os.Stat(settings.SettingsFilePath())
	if err == nil {
		// settings file already exists
		log.Fatal("Workspace settings already exists !")
	}

	settings.store()
}

func InitConfig() {
	workspacePath := getWorkDirPath()
	CreateNewDirectory(workspacePath, configDir)
}

func seekSettingsPathRecurse(dirPath string) (string, error) {
	//log.Printf("Seek Settings in dir: %s ...\n", dirPath)
	if dirPath == "/" {
		return "", errors.New("Unable to found settings path")
	}
	settingsFilePath := filepath.Join(dirPath, settingsFilePath)

	_, err := os.Stat(settingsFilePath)
	if err == nil {
		// settings file exists
		return settingsFilePath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// settings file does *not* exist

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	}

	parentDirPath := filepath.Dir(dirPath)
	return seekSettingsPathRecurse(parentDirPath)
}

func seekSettingsPath() string {
	workDirPath := getWorkDirPath()

	settingsPath, err := seekSettingsPathRecurse(workDirPath)

	if settingsPath != "" && err == nil {
		return settingsPath
	}
	return ""
}

func loadSettings() (*Settings, bool) {
	settingsFilePath := seekSettingsPath()
	if settingsFilePath == "" {
		// Workspace settings does not exists yet
		return nil, false
	}

	yfile, err := ioutil.ReadFile(settingsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	settings := Settings{}

	err = yaml.Unmarshal(yfile, &settings)
	if err != nil {
		log.Fatal(err)
	}

	return &settings, true
}

// --- SettingsService ---

type SettingsService struct {
	settings *Settings
}

// constructor
func newSettingsService() SettingsService {
	settings, ok := loadSettings()
	if !ok {
		log.Fatal("Unable to load workspace settings !")
	}
	settingsService := SettingsService{settings}
	return settingsService
}

// settings getter
func (s SettingsService) Settings() *Settings {
	return s.settings
}

// singleton
var lock = &sync.Mutex{}

var settingsService *SettingsService

func GetSettingsService() *SettingsService {
	if settingsService == nil {
		lock.Lock()
		defer lock.Unlock()
		if settingsService == nil {
			service := newSettingsService()
			settingsService = &service
		}
	}
	return settingsService
}


