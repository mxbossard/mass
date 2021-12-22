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

func ListImages() (images []resources.Image, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	images, err = resources.ScanImages(ss.ProjectsDir())
	return
}

func GetProject(name string) (p resources.Project, ok bool, err error) {
	projects, err := ListProjects()
	for _, p = range projects {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}

func GetEnv(name string) (r resources.Env, ok bool, err error) {
	envs, err := ListEnvs()
	for _, r = range envs {
		if r.Name() == name {
			ok = true
			return
		}
	}
	return
}

func GetImage(projectName, imageName string) (r resources.Image, ok bool, err error) {
	images, err := ListImages()
	for _, r = range images {
		if r.Name() == projectName + "/" + imageName {
			ok = true
			return
		}
	}
	return
}

