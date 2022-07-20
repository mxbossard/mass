package resources

import (
	"fmt"
	"os"

	"mby.fr/mass/internal/settings"
)

func ListEnvs() (envs []Env, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	envs, err = ScanEnvs(ss.EnvsDir())
	if err == os.ErrNotExist {
		err = nil
	} else if err != nil {
		fmt.Printf("Error from ListEnvs: %s\n", err)
	}
	return
}

func ListProjects() (projects []Project, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	projects, err = ScanProjects(ss.ProjectsDir())
	return
}

func ListImages() (images []Image, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return
	}
	images, err = ScanImages(ss.ProjectsDir())
	return
}

func GetProject(name string) (p Project, ok bool, err error) {
	projects, err := ListProjects()
	for _, p = range projects {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}

func GetEnv(name string) (r Env, ok bool, err error) {
	envs, err := ListEnvs()

	for _, r = range envs {
		if r.Name() == name {
			ok = true
			return
		}
	}
	return
}

func GetImage(projectName, imageName string) (r Image, ok bool, err error) {
	images, err := ListImages()
	for _, r = range images {
		if r.Name() == projectName+"/"+imageName {
			ok = true
			return
		}
	}
	return
}
