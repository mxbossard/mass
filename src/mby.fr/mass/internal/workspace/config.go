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

const configDir = ".mass/"
const configFilePath = "config.yaml"

// --- Config ---

type Config struct {
	Name string
	WorkspacePath string
}

func (c Config) String() string {
	return fmt.Sprintf("Config name: %s ; workspacePath: %s.", c.Name, c.WorkspacePath)
}

func (c Config) ConfigDirPath() string {
	return filepath.Join(c.WorkspacePath, configDir)
}

func (c Config) ConfigFilePath() string {
	return filepath.Join(c.ConfigDirPath(), "config.yaml")
}

func (c Config) store() {
	configFilePath := c.ConfigFilePath()
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(c.ConfigDirPath(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(configFilePath, data, 0755)
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

func InitConfig() {
	workspacePath := getWorkDirPath()
	name := filepath.Base(workspacePath)
	config := Config{name, workspacePath}

	_, err := os.Stat(config.ConfigFilePath())
	if err == nil {
		// config file already exists
		log.Fatal("Workspace config already exists !")
	}

	config.store()
}

func seekConfigPathRecurse(dirPath string) (string, error) {
	//log.Printf("Seek Config in dir: %s ...\n", dirPath)
	if dirPath == "/" {
		return "", errors.New("Unable to found config path")
	}
	configFilePath := filepath.Join(dirPath, "/" + configFilePath)

	_, err := os.Stat(configFilePath)
	if err == nil {
		// config file exists
		return configFilePath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// config file does *not* exist

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	}

	parentDirPath := filepath.Dir(dirPath)
	return seekConfigPathRecurse(parentDirPath)
}

func seekConfigPath() string {
	workDirPath := getWorkDirPath()

	configPath, err := seekConfigPathRecurse(workDirPath)

	if configPath != "" && err == nil {
		return configPath
	}
	return ""
}

func loadConfig() (*Config, bool) {
	configFilePath := seekConfigPath()
	if configFilePath == "" {
		// Workspace config does not exists yet
		return nil, false
	}

	yfile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config, true
}

// --- ConfigService ---

type ConfigService struct {
	config *Config
}

// constructor
func newConfigService() ConfigService {
	config, ok := loadConfig()
	if !ok {
		log.Fatal("Unable to load workspace config !")
	}
	configService := ConfigService{config}
	return configService
}

// config getter
func (s ConfigService) Config() *Config {
	return s.config
}

// singleton
var lock = &sync.Mutex{}

var configService *ConfigService

func GetConfigService() *ConfigService {
	if configService == nil {
		lock.Lock()
		defer lock.Unlock()
		if configService == nil {
			service := newConfigService()
			configService = &service
		}
	}
	return configService
}


