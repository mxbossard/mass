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

func buildResource(res resources.Resource, noCache bool) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	err = builder.Build(noCache)
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
			err := buildResource(r, noCache)
			if err != nil {
				log.Fatal(err)
			}
		}(r, noCache)
	}
	wg.Wait()

	d.Flush()
	d.Info("Build finished")
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
	d.Info("Starting up ...")

	res, errors := ResolveExpression(args)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			err := upResource(r)
			if err != nil {
				log.Fatal(err)
			}
		}(r)
	}
	wg.Wait()

	d.Flush()
	d.Info("Up finished")
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
		header := fmt.Sprintf("--- Config of %s %s\n", r.Kind(), r.AbsoluteName())
		footer := "---\n"
		d.Display(header, *config, footer)
	}
	d.Info("Config finished")
	d.Flush()
}
