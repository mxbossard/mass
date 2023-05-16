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
	Build(onlyIfChange bool, noCache bool, forcePull bool) error
}

func New(r resources.Resourcer) (Builder, error) {
	var images []resources.Image
	switch res := r.(type) {
	case *resources.Project:
		return New(*res)
	case *resources.Image:
		return New(*res)
	case resources.Project:
		imgs, err := res.Images()
		if err != nil {
			return nil, err
		}
		for _, i := range imgs {
			images = append(images, *i)
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

func (b DockerBuilder) Build(onlyIfChange bool, noCache bool, forcePull bool) (err error) {
	buildCount := len(b.images)
	errors := make(chan error, buildCount*2)
	var wg sync.WaitGroup

	for _, image := range b.images {
		wg.Add(1)
		go func(image resources.Image) {
			defer wg.Done()
			buildDockerImage(b.binary, image, onlyIfChange, noCache, forcePull, errors)
		}(image)
	}

	// Wait for all build to finish
	wg.Wait()

	// Use select to not block if no error in channel
	select {
	case err = <-errors:
	default:
	}

	return err
}

func buildDockerImage(binary string, image resources.Image, onlyIfChange bool, noCache bool, forcePull bool, errors chan error) {
	d := display.Service()
	logger := d.BufferedActionLogger("build", image.FullTaggedName())

	err := change.Init()
	if err != nil {
		errors <- err
		return
	}

	// If forcePull => force build
	if !forcePull && onlyIfChange {
		changed, _, err := change.DoesImageChanged(image)
		if err != nil {
			errors <- err
		} else if !changed {
			logger.Info("Image: %s did not changed. Do not build it.", image.FullTaggedName())
			return
		}
	}

	logger.Info("Building image: %s ...", image.FullTaggedName())

	var buildParams []string
	buildParams = append(buildParams, "build", "-t", image.FullTaggedName())

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
		err := fmt.Errorf("Error building image %s : %w", image.FullTaggedName(), err)
		errors <- err
	}

	change.StoreImageSignature(image)

	logger.Info("Build finished for image: %s .", image.FullTaggedName())
}
