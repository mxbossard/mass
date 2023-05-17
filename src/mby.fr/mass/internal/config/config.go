package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/templates"
)

const DefaultConfigFile = templates.ConfigTemplate

type LabelsConfig map[string]string
type TagsConfig map[string]string
type EnvConfig map[string]string
type BuildArgsConfig map[string]string
type EntrypointConfig []string
type CommandArgsConfig []string

type Config struct {
	Labels      LabelsConfig      `yaml:"labels"`
	Tags        TagsConfig        `yaml:"tags"`
	Environment EnvConfig         `yaml:"environment"`
	BuildArgs   BuildArgsConfig   `yaml:"buildArgs"`
	Entrypoint  EntrypointConfig  `yaml:"entrypoint"`
	CommandArgs CommandArgsConfig `yaml:"commandArgs"`
}

// Init config in a directory path
func Init(path string, data interface{}) (err error) {
	configFilepath := filepath.Join(path, DefaultConfigFile)
	_, err = os.Stat(configFilepath)
	if os.IsNotExist(err) {
		err = nil
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
	var merged map[string]string
	if base == nil {
		merged = map[string]string{}
	} else {
		merged = base
	}
	for k, v := range replace {
		merged[k] = v
	}
	return merged
}

func mergeStringArrays(base, replace []string) []string {
	var merged []string
	if base == nil {
		return replace
	} else {
		merged = base
	}
	for k, v := range replace {
		merged[k] = v
	}
	return merged
}

// Merge several config from lowest priority to highest priority
func Merge(configs ...Config) Config {
	mergedConfig := configs[0]

	for _, c := range configs[1:] {
		mergedConfig.Labels = mergeStringMaps(mergedConfig.Labels, c.Labels)
		mergedConfig.Tags = mergeStringMaps(mergedConfig.Tags, c.Tags)
		mergedConfig.Environment = mergeStringMaps(mergedConfig.Environment, c.Environment)
		mergedConfig.BuildArgs = mergeStringMaps(mergedConfig.BuildArgs, c.BuildArgs)
		if len(c.Entrypoint) > 0 {
			mergedConfig.Entrypoint = c.Entrypoint
		}
		mergedConfig.CommandArgs = mergeStringArrays(mergedConfig.CommandArgs, c.CommandArgs)
	}

	return mergedConfig
}
