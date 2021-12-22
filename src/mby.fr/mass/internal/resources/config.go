package resources

import (

	"mby.fr/mass/internal/settings"
	"mby.fr/mass/internal/config"
)

func MergedConfig(resourceExpr string) (configs []config.Config, err error) {
        ss, err := settings.GetSettingsService()
        if err != nil {
                return nil, err
        }

        workingEnv, err := ss.WorkingEnv()
        if err != nil {
                return nil, err
        }

        resources, err := ResolveExpression(resourceExpr)
        if err != nil {
                return nil, err
        }

        envResources, err := ResolveExpression(workingEnv, EnvKind)
        if err != nil {
                return nil, err
        }

	envConfig, err := envResources[0].Config()
        if err != nil {
                return nil, err
        }

        for _, res := range resources {
		switch r := res.(type) {
                case Env:
                        ec, err := r.Config()
			if err != nil {
				return nil, err
			}
                        configs = append(configs, ec)
                case Project:
			pc, err := r.Config()
			if err != nil {
				return nil, err
			}
                        c := config.Merge(envConfig, pc)
                        configs = append(configs, c)
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
                        configs = append(configs, c)
                }
        }

	return
}

