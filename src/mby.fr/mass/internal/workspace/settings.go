package workspace

import (
	"fmt"
	"sync"
	"path/filepath"
	"log"
	"os"
	"errors"

	"github.com/spf13/viper"
)

const settingsDir = ".mass/"
const settingsFilePath = settingsDir + "settings.yaml"

// --- Settings ---

func initViper(workspacePath string) {
	settingsDirPath := filepath.Join(workspacePath, settingsDir)
	workspaceName := filepath.Base(workspacePath)

	//log.Println("Initialize viper ...", workspaceName, settingsDirPath)

	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(settingsDirPath)

	viper.SetDefault("Name", workspaceName)
	viper.SetDefault("ConfigDir", "config")
	viper.SetDefault("Environments", []string{"dev", "stage", "prod"})
}

// Store settings erasing previous settings
func storeSettings() {
	//log.Println("Store settings ...")
	err := viper.WriteConfig()
	if err != nil {
		log.Fatal("Unable to store settings !", err)
	}
}

func readSettings() (Settings) {
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			storeSettings()
		} else {
			// Config file was found but another error was produced
		}
	}

	var s Settings

	err := viper.Unmarshal(&s)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
	}
	// Config file found and successfully parsed
	return s
}

func initSettings(workspacePath string) {
	//log.Println("Initialize settings ...")
	initViper(workspacePath)
	os.MkdirAll(filepath.Join(workspacePath, settingsDir), 0755)
	path := filepath.Join(workspacePath, settingsFilePath)
	err := viper.SafeWriteConfigAs(path)
	//err := viper.SafeWriteConfig()
	if err != nil {
		log.Fatal("Unable to initialize settings !", err)
	}
}

type Settings struct {
	Name string
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings name: %s", s.Name)
}

func seekSettingsFilePathRecurse(dirPath string) (string, error) {
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
	return seekSettingsFilePathRecurse(parentDirPath)
}

func seekSettingsFilePath() (string, error) {
	workDirPath := GetWorkDirPath()

	settingsPath, err := seekSettingsFilePathRecurse(workDirPath)
	return settingsPath, err
}

// --- SettingsService ---

type SettingsService struct {
	workspacePath string
	settings *Settings
}

// constructor
func newSettingsService() SettingsService {
	settingsFilePath, err := seekSettingsFilePath()
	if err != nil {
		log.Fatal(err)
	}
	workspacePath := filepath.Dir(filepath.Dir(settingsFilePath))
	initViper(workspacePath)
	settings := readSettings()

	settingsService := SettingsService{settings: &settings, workspacePath: workspacePath}
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


