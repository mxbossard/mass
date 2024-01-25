package resources

//"path/filepath"

type Project struct {
	directoryBase   `yaml:"base,inline"`
	configurableDir `yaml:"-"` // Ignore this field for yaml marshalling
	testable        `yaml:"testable,inline"`

	DeploymentDeps []string `yaml:"dependencies.deployments"`

	images      []*Image      `yaml:"-"` // Ignore this field for yaml marshalling
	deployments []*Deployment `yaml:"-"` // Ignore this field for yaml marshalling
}

func (p Project) init() (err error) {
	err = p.directoryBase.init()
	if err != nil {
		return
	}
	err = p.testable.init()
	if err != nil {
		return
	}

	return
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

func buildProject(projectDir string) (p Project, err error) {
	b, err := buildDirectoryBase(ProjectKind, projectDir, resourceName(projectDir))
	if err != nil {
		return
	}
	p = Project{directoryBase: b}

	p.configurableDir = buildConfigurableDir(b)
	t, err := buildTestable(p)
	if err != nil {
		return
	}
	p.testable = t

	return
}
