package testing

import (
	"fmt"

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
		CmdArgs: []string{"echo", "Dummy venom runner"},
		Remove:  true,
	}
)

func RunVenomTests(d display.Displayer, res resources.Resource) (err error) {
	/*
		testableType := reflect.TypeOf((*resources.Testable)(nil))
		ok := reflect.TypeOf(res).Implements(testableType)
		if !ok {
			return fmt.Errorf("Resource %s is not testable !", res.QualifiedName())
		}
	*/

	testable, ok := res.(resources.Testable)
	if !ok {
		//return fmt.Errorf("Resource %s is not testable !", res.QualifiedName())
		d.Warn(fmt.Sprintf("Resource %s is not testable !", res.QualifiedName()))
		return
	}

	//testable := reflect.
	testDirMount := testable.TestDir() + ":/venom:ro"

	runner := dummyRunner
	runner.Volumes = []string{testDirMount}

	logger := d.BufferedActionLogger("test", res.QualifiedName())

	err = runner.Wait(logger.Out(), logger.Err())

	return
}
