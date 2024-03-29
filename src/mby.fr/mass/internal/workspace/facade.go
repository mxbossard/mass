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
	BumpMinor    bool
	BumpMajor    bool
)

func printErrors(errors errorz.Aggregated) {
	if errors.GotError() {
		display := display.Service()
		display.Display(errors)
	}
}

func ResolveExpression(args []string, kinds ...resources.Kind) []resources.Resourcer {
	resourceExpr := strings.Join(args, " ")
	res, errors := resources.ResolveExpression(resourceExpr, kinds...)
	printErrors(errors)
	return res
}

func DisplayResourcesConfig(args []string) {
	d := display.Service()
	d.Info("Config starting ...")

	res := ResolveExpression(args, resources.AllKind)
	for _, r := range res {
		config, err := resources.MergedConfig(r)
		if err != nil {
			d.Error(fmt.Sprintf("Error merging config: %s !", err))
		}
		header := fmt.Sprintf("--- Config of %s\n", r.QualifiedName())
		footer := "---\n"
		d.Display(header, *config, footer)
	}
	d.Flush()
	d.Info("Config finished")
}

func DisplayResourcesVersion(args []string) {
	d := display.Service()
	d.Info("Version starting ...")

	res := ResolveExpression(args, resources.ImageKind)
	for _, r := range res {
		var msg string
		switch v := r.(type) {
		case *resources.Image:
			msg = fmt.Sprintf("Version of %s: %s\n", v.QualifiedName(), v.Version())
		default:
			msg = fmt.Sprintf("Resource %s is not versionable.\n", r.QualifiedName())
		}
		d.Display(msg)

	}
	d.Flush()
	d.Info("Version finished")
}

func forgeVersionBumpMessage(fromVer, toVer string) (string) {
	return fmt.Sprintf("%s => %s", fromVer, toVer)
}

func BumpResources(args []string) {
	d := display.Service()
	d.Info("Bump starting ...")

	res := ResolveExpression(args, resources.ImageKind)
	for _, r := range res {
		var i interface{} = r
		vb, ok := i.(resources.VersionBumper)
		if ok {
			toVer, fromVer, err := vb.Bump(BumpMinor, BumpMajor)
			if err != nil {
				d.Warn(fmt.Sprintf("Error bumping resource %s: %s\n", r.QualifiedName(), err))
			} else {
				msg := forgeVersionBumpMessage(fromVer, toVer)
				msg = fmt.Sprintf("Bumped resource %s: %s\n", r.QualifiedName(), msg)
				d.Display(msg)
			}
		}
	}

	d.Flush()
	d.Info("Bump finished")
}

func PromoteResources(args []string) {
	d := display.Service()
	d.Info("Promote starting ...")

	res := ResolveExpression(args, resources.ImageKind)
	for _, r := range res {
		var i interface{} = r
		vb, ok := i.(resources.VersionBumper)
		if ok {
			toVer, fromVer, err := vb.Promote()
			if err != nil {
				d.Warn(fmt.Sprintf("Error promoting resource %s: %s\n", r.QualifiedName(), err))
			} else {
				msg := forgeVersionBumpMessage(fromVer, toVer)
				msg = fmt.Sprintf("Promoted resource %s: %s\n", r.QualifiedName(), msg)
				d.Display(msg)
			}
		}
	}

	d.Flush()
	d.Info("Promote finished")
}

func ReleaseResources(args []string) {
	d := display.Service()
	d.Info("Release starting ...")

	res := ResolveExpression(args, resources.ImageKind)
	for _, r := range res {
		var i interface{} = r
		vb, ok := i.(resources.VersionBumper)
		if ok {
			toVer, fromVer, err := vb.Release()
			if err != nil {
				d.Warn(fmt.Sprintf("Error releasing resource %s: %s\n", r.QualifiedName(), err))
			} else {
				msg := forgeVersionBumpMessage(fromVer, toVer)
				msg = fmt.Sprintf("Released resource %s: %s\n", r.QualifiedName(), msg)
				d.Display(msg)
			}
		}
	}

	d.Flush()
	d.Info("Release finished")
}

func buildResource(res resources.Resourcer) error {
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

	res := ResolveExpression(args, resources.AllKind)
	builder := func(r resources.Resourcer) (void interface{}, err error) {
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

func pullResource(res resources.Resourcer) error {
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

	res := ResolveExpression(args, resources.AllKind)
	puller := func(r resources.Resourcer) (void interface{}, err error) {
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

func upResource(res resources.Resourcer) error {
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

	res := ResolveExpression(args, resources.AllKind)
	upper := func(r resources.Resourcer) (void interface{}, err error) {
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

func downResource(res resources.Resourcer) error {
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

	res := ResolveExpression(args, resources.AllKind)
	downer := func(r resources.Resourcer) (void interface{}, err error) {
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

	res := ResolveExpression(args, resources.AllKind)
	d.Info(fmt.Sprintf("Will test resources:"))
	for _, r := range res {
		d.Info(fmt.Sprintf(" - %s", r.QualifiedName()))
	}

	tester := func(r resources.Resourcer) (void interface{}, err error) {
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
