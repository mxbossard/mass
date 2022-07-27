package workspace

import (
	"fmt"
	"strings"

	//"fmt"

	"mby.fr/mass/internal/build"
	"mby.fr/mass/internal/deploy"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/resources"
	"mby.fr/mass/testing"
	"mby.fr/utils/concurrent"
	"mby.fr/utils/errorz"
)

var (
	NoCacheBuild bool
	ForceBuild   bool
	ForcePull    bool
	RmVolumes    bool
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
	printErrors(errors)

	builder := func(r resources.Resource) (void interface{}, err error) {
		err = buildResource(r)
		return
	}
	_, err := concurrent.RunWaiting(builder, res...)
	if err != nil {
		d.Fatal(fmt.Sprintf("Encountered error during build phase: %s", err))
	}

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

func PullResources(args []string) {
	d := display.Service()
	d.Info("Pull starting ...")

	res, errors := ResolveExpression(args...)
	printErrors(errors)

	puller := func(r resources.Resource) (void interface{}, err error) {
		err = pullResource(r)
		return
	}
	_, err := concurrent.RunWaiting(puller, res...)
	if err != nil {
		d.Fatal(fmt.Sprintf("Encountered error during pull phase: %s", err))
	}

	d.Flush()
	d.Info("Pull finished")
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
	if ForcePull {
		PullResources(args)
	} else {
		BuildResources(args)
	}

	d := display.Service()
	d.Info("Up starting ...")

	res, errors := ResolveExpression(args...)
	printErrors(errors)

	upper := func(r resources.Resource) (void interface{}, err error) {
		err = upResource(r)
		return
	}
	_, err := concurrent.RunWaiting(upper, res...)
	if err != nil {
		d.Fatal(fmt.Sprintf("Encountered error during up phase: %s", err))
	}

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
	printErrors(errors)

	downer := func(r resources.Resource) (void interface{}, err error) {
		err = downResource(r)
		return
	}
	_, err := concurrent.RunWaiting(downer, res...)
	if err != nil {
		d.Fatal(fmt.Sprintf("Encountered error during down phase: %s", err))
	}

	d.Flush()
	d.Info("Down finished")
}

func TestResources(args []string) {
	UpResources(args)

	d := display.Service()
	d.Info("Test starting ...")

	res, errors := ResolveExpression(args...)
	printErrors(errors)

	d.Info(fmt.Sprintf("Will test resources:"))
	for _, r := range res {
		d.Info(fmt.Sprintf(" - %s", r.QualifiedName()))
	}

	tester := func(r resources.Resource) (void interface{}, err error) {
		err = testing.VenomTests(d, r)
		return
	}
	_, err := concurrent.RunWaiting(tester, res...)
	if err != nil {
		d.Fatal(fmt.Sprintf("Encountered error during test phase: %s", err))
	}

	d.Flush()
	d.Info("Test finished")
}
