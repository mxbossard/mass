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

const settingsDir = ".mass/"
const settingsFilePath = "./mass.yaml"

const configDir = "config"

// --- Settings ---

type Settings struct {
	Name string
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings name: %s", s.Name)
}

func InitConfig() {
	workspacePath := GetWorkDirPath()
	CreateNewSubDirectory(workspacePath, configDir)
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
	workDirPath := GetWorkDirPath()

	settingsPath, err := seekSettingsPathRecurse(workDirPath)
	if err != nil {
		log.Fatal(err)
	}
	return settingsPath
}

func loadSettings() (*Settings, string, bool) {
	settingsFilePath := seekSettingsPath()
	if settingsFilePath == "" {
		// Workspace settings does not exists yet
		return nil, "", false
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

	return &settings, settingsFilePath, true
}

// --- SettingsService ---

type SettingsService struct {
	workspacePath string
	settings *Settings
}

// constructor
func newSettingsService() SettingsService {
	settings, settingsFilePath, ok := loadSettings()
	if !ok {
		log.Fatal("Unable to load workspace settings !")
	}
	workspacePath := filepath.Dir(settingsFilePath)
	settingsService := SettingsService{settings: settings, workspacePath: workspacePath}
	return settingsService
}

// workspacePath getter
func (s SettingsService) WorkspacePath() string {
	return s.workspacePath
}

// settings getter
func (s SettingsService) Settings() *Settings {
	return s.settings
}

func (s SettingsService) SettingsDirPath() string {
	return filepath.Join(s.workspacePath, settingsDir)
}

func (s SettingsService) SettingsFilePath() string {
	return filepath.Join(s.workspacePath, settingsFilePath)
}

func (s SettingsService) ConfigDirPath() string {
	return filepath.Join(s.workspacePath, configDir)
}

func (s *SettingsService) InitSettings() {
	name := filepath.Base(s.workspacePath)
	s.settings = &Settings{Name: name}

	_, err := os.Stat(s.SettingsFilePath())
	if err == nil {
		// settings file already exists
		log.Fatal("Workspace settings already exists !")
	}

	s.store()
}

func (s SettingsService) store() {
	settingsFilePath := s.SettingsFilePath()
	data, err := yaml.Marshal(s.settings)
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


