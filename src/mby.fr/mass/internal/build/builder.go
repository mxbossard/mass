package build

import (
	//"bytes"

	"fmt"
	"io"
	"os/exec"
	"sync"

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
	buildCount := len(b.images)
	errors := make(chan error, buildCount)
	var wg sync.WaitGroup

	for _, image := range b.images {
		wg.Add(1)
		go func(image resources.Image) {
			defer wg.Done()
			buildDockerImage(b.binary, image, noCache, errors)
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

func stream(r io.Reader, w io.Writer, errors chan error) {
	_, err := io.Copy(w, r)
	if err != nil {
		errors <- err
	}
}

func buildDockerImage(binary string, image resources.Image, noCache bool, errors chan error) {
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
	configs, err := resources.MergedConfig(image)
	if err != nil {
		errors <- err
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errors <- err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errors <- err
	}

	streamErrors := make(chan error, 10)
	go stream(stdout, logger.Out(), streamErrors)
	go stream(stderr, logger.Err(), streamErrors)

	err = cmd.Start()
	if err != nil {
		logger.Flush()
		errors <- err
	}
	err = cmd.Wait()
	if err == nil {
		// Use select to not block if no error in channel
		select {
		case err = <-streamErrors:
		default:
		}
	}
	if err != nil {
		logger.Flush()
		err := fmt.Errorf("Error building image %s : %w", image.Name(), err)
		errors <- err
	}
	logger.Info("Build finished for image: %s .", image.Name())
}
