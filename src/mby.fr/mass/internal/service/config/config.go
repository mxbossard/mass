package config

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

// --- MassConfig ---

type MassConfig struct {
	Name string
	WorkspacePath string
}

func (c MassConfig) String() string {
	return fmt.Sprintf("MassConfig name: %s ; workspacePath: %s.", c.Name, c.WorkspacePath)
}

func (c MassConfig) ConfigDirPath() string {
	return filepath.Join(c.WorkspacePath, "/.mass/")
}

func (c MassConfig) ConfigFilePath() string {
	return filepath.Join(c.ConfigDirPath(), "config.yaml")
}

func (c MassConfig) store() {
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

func InitMassConfig() {
	workspacePath := getWorkDirPath()
	name := filepath.Base(workspacePath)
	config := MassConfig{name, workspacePath}

	_, err := os.Stat(config.ConfigFilePath())
	if err == nil {
		// config file already exists
		log.Fatal("mass config already exists !")
	}

	config.store()
}

func seekMassConfigPathRecurse(dirPath string) (string, error) {
	//log.Printf("Seek MassConfig in dir: %s ...\n", dirPath)
	if dirPath == "/" {
		return "", errors.New("Unable to found config path")
	}
	configFilePath := filepath.Join(dirPath, "/.mass/config.yaml")

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
	return seekMassConfigPathRecurse(parentDirPath)
}

func seekMassConfigPath() string {
	workDirPath := getWorkDirPath()

	configPath, err := seekMassConfigPathRecurse(workDirPath)

	if configPath != "" && err == nil {
		return configPath
	}
	return ""
}

func loadMassConfig() (*MassConfig, bool) {
	configFilePath := seekMassConfigPath()
	if configFilePath == "" {
		// Mass config does not exists yet
		return nil, false
	}

	yfile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := MassConfig{}

	err = yaml.Unmarshal(yfile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config, true
}

// --- MassConfigService ---

type MassConfigService struct {
	config *MassConfig
}

// constructor
func newMassConfigService() MassConfigService {
	massConfig, ok := loadMassConfig()
	if !ok {
		log.Fatal("Unable to load mass config !")
	}
	configService := MassConfigService{massConfig}
	return configService
}

// config getter
func (s MassConfigService) Config() *MassConfig {
	return s.config
}

// singleton
var lock = &sync.Mutex{}

var maasConfigService *MassConfigService

func GetMassConfigService() *MassConfigService {
	if maasConfigService == nil {
		lock.Lock()
		defer lock.Unlock()
		if maasConfigService == nil {
			service := newMassConfigService()
			maasConfigService = &service
		}
	}
	return maasConfigService
}


