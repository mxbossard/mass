package config

import(
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"mby.fr/mass/internal/templates"
	"mby.fr/mass/internal/settings"
)

const DefaultConfigFile = templates.ConfigTemplate

type EnvConfig map[string]string

type Config struct {
	Environment EnvConfig
}

// Init config in a directory path
func Init(path string, data interface{}) (err error) {
	configFilepath := filepath.Join(path, DefaultConfigFile)
	_, err = os.Stat(configFilepath)
	if os.IsNotExist(err) {
		ss, err := settings.GetSettingsService()
		if err != nil {
			return err
		}
		renderer := ss.TemplatesRenderer()
		err = renderer.RenderToFile(templates.ConfigTemplate, configFilepath, data)
		if err != nil {
			return err
		}
	}
	return
}

// Read config from dir or file
func Read(path string) (c Config, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}
	if info.IsDir() {
		path = filepath.Join(path, DefaultConfigFile)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return
	}

	return
}

func mergeStringMaps(base, replace map[string]string) map[string]string {
	merged := base
	for k, v := range replace {
		merged[k] = v
	}
	return merged
}

// Merge several config from lowest priority to highest priority
func Merge(configs ...Config) (Config) {
	mergedConfig := configs[0]

	for _, c := range configs[1:] {
		mergedEnv := mergeStringMaps(mergedConfig.Environment, c.Environment)
		mergedConfig.Environment = mergedEnv
	}

	return mergedConfig
}

