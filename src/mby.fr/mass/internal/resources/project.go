package resources

import (
	"mby.fr/mass/internal/settings"
	"mby.fr/utils/file"
)

type Project struct {
	base     `yaml:"base,inline"`
	testable `yaml:"testable,inline"`

	images     []*Image
	DeployFile string `yaml:"deployFile"`
}

func (p Project) init() (err error) {
	err = p.base.init()
	if err != nil {
		return
	}
	err = p.testable.init()
	if err != nil {
		return
	}

	// Init Deploy file
	deployfileContent := ""
	//buildfileContent := "FROM alpine\n"
	_, err = file.SoftInitFile(p.DeployFile, deployfileContent)

	return
}

func (p Project) AbsoluteName() (name string, err error) {
	ss, err := settings.GetSettingsService()
	if err != nil {
		return "", err
	}
	name = ss.Settings().Name + "-" + p.Name()
	return
}

func (p Project) AbsDeployFile() string {
	return absResourcePath(p.Dir(), p.DeployFile)
}

func (p *Project) Images() ([]*Image, error) {
	var err error = nil
	if len(p.images) == 0 {
		images, err := Scan[Image](p.Dir())
		var imagesPtrs []*Image
		for i := 0; i < len(images); i++ {
			imagesPtrs = append(imagesPtrs, &images[i])
		}
		if err != nil {
			return []*Image{}, err
		}
		p.images = imagesPtrs
	}
	return p.images, err
}

func buildProject(path string) (p Project, err error) {
	deployfile := DefaultDeployFile

	b, err := buildBase(ProjectKind, path, DefaultResourceFile)
	if err != nil {
		return
	}
	p = Project{base: b, DeployFile: deployfile}

	t, err := buildTestable(b, path)
	if err != nil {
		return
	}
	p.testable = t

	return
}
