package build

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"mby.fr/mass/internal/resources"
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
	fmt.Fprintf(os.Stdout, "Building image: %s ...\n", image.Name())
	cmd := exec.Command(binary, "build", ".")
	cmd.Dir = image.Dir()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	fmt.Fprint(os.Stdout, stdout.String())
	fmt.Fprint(os.Stderr, stderr.String())
	if err != nil {
		return fmt.Errorf("Error building image %s : %w", image.Name(), err)
	}
	return
}
