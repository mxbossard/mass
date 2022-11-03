package deploy

import (
	//"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"mby.fr/mass/internal/command"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/logger"
	"mby.fr/mass/internal/resources"

	"mby.fr/utils/errorz"
)

const composeImage = "docker/compose:alpine-1.29.2"
const defaultComposeDownTimeoutInSec = "60"

var NotDeployableResource error = fmt.Errorf("Not deployable resource")

type Deployer interface {
	Pull() error
	Deploy() error
	Undeploy(rmVolumes bool) error
}

func New(r resources.Resourcer) (Deployer, error) {
	switch res := r.(type) {
	case *resources.Project:
		return New(*res)
	case *resources.Image:
		return New(*res)
	case resources.Project:
		return DockerComposeProjectsDeployer{"docker", []string{}, []resources.Project{res}}, nil
	case resources.Image:
		return DockerImagesDeployer{"docker", []string{}, []resources.Image{res}}, nil
	default:
		return nil, fmt.Errorf("%w: %s", NotDeployableResource, r.QualifiedName())
	}
}

type DockerImagesDeployer struct {
	binary string
	args   []string
	images []resources.Image
}

func (d DockerImagesDeployer) Pull() (err error) {
	for _, image := range d.images {
		err = pullImage(d.binary, image)
		if err != nil {
			return
		}
	}
	return
}

func (d DockerImagesDeployer) Deploy() (err error) {
	for _, image := range d.images {
		err = runImage(d.binary, image)
		if err != nil {
			return
		}
	}
	return
}

func (d DockerImagesDeployer) Undeploy(rmVolumes bool) (err error) {
	// FIXME: remove persistent volumes
	return undeployContainers(d.binary, d.images)
}

func absContainerName(image resources.Image) (name string, err error) {
	imageAbsName, err := image.AbsoluteName()
	if err != nil {
		return
	}
	re := regexp.MustCompile("[/ ]")
	name = re.ReplaceAllString(imageAbsName, "_")
	return
}

func pullImage(binary string, image resources.Image) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("pull", image.FullName())

	var pullParams []string
	pullParams = append(pullParams, "pull", image.FullName())

	log.Debug("pull params: %s", pullParams)
	cmd := exec.Command(binary, pullParams...)
	//cmd.Dir = image.Dir()

	err = command.RunLogging(cmd, log)
	if err != nil {
		flushErr := d.Flush()
		err = fmt.Errorf("Error pulling image %s : %w", image.FullName(), err)
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	log.Info("Pull finished for image: %s .", image)
	return
}

func runDockerImage(log logger.ActionLogger, binary string, runArgs []string, name string, image string, cmdArgs ...string) (err error) {
	log.Info("Running image: %s as: %s ...", image, name)

	var runParams []string
	runParams = append(runParams, "run")

	runParams = append(runParams, runArgs...)

	if name != "" {
		runParams = append(runParams, "--name", name)
	}

	runParams = append(runParams, image)

	runParams = append(runParams, cmdArgs...)

	log.Debug("run params: %s", runParams)
	cmd := exec.Command(binary, runParams...)
	//cmd.Dir = image.Dir()

	err = command.RunLogging(cmd, log)
	if err != nil {
		// flushErr := log.Flush()
		// agg := errorz.NewAggregated(err, flushErr)
		return fmt.Errorf("Error running image %s : %w", image, err)
	}
	log.Info("Run finished for image: %s .", image)
	return
}

func runImage(binary string, image resources.Image) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("run", image.FullName())

	var runArgs []string
	config, errors := resources.MergedConfig(image)
	if errors != nil {
		return errors
	}
	// Add envVars
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runArgs = append(runArgs, envArg)
	}

	//runArgs = append(runArgs, "badArg")

	var cmdArgs []string
	// Add runParams
	for _, argValue := range config.RunArgs {
		cmdArgs = append(cmdArgs, argValue)
	}

	ctName, err := absContainerName(image)
	if err != nil {
		return
	}

	err = runDockerImage(log, binary, runArgs, ctName, image.FullName(), cmdArgs...)
	if err != nil {
		flushErr := d.Flush()
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	return
}

func rmDockerContainers(log logger.ActionLogger, binary string, names ...string) (err error) {
	var rmParams []string
	rmParams = append(rmParams, "rm", "-f")
	rmParams = append(rmParams, names...)

	cmd := exec.Command(binary, rmParams...)
	//cmd.Dir = image.Dir()

	err = command.RunLogging(cmd, log)
	if err != nil {
		// flushErr := log.Flush()
		// agg := errorz.NewAggregated(err, flushErr)
		return fmt.Errorf("Error removing container %s : %w", names, err)
	}

	return
}

func undeployContainers(binary string, images []resources.Image) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("rm", "")

	names := []string{}
	for _, image := range images {
		ctName, err := absContainerName(image)
		if err != nil {
			return err
		}
		names = append(names, ctName)
	}
	log.Info("Removing containers: %s ...", names)
	err = rmDockerContainers(log, binary, names...)
	if err != nil {
		flushErr := d.Flush()
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	log.Info("Removing finished")
	return
}

type DockerComposeProjectsDeployer struct {
	binary   string
	args     []string
	projects []resources.Project
}

func (d DockerComposeProjectsDeployer) Pull() (err error) {
	for _, project := range d.projects {
		err = pullDockerComposeProject(project, d.binary)
		if err != nil {
			return
		}
	}
	return
}

