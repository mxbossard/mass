package build

import (
	//"bytes"

	"fmt"
	"os/exec"
	"sync"

	"mby.fr/mass/internal/change"
	"mby.fr/mass/internal/command"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"
)

var NotBuildableResource error = fmt.Errorf("Not buildable resource")

type Builder interface {
	Build(noCache bool, force bool, forcePull bool) error
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
		return nil, fmt.Errorf("%w: %s", NotBuildableResource, r.QualifiedName())
	}
	return DockerBuilder{"docker", images}, nil
}

type DockerBuilder struct {
	binary string
	images []resources.Image
}

func (b DockerBuilder) Build(noCache bool, force bool, forcePull bool) (err error) {
	buildCount := len(b.images)
	errors := make(chan error, buildCount)
	var wg sync.WaitGroup

	for _, image := range b.images {
		doBuild := force
		if !doBuild {
			err = change.Init()

			if err != nil {
				return
			}
			doBuild, _, err = change.DoesImageChanged(image)
		}

		if doBuild {
			wg.Add(1)
			go func(image resources.Image) {
				defer wg.Done()
				buildDockerImage(b.binary, image, noCache, forcePull, errors)
			}(image)
		}
	}

	// Wait for all build to finish
	wg.Wait()

	// Use select to not block if no error in channel
	select {
	case err = <-errors:
	default:
	}

	for _, image := range b.images {
		change.StoreImageSignature(image)
	}

	return err
}

func buildDockerImage(binary string, image resources.Image, noCache bool, forcePull bool, errors chan error) {
	d := display.Service()
	logger := d.BufferedActionLogger("build", image.Name())
	logger.Info("Building image: %s ...", image.Name())

	var buildParams []string
	buildParams = append(buildParams, "build", "-t", image.FullName())

	// Add --no-cache option
	if noCache {
		buildParams = append(buildParams, "--no-cache")
	}

	if forcePull {
		buildParams = append(buildParams, "--pull")
	}

	// Forge build-args
	config, err := resources.MergedConfig(image)
	if err != nil {
		errors <- err
	}

	for argKey, argValue := range config.BuildArgs {
		var buildArg string = "--build-arg=" + argKey + "=" + argValue
		buildParams = append(buildParams, buildArg)
	}

	// Add dot folder as last param
	buildParams = append(buildParams, ".")

	logger.Debug("build params: %s", buildParams)
	cmd := exec.Command(binary, buildParams...)
	cmd.Dir = image.Dir()

	err = command.RunLogging(cmd, logger)
	if err != nil {
		logger.Flush()
		err := fmt.Errorf("Error building image %s : %w", image.Name(), err)
		errors <- err
	}

	logger.Info("Build finished for image: %s .", image.Name())
}
