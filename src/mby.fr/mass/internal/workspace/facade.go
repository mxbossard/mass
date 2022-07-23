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

func ResolveExpression(args []string) ([]resources.Resource, errorz.Aggregated) {
	resourceExpr := strings.Join(args, " ")
	return resources.ResolveExpression(resourceExpr)
}

func GetResourcesConfig(args []string) {
	res, errors := ResolveExpression(args)
	printErrors(errors)
	d := display.Service()
	for _, r := range res {
		config, err := resources.MergedConfig(r)
		if err != nil {
			printErrors(errors)
		}
		header := fmt.Sprintf("--- Config of %s %s\n", r.Kind(), r.QualifiedName())
		footer := "---\n"
		d.Display(header, *config, footer)
	}
	d.Info("Config finished")
	d.Flush()
}

func buildResource(res resources.Resource, noCache bool, force bool) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	err = builder.Build(noCache, force)
	//fmt.Println("Build finished")
	return err
}

func BuildResources(args []string, noCache bool) {
	d := display.Service()
	d.Info("Starting build ...")

	res, errors := ResolveExpression(args)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource, noCache bool) {
			defer wg.Done()
			err := buildResource(r, noCache, true)
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
	d.Info("Starting up ...")

	res, errors := ResolveExpression(args)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			if !forcePull {
				// Never build if forcePull
				err := buildResource(r, noCacheBuild, forceBuild)
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
	d.Info("Starting down ...")

	res, errors := ResolveExpression(args)
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
