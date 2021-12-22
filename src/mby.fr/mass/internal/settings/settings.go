package settings

import (
	"fmt"
	"sync"
	"path/filepath"
	"log"
	"os"
	"errors"

	"github.com/spf13/viper"

	"mby.fr/mass/internal/templates"
)

const defaultSettingsDir = ".mass"
var settingsFile = filepath.Join(defaultSettingsDir, "settings.yaml")

var PathNotFound = fmt.Errorf("Unable to found settings path")
var NotExistingEnv = fmt.Errorf("Env don't exists")

// Default settings
const defaultEnvsDir = "envs"
const defaultProjectsDir = "."
const defaultCacheDir = ".cache"
const defaultTemplatesDir = ".templates"
const defaultEnvToUse = "dev"
var defaultEnvs = []string{"dev", "stage", "prod"}

var SelectedEnvironment string = ""

// --- Settings ---

type Settings struct {
	Name string `yaml:"workspace"`
	EnvsDir string `yaml:"envsDirectory"`
	ProjectsDir string `yaml:"projectsDirectory"`
	CacheDir string `yaml:"cacheDirectory"`
	TemplatesDir string `yaml:"templatesDirectory"`
	Environments []string `yaml:"environments"`
	DefaultEnvironment string `yaml:"defaultEnvironment"`
}

func Default() Settings {
	return Settings{
		EnvsDir: defaultEnvsDir,
		ProjectsDir: defaultProjectsDir,
		CacheDir: defaultCacheDir,
		TemplatesDir: defaultTemplatesDir,
		Environments: defaultEnvs,
		DefaultEnvironment: defaultEnvToUse,
	}
}

func initViper(workspacePath string) {
	//settingsDirPath := filepath.Join(workspacePath, defaultSettingsDir)
	settingsFilePath := filepath.Join(workspacePath, settingsFile)
	workspaceName := filepath.Base(workspacePath)

	viper.SetConfigType("yaml")
	// Do not use SetConfigName, use SetConfigFile instead
	//viper.SetConfigName("settings")
	viper.SetConfigFile(settingsFilePath)

	viper.SetDefault("Name", workspaceName)
	viper.SetDefault("EnvsDir", defaultEnvsDir)
	viper.SetDefault("ProjectsDir", defaultProjectsDir)
	viper.SetDefault("CacheDir", defaultCacheDir)
	viper.SetDefault("TemplatesDir", defaultTemplatesDir)
	viper.SetDefault("Environments", defaultEnvs)
	viper.SetDefault("DefaultEnvironment", defaultEnvToUse)
}

// Store settings erasing previous settings
func storeSettings() (err error) {
	log.Println("Store settings in:", viper.ConfigFileUsed())
	err = viper.WriteConfig()
	if err != nil {
		fmt.Errorf("Unable to store settings: %w !", err)
	}
	return
}

func readSettings() (s *Settings, err error) {
	// Find and read the config file
	if err = viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			err = fmt.Errorf("Settings file not found: %w !", err)
			return
		} else {
			// Config file was found but another error was produced
			err = fmt.Errorf("Unable to read settings: %w !", err)
			return
		}
	}

	err = viper.Unmarshal(&s)
	if err != nil {
		err = fmt.Errorf("Unable to unmarshal settings: %w !", err)
		return
	}
	// Config file found and successfully parsed
	return
}

func Init(workspacePath string) (err error) {
	// FIXME: use built Settings to init the settings FS ?
	//log.Println("Initialize settings ...", viper.ConfigFileUsed())
	initViper(workspacePath)
	newSettingsDir := filepath.Join(workspacePath, defaultSettingsDir)
	err = os.MkdirAll(newSettingsDir, 0755)
	if err != nil {
		return
	}
	err = viper.WriteConfig()
	if err != nil {
		err = fmt.Errorf("Unable to initialize settings: %w", err)
		return
	}

	newTemplatesDir := filepath.Join(newSettingsDir, defaultTemplatesDir)
	err = os.MkdirAll(newTemplatesDir, 0755)
	if err != nil {
		return
	}
	err = templates.Init(newTemplatesDir)
	if err != nil {
		err = fmt.Errorf("Unable to initialize settings: %w", err)
		return
	}

	fmt.Println("Initialized settings in:", viper.ConfigFileUsed())
	return
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings worspace name: %s", s.Name)
}

func seekSettingsFilePathRecurse(dirPath string) (string, error) {
	//log.Printf("Seek Settings in dir: %s ...\n", dirPath)
	if dirPath == "/" {
		return "", nil
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

func SeekSettingsFilePath(path string) (settingsPath string, err error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	settingsPath, err = seekSettingsFilePathRecurse(absolutePath)
	return
}

// --- SettingsService ---

type SettingsService struct {
	workspacePath string
	settings *Settings
}

// constructor
func newSettingsService() (service *SettingsService, err error) {
	workDirPath, err := os.Getwd()
	if err != nil {
		return
	}
	settingsFilePath, err := SeekSettingsFilePath(workDirPath)
	if err != nil {
		return
	}
	if settingsFilePath == "" {
		err = PathNotFound
		return
	}
	workspacePath := filepath.Dir(filepath.Dir(settingsFilePath))
	initViper(workspacePath)
	settings, err := readSettings()
	if err != nil {
		return
	}

	service = &SettingsService{settings: settings, workspacePath: workspacePath}
	return
}

// workspacePath getter
func (s SettingsService) WorkspaceDir() string {
	return s.workspacePath
}

// settings getter
func (s SettingsService) Settings() *Settings {
	return s.settings
}

func (s SettingsService) SettingsDir() string {
	return filepath.Join(s.workspacePath, defaultSettingsDir)
}

func (s SettingsService) TemplatesDir() string {
	return filepath.Join(s.SettingsDir(), s.settings.TemplatesDir)
}

func (s SettingsService) TemplatesRenderer() templates.Renderer {
	renderer := templates.New(s.TemplatesDir())
	return renderer
}

func (s SettingsService) SettingsFile() string {
	return filepath.Join(s.workspacePath, settingsFile)
}

func (s SettingsService) EnvsDir() string {
	return filepath.Join(s.workspacePath, s.settings.EnvsDir)
}

func (s SettingsService) ProjectsDir() string {
	return filepath.Join(s.workspacePath, s.settings.ProjectsDir)
}

func (s SettingsService) CacheDir() string {
	return filepath.Join(s.workspacePath, s.settings.CacheDir)
}

func (s SettingsService) WorkingEnv() (string, error) {
	envToUse := s.settings.DefaultEnvironment
	if SelectedEnvironment != "" {
		// User specified an environment
		envToUse = SelectedEnvironment
	}

	envExists := false
	for _, e := range s.settings.Environments {
		//fmt.Printf("testing env %s == %s\n", envToUse, e)
		if e == envToUse {
			envExists = true
			break
		}
	}

	if !envExists {
		return "", NotExistingEnv
	}

	return envToUse, nil
}

// singleton
var lock = &sync.Mutex{}

var settingsService *SettingsService

func GetSettingsService() (service *SettingsService, err error) {
	lock.Lock()
	defer lock.Unlock()
	// FIXME: disable singleton because unitest are failing.
	if settingsService == nil || true {
		service, err = newSettingsService()
	}
	return
}

