package resources

import (
	"errors"
	"io/fs"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/errorz"
)

func MergedConfig(res Resourcer) (conf *config.Config, err error) {
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
	/*
		if !ok {
			err := fmt.Errorf("working env %s not found", workingEnv)
			return nil, err
		}
	*/
	var envConfig config.Config
	if ok {
		envConfig, err = workingEnvRes.Config()
		if err != nil {
			return nil, err
		}
	}

	switch r := res.(type) {
	//case *Env, *Project, *Image:
	//	return MergedConfig(*r)
	case *Env:
		return MergedConfig(*r)
	case *Project:
		return MergedConfig(*r)
	case *Image:
		return MergedConfig(*r)
	case Env:
		ec, err := r.Config()
		if errors.Is(err, fs.ErrNotExist) {
			// swallow config not found error
			err = nil
		} else if err != nil {
			return nil, err
		} else {
			conf = &ec
		}
	case Project:
		pc, err := r.Config()
		if errors.Is(err, fs.ErrNotExist) {
			// swallow config not found error
			err = nil
		} else if err != nil {
			return nil, err
		} else {
			c := config.Merge(envConfig, pc)
			conf = &c
		}
	case Image:
		var c config.Config
		pc, err := r.Project.Config()
		if errors.Is(err, fs.ErrNotExist) {
			// swallow config not found error
			err = nil
		} else if err != nil {
			return nil, err
		} else {
			c = config.Merge(envConfig, pc)
		}
		ic, err := r.Config()
		if errors.Is(err, fs.ErrNotExist) {
			// swallow config not found error
			err = nil
		} else if err != nil {
			return nil, err
		} else {
			c = config.Merge(c, ic)
		}
		conf = &c
	}
	return
}

func MergedConfigs(resources []Resourcer) (configs []config.Config, errors errorz.Aggregated) {
	for _, res := range resources {
		c, err := MergedConfig(res)
		if err != nil {
			errors.Add(err)
		}
		configs = append(configs, *c)
	}

	return
}
