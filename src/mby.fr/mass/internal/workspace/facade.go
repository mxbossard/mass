package workspace

import (
	"fmt"
	"strings"

	//"fmt"

	"sync"

	"mby.fr/mass/internal/build"
	"mby.fr/mass/internal/deploy"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/testing"
	"mby.fr/utils/errorz"
)

var (
	NoCacheBuild bool
	//BuildOnlyIfChange bool
	ForceBuild bool
	ForcePull  bool
	RmVolumes  bool
)

func printErrors(errors errorz.Aggregated) {
	if errors.GotError() {
		display := display.Service()
		display.Display(errors)
	}
}

func ResolveExpression(args ...string) ([]resources.Resource, errorz.Aggregated) {
	resourceExpr := strings.Join(args, " ")
	return resources.ResolveExpression(resourceExpr)
}

func GetResourcesConfig(args []string) {
	d := display.Service()
	d.Info("Config starting ...")

	res, errors := ResolveExpression(args...)
	printErrors(errors)
	for _, r := range res {
		config, err := resources.MergedConfig(r)
		if err != nil {
			printErrors(errors)
		}
		header := fmt.Sprintf("--- Config of %s\n", r.QualifiedName())
		footer := "---\n"
		d.Display(header, *config, footer)
	}
	d.Flush()
	d.Info("Config finished")
}

func buildResource(res resources.Resource) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	err = builder.Build(!ForceBuild, NoCacheBuild, ForcePull)
	//fmt.Println("Build finished")
	return err
}

func BuildResources(args []string) {
	d := display.Service()
	d.Info("Build starting ...")

	res, errors := ResolveExpression(args...)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			ForceBuild = true

			err := buildResource(r)
			if err != nil {
				d.Fatal(fmt.Sprintf("%v", err))
			}
		}(r)
	}
	wg.Wait()

	d.Flush()
	d.Info("Build finished")
}

func pullResource(res resources.Resource) error {
	deployer, err := deploy.New(res)
	if err != nil {
		return err
	}

	err = deployer.Pull()
	//fmt.Println("Build finished")
	return err
}

func upResource(res resources.Resource) error {
	deployer, err := deploy.New(res)
	if err != nil {
		return err
	}

	err = deployer.Deploy()
	//fmt.Println("Build finished")
	return err
}

func UpResources(args []string) {
	d := display.Service()
	d.Info("Up starting ...")

	res, errors := ResolveExpression(args...)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			if ForcePull {
				// Never build if forcePull
				err := pullResource(r)
				if err != nil {
					d.Error("Encountered error during pull phase !")
					d.Fatal(fmt.Sprintf("%v", err))
				}
			} else {
				err := buildResource(r)
				if err != nil {
					d.Error("Encountered error during build phase !")
					d.Fatal(fmt.Sprintf("%v", err))
				}
			}
			err := upResource(r)
			if err != nil {
				d.Error("Encountered error during Up !")
				d.Fatal(fmt.Sprintf("%v", err))
			}
		}(r)
	}
	wg.Wait()

	d.Flush()
	d.Info("Up finished")
}

func downResource(res resources.Resource) error {
	deployer, err := deploy.New(res)
	if err != nil {
		return err
	}

	err = deployer.Undeploy(RmVolumes)
	return err
}

func DownResources(args []string) {
	d := display.Service()
	d.Info("Down starting ...")

	res, errors := ResolveExpression(args...)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			err := downResource(r)
			if err != nil {
				d.Fatal(fmt.Sprintf("%v", err))
			}
		}(r)
	}
	wg.Wait()

	d.Flush()
	d.Info("Down finished")
}

func TestResources(args []string) {
	d := display.Service()
	d.Info("Test starting ...")

	res, errors := ResolveExpression(args...)
	printErrors(errors)

	d.Info(fmt.Sprintf("Will test resources:"))
	for _, r := range res {
		d.Info(fmt.Sprintf(" - %s", r.QualifiedName()))
		err := testing.VenomTests(d, r)
		if err != nil {
			d.Fatal(fmt.Sprintf("%v", err))
		}
	}

	d.Flush()
	d.Info("Test finished")
}
