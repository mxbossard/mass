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
const settingsFile = settingsDir + "settings.yaml"

// Default settings
const defaultConfigDir = "config"
var defaultEnvs = []string{"dev", "stage", "prod"}


// --- Settings ---

type Settings struct {
	Name string
	ConfigDir string
	Environments []string
}

func initViper(workspacePath string) {
	//settingsDirPath := filepath.Join(workspacePath, settingsDir)
	settingsFilePath := filepath.Join(workspacePath, settingsFile)
	workspaceName := filepath.Base(workspacePath)

	viper.SetConfigType("yaml")
	//viper.SetConfigName("settings")
	viper.SetConfigFile(settingsFilePath)

	viper.SetDefault("Name", workspaceName)
	viper.SetDefault("ConfigDir", defaultConfigDir)
	viper.SetDefault("Environments", defaultEnvs)
}

// Store settings erasing previous settings
func storeSettings() {
	log.Println("Store settings in:", viper.ConfigFileUsed())
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
			log.Fatal("Settings file not found:", err)
		} else {
			// Config file was found but another error was produced
			log.Fatal("Unable to read settings:", err)
		}
	}

	var s Settings

	err := viper.Unmarshal(&s)
	if err != nil {
		log.Fatal("Unable to unmarshal settings:", err)
	}
	// Config file found and successfully parsed
	return s
}

func initSettings(workspacePath string) {
	//log.Println("Initialize settings ...", viper.ConfigFileUsed())
	initViper(workspacePath)
	os.MkdirAll(filepath.Join(workspacePath, settingsDir), 0755)
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Unable to initialize settings: %v", err)
	}
	fmt.Println("Initialized settings in:", viper.ConfigFileUsed())
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings name: %s", s.Name)
}

func seekSettingsFilePathRecurse(dirPath string) (string, error) {
	//log.Printf("Seek Settings in dir: %s ...\n", dirPath)
	if dirPath == "/" {
		return "", errors.New("Unable to found settings path")
	}
	settingsFilePath := filepath.Join(dirPath, settingsFile)

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
	return filepath.Join(s.workspacePath, settingsFile)
}

func (s SettingsService) ConfigDirPath() string {
	return filepath.Join(s.workspacePath, s.settings.ConfigDir)
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