func (d DockerComposeProjectsDeployer) Deploy() (err error) {
	for _, project := range d.projects {
		err = upDockerComposeProject(project, d.binary, d.args...)
		if err != nil {
			return
		}
	}
	return
}

func (d DockerComposeProjectsDeployer) Undeploy(rmVolumes bool) (err error) {
	for _, p := range d.projects {
		err = downDockerComposeProject(p, d.binary, rmVolumes)
		if err != nil {
			return
		}
	}
	return
}

func pullDockerComposeProject(project resources.Project, binary string) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("pull", project.Name())
	log.Info("Pulling project: %s ...", project.Name())

	// Docker run level config
	projectVol := project.Dir() + ":/code:ro"
	runComposeOnDockerArgs := []string{"--rm",
		"-v", "/var/run/docker.sock:/var/run/docker.sock:rw", // Mount Docker daemon socket
		"-v", "/var/lib/docker/image:/var/lib/docker/image:rw", // Mount docker image cache
		"-v", projectVol, "--workdir", "/code", // Mount project code
	}

	config, errors := resources.MergedConfig(project)
	if errors != nil {
		return errors
	}
	// Supply environment at docker run level
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runComposeOnDockerArgs = append(runComposeOnDockerArgs, envArg)
	}

	// Set project name
	absoluteName, err := project.AbsoluteName()
	if err != nil {
		return err
	}

	// Up level config
	var cmdParams []string
	cmdParams = append(cmdParams, "--project-name", absoluteName)
	cmdParams = append(cmdParams, "pull")

	// Add default params
	//cmdParams = append(cmdParams,)

	err = runDockerImage(log, binary, runComposeOnDockerArgs, "", composeImage, cmdParams...)
	if err != nil {
		flushErr := d.Flush()
		err = fmt.Errorf("Error pulling project %s : %w", project.Name(), err)
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	log.Info("Pull finished for project: %s .", project.Name())
	return
}

func upDockerComposeProject(project resources.Project, binary string, args ...string) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("up", project.Name())
	log.Info("Upping project: %s ...", project.Name())

	// Docker run level config
	projectVol := project.Dir() + ":/code:ro"
	runComposeOnDockerArgs := []string{"--rm",
		"-v", "/var/run/docker.sock:/var/run/docker.sock:rw", // Mount Docker daemon socket
		"-v", "/var/lib/docker/image:/var/lib/docker/image:rw", // Mount docker image cache
		"-v", projectVol, "--workdir", "/code", // Mount project code
	}

	config, errors := resources.MergedConfig(project)
	if errors != nil {
		return errors
	}
	// Supply environment at docker run level
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runComposeOnDockerArgs = append(runComposeOnDockerArgs, envArg)
	}

	// compose interesting Options
	// --project-name
	// --profile (to deploy part of compose file) ?
	// --context (to deploy on different hosts) ?
	// --verbose
	// --log-level
	// --host (daemon socket to connect to)
	// --env-file ?

	// Set project name
	absoluteName, err := project.AbsoluteName()
	if err != nil {
		return err
	}

	// Up level config
	var cmdParams []string
	cmdParams = append(cmdParams, "--project-name", absoluteName)
	//cmdParams = append(cmdParams, "--verbose")
	cmdParams = append(cmdParams, "up")
	// up interesting Options
	// --force-recreate
	// --no-recreate
	// --no-build
	// --detach ?
	// --timeout ?
	// --renew-anon-volumes ?
	// --remove-orphans
	// --scale ?

	// Add default params
	cmdParams = append(cmdParams, "--detach", "--no-build", "--remove-orphans", "--force-recreate")

	cmdParams = append(cmdParams, args...)

	err = runDockerImage(log, binary, runComposeOnDockerArgs, "", composeImage, cmdParams...)
	if err != nil {
		flushErr := d.Flush()
		err = fmt.Errorf("Error upping project %s : %w", project.Name(), err)
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	log.Info("Up finished for project: %s .", project.Name())
	return
}

func downDockerComposeProject(project resources.Project, binary string, rmVolumes bool) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("down", project.Name())
	log.Info("Downing project: %s ...", project.Name())

	// Docker run level config
	projectVol := project.Dir() + ":/code:ro"
	runComposeOnDockerArgs := []string{"--rm",
		"-v", "/var/run/docker.sock:/var/run/docker.sock:rw", // Mount Docker daemon socket
		"-v", "/var/lib/docker/image:/var/lib/docker/image:rw", // Mount docker image cache
		"-v", projectVol, "--workdir", "/code", // Mount project code
	}

	config, errors := resources.MergedConfig(project)
	if errors != nil {
		return errors
	}
	// Supply environment at docker run level
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runComposeOnDockerArgs = append(runComposeOnDockerArgs, envArg)
	}

	// Set project name
	absoluteName, err := project.AbsoluteName()
	if err != nil {
		return err
	}

	// Up level config
	var cmdParams []string
	cmdParams = append(cmdParams, "--project-name", absoluteName)
	cmdParams = append(cmdParams, "down")

	// Add default params
	cmdParams = append(cmdParams, "--remove-orphans", "--timeout", defaultComposeDownTimeoutInSec)

	if rmVolumes {
		cmdParams = append(cmdParams, "--volumes")
	}

	err = runDockerImage(log, binary, runComposeOnDockerArgs, "", composeImage, cmdParams...)
	if err != nil {
		flushErr := d.Flush()
		err = fmt.Errorf("Error downing project %s : %w", project.Name(), err)
		agg := errorz.NewAggregated(err, flushErr)
		return agg
	}

	log.Info("Down finished for project: %s .", project.Name())
	return
}
