package build

import (
	//"bytes"
	"fmt"
	"os/exec"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/display"
)

var NotBuildableResource error = fmt.Errorf("Not buildable resource")

type Builder interface {
	Build() error
}

func New(r resources.Resource) (Builder, error) {
	var images []resources.Image
	var err error
	switch res := r.(type) {
	case resources.Project:
		images, err = res.Images()
		if err != nil {
			return nil, err
		}
	case resources.Image:
		images = append(images, res)
	default:
		return nil, fmt.Errorf("%w: %s", NotBuildableResource, r.AbsoluteName())
	}
	return DockerBuilder{"docker", images}, nil
}

type DockerBuilder struct {
	binary string
	images []resources.Image
}

func (b DockerBuilder) Build() (err error) {
	for _, image := range b.images {
		err = buildDockerImage(b.binary, image)
		if err != nil {
			return
		}
	}
	return
}

func buildDockerImage(binary string, image resources.Image) (err error) {
	d := display.Service()
	logger := d.BufferedActionLogger("build", image.Name())
	logger.Info("Building image: %s ...", image.Name())
	cmd := exec.Command(binary, "build", ".")
	cmd.Dir = image.Dir()
	cmd.Stdout = logger.Out()
	cmd.Stderr = logger.Err()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error building image %s : %w", image.Name(), err)
	}
	logger.Info("Build finished for image: %s .", image.Name())
	return
}
