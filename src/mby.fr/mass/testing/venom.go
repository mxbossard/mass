package testing

import (
	"fmt"
	"sync"

	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"

	"mby.fr/utils/container"
)

var (
	venomImage = "venom:1.0.1"

	venomRunner = container.Runner{
		Image: venomImage,
		//Volumes: []string{testDirMount},
		CmdArgs: []string{"run"},
		Remove:  true,
	}

	dummyRunner = container.Runner{
		Image:   "alpine:3.16",
		CmdArgs: []string{"sh", "-c", "echo 'Dummy venom runner '; for i in $(seq 3); do sleep 1; echo .; done"},
		Remove:  true,
	}
)

func RunProjectVenomTests(d display.Displayer, p resources.Project) (err error) {
	images, err := p.Images()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(images))
	for _, i := range images {
		wg.Add(1)
		go func(i resources.Image) {
			defer wg.Done()
			err = RunImageVenomTests(d, i)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	// Wait for all tests to finish
	wg.Wait()

	// Use select to not block if no error in channel
	select {
	case err = <-errors:
	default:
	}
	if err != nil {
		return
	}

	err = RunVenomTests(d, p)
	return
}

func RunImageVenomTests(d display.Displayer, i resources.Image) (err error) {
	return RunVenomTests(d, i)
}

func RunVenomTests(d display.Displayer, res resources.Tester) (err error) {
	testDirMount := res.TestDir() + ":/venom:ro"

	runner := dummyRunner
	runner.Volumes = []string{testDirMount}

	logger := d.BufferedActionLogger("test", res.QualifiedName())
	//defer logger.Close()

	err = runner.Wait(logger.Out(), logger.Err())
	return
}

func VenomTests(d display.Displayer, res resources.Resource) (err error) {
	switch v := res.(type) {
	case resources.Project:
		RunProjectVenomTests(d, v)
	case resources.Image:
		RunImageVenomTests(d, v)
	default:
		d.Warn(fmt.Sprintf("Resource %s is not testable !", res.QualifiedName()))
		return
	}

	return
}
