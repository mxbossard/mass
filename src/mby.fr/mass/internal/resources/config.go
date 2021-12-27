package resources

import (
	"fmt"

	"mby.fr/utils/errorz"
	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/config"
)

func MergedConfig(resourceExpr string) (configs []config.Config, errors errorz.Aggregated) {
        ss, err := settings.GetSettingsService()
        if err != nil {
		errors.Add(err)
                return nil, errors
        }

        workingEnv, err := ss.WorkingEnv()
        if err != nil {
		errors.Add(err)
                return nil, errors
        }

        resources, errs := ResolveExpression(resourceExpr)
        if errs.GotError() {
                return nil, errs
        }

        workingEnvRes, ok, err := GetEnv(workingEnv)
        if err != nil {
		errors.Add(err)
                return nil, errors
        }
	if !ok {
		err := fmt.Errorf("working env %s not found", workingEnv)
		errors.Add(err)
		return nil, errors
	}

	envConfig, err := workingEnvRes.Config()
        if err != nil {
		errors.Add(err)
                return nil, errors
        }

        for _, res := range resources {
		switch r := res.(type) {
                case Env:
                        ec, err := r.Config()
			if err != nil {
				errors.Add(err)
				return nil, errors
			}
                        configs = append(configs, ec)
                case Project:
			pc, err := r.Config()
			if err != nil {
				errors.Add(err)
				return nil, errors
			}
                        c := config.Merge(envConfig, pc)
                        configs = append(configs, c)
                case Image:
			pc, err := r.Project.Config()
			if err != nil {
				errors.Add(err)
				return nil, errors
			}
			ic, err := r.Config()
			if err != nil {
				errors.Add(err)
				return nil, errors
			}
                        c := config.Merge(envConfig, pc, ic)
                        configs = append(configs, c)
                }
        }

	return
}

