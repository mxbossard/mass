package testing

import (
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"

	"mby.fr/utils/container"
)

var (
	venomImage = "venom:1.0.1"
)

func RunVenomTests(d display.Displayer, res resources.Resource) (err error) {
	/*
		testableType := reflect.TypeOf((*resources.Testable)(nil))
		ok := reflect.TypeOf(res).Implements(testableType)
		if !ok {
			return fmt.Errorf("Resource %s is not testable !", res.QualifiedName())
		}
	*/

	//testable := reflect.
	testDirMount := res.Dir() + ":/venom:ro"

	run := container.Run{
		Image:   venomImage,
		Volumes: []string{testDirMount},
		CmdArgs: []string{"run"},
		Remove:  true,
	}

	logger := d.BufferedActionLogger("test", res.QualifiedName())

	err = run.Wait(logger.Out(), logger.Err())

	return
}
