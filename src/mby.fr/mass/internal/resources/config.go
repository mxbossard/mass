package resources

import (
	"fmt"

	"mby.fr/utils/errorz"
	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/config"
)

func MergedConfig(res Resource) (conf *config.Config, err error) {
        ss, err := settings.GetSettingsService()
        if err != nil {
                return nil, err
        }

        workingEnv, err := ss.WorkingEnv()
        if err != nil {
                return nil, err
        }

        workingEnvRes, ok, err := GetEnv(workingEnv)
        if err != nil {
                return nil, err
        }
	if !ok {
		err := fmt.Errorf("working env %s not found", workingEnv)
		return nil, err
	}

	envConfig, err := workingEnvRes.Config()
        if err != nil {
                return nil, err
        }

	switch r := res.(type) {
	case Env:
		ec, err := r.Config()
		if err != nil {
			return nil, err
		}
		conf = &ec
	case Project:
		pc, err := r.Config()
		if err != nil {
			return nil, err
		}
		c := config.Merge(envConfig, pc)
		conf = &c
	case Image:
		pc, err := r.Project.Config()
		if err != nil {
			return nil, err
		}
		ic, err := r.Config()
		if err != nil {
			return nil, err
		}
		c := config.Merge(envConfig, pc, ic)
		conf = &c
	}
	return
}

func MergedConfigs(resources []Resource) (configs []config.Config, errors errorz.Aggregated) {
        for _, res := range resources {
		c, err := MergedConfig(res)
		if err != nil {
			errors.Add(err)
		}
		configs = append(configs, *c)
        }

	return
}

