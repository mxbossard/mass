package build

import (
	//"bytes"
	"fmt"
	"os/exec"

	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"
)

var NotBuildableResource error = fmt.Errorf("Not buildable resource")

type Builder interface {
	Build(noCache bool) error
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

func (b DockerBuilder) Build(noCache bool) (err error) {
	for _, image := range b.images {
		err = buildDockerImage(b.binary, image, noCache)
		if err != nil {
			return
		}
	}
	return
}

func buildDockerImage(binary string, image resources.Image, noCache bool) (err error) {
	d := display.Service()
	logger := d.BufferedActionLogger("build", image.Name())
	logger.Info("Building image: %s ...", image.Name())

	var buildParams []string
	buildParams = append(buildParams, "build", "-t", image.FullName())

	// Add --no-cache option
	if noCache {
		buildParams = append(buildParams, "--no-cache")
	}

	// Forge build-args
	configs, errors := resources.MergedConfig(image)
	if errors != nil {
		return errors
	}

	for argKey, argValue := range configs.BuildArgs {
		var buildArg string = "--build-arg=" + argKey + "=" + argValue
		buildParams = append(buildParams, buildArg)
	}

	// Add dot folder as last param
	buildParams = append(buildParams, ".")

	logger.Debug("build params: %s", buildParams)
	cmd := exec.Command(binary, buildParams...)
	cmd.Dir = image.Dir()
	cmd.Stdout = logger.Out()
	cmd.Stderr = logger.Err()
	err = cmd.Run()
	if err != nil {
		logger.Flush()
		return fmt.Errorf("Error building image %s : %w", image.Name(), err)
	}
	logger.Info("Build finished for image: %s .", image.Name())
	return
}
