package workspace

import (
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/settings"
)

func ListEnvs() (envs []resources.Env, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	envs, err = resources.ScanEnvs(ss.EnvsDir())
	return
}

func ListProjects() (projects []resources.Project, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	projects, err = resources.ScanProjects(ss.ProjectsDir())
	return
}

