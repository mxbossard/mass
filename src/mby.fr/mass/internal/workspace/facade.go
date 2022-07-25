package workspace

import (
	"fmt"
	"strings"

	//"fmt"
	"log"
	"sync"

	"mby.fr/mass/internal/build"
	"mby.fr/mass/internal/deploy"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"
	"mby.fr/utils/errorz"
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

func buildResource(res resources.Resource, onlyIfChange bool, noCache bool, force bool, forcePull bool) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	err = builder.Build(onlyIfChange, noCache, force, forcePull)
	//fmt.Println("Build finished")
	return err
}

func BuildResources(args []string, noCache bool, forcePull bool) {
	d := display.Service()
	d.Info("Build starting ...")

	res, errors := ResolveExpression(args...)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource, noCache bool) {
			defer wg.Done()
			err := buildResource(r, true, noCache, true, forcePull)
			if err != nil {
				log.Fatal(err)
			}
		}(r, noCache)
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

func UpResources(args []string, noCacheBuild bool, forceBuild bool, forcePull bool) {
	d := display.Service()
	d.Info("Uo starting ...")

	res, errors := ResolveExpression(args...)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			if !forcePull {
				// Never build if forcePull
				err := buildResource(r, false, noCacheBuild, forceBuild, false)
				if err != nil {
					d.Error("Encountered error during build phase !")
					log.Fatal(err)
				}
			} else {
				err := pullResource(r)
				if err != nil {
					d.Error("Encountered error during pull phase !")
					log.Fatal(err)
				}
			}

			err := upResource(r)
			if err != nil {
				d.Error("Encountered error during Up !")
				log.Fatal(err)
			}
		}(r)
	}
	wg.Wait()

	d.Flush()
	d.Info("Up finished")
}

func downResource(res resources.Resource, rmVolumes bool) error {
	deployer, err := deploy.New(res)
	if err != nil {
		return err
	}

	err = deployer.Undeploy(rmVolumes)
	return err
}

func DownResources(args []string, rmVolumes bool) {
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
			err := downResource(r, rmVolumes)
			if err != nil {
				log.Fatal(err)
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
	}

	d.Flush()
	d.Info("Test finished")
}
