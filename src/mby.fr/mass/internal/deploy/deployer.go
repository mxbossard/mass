package deploy

import (
	//"bytes"
	"fmt"
	"os/exec"

	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/display"
)

var NotDeployableResource error = fmt.Errorf("Not deployable resource")

type Deployer interface {
	Deploy() error
}

func New(r resources.Resource) (Deployer, error) {
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
		return nil, fmt.Errorf("%w: %s", NotDeployableResource, r.AbsoluteName())
	}
	return DockerDeployer{"docker", images}, nil
}

type DockerDeployer struct {
	binary string
	images []resources.Image
}

func (d DockerDeployer) Deploy() (err error) {
	for _, image := range d.images {
		err = runDockerImage(d.binary, image)
		if err != nil {
			return
		}
	}
	return
}

func runDockerImage(binary string, image resources.Image) (err error) {
	d := display.Service()
	logger := d.BufferedActionLogger("run", image.Name())
	logger.Info("Running image: %s ...", image.FullName())
	cmd := exec.Command(binary, "run", image.FullName(), "finalArg")
	cmd.Dir = image.Dir()
	cmd.Stdout = logger.Out()
	cmd.Stderr = logger.Err()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running image %s : %w", image.FullName(), err)
	}
	logger.Info("Run finished for image: %s .", image.Name())
	return
}
