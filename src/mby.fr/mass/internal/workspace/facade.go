package workspace

import (
	"strings"
	"fmt"
	"log"

	"mby.fr/utils/errorz"
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/internal/build"
	"mby.fr/mass/internal/display"
)

func ResolveExpression(args []string) ([]resources.Resource, errorz.Aggregated) {
	resourceExpr := strings.Join(args, " ")
	return resources.ResolveExpression(resourceExpr)
}

func buildResource(res resources.Resource) error {
	builder, err := build.New(res)
	if err != nil {
		return err
	}

	fmt.Println("Build finished")
	err = builder.Build()
	fmt.Println("Build finished")
	return err
}

func printErrors(errors errorz.Aggregated) {
	if errors.GotError() {
		display := display.Service()
		display.Display(errors)
	}
}

func BuildResources(args []string) {
	res, errors := ResolveExpression(args)
	//fmt.Println(res)
	printErrors(errors)
	for _, r := range res {
		err := buildResource(r)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Build finished")
}

func GetResourcesConfig(args []string) {
	res, errors := ResolveExpression(args)
	printErrors(errors)
	configs, errors := resources.MergedConfigs(res)
	display := display.Service()
	display.Display(configs, errors)
}
