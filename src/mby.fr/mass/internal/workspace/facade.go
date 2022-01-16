package workspace

import (
	"strings"
	//"fmt"
	"log"
	"sync"

	"mby.fr/utils/errorz"
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/build"
	"mby.fr/mass/internal/display"
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

func buildResource(res resources.Resource) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	err = builder.Build()
	//fmt.Println("Build finished")
	return err
}

func BuildResources(args []string) {
	res, errors := ResolveExpression(args)
	//fmt.Println(res)
	printErrors(errors)
	var wg sync.WaitGroup
	for _, r := range res {
		wg.Add(1)
		go func(r resources.Resource) {
			defer wg.Done()
			err := buildResource(r)
			if err != nil {
				log.Fatal(err)
			}
		}(r)
	}
	wg.Wait()

	d := display.Service()
	d.Flush()
	d.Info("Build finished")
}

func GetResourcesConfig(args []string) {
	res, errors := ResolveExpression(args)
	printErrors(errors)
	configs, errors := resources.MergedConfigs(res)
	display := display.Service()
	display.Display(configs, errors)
}
