package deploy

import (
	//"bytes"
	"fmt"
	"os/exec"

	"mby.fr/mass/internal/command"
	"mby.fr/mass/internal/display"
	"mby.fr/mass/internal/logger"
	"mby.fr/mass/internal/resources"
)

const composeImage = "docker/compose:alpine-1.29.2"

var NotDeployableResource error = fmt.Errorf("Not deployable resource")

type Deployer interface {
	Deploy() error
}

func New(r resources.Resource) (Deployer, error) {
	switch res := r.(type) {
	case resources.Project:
		return DockerComposeProjectsDeployer{"docker", []string{}, []resources.Project{res}}, nil
	case resources.Image:
		return DockerImagesDeployer{"docker", []string{}, []resources.Image{res}}, nil
	default:
		return nil, fmt.Errorf("%w: %s", NotDeployableResource, r.AbsoluteName())
	}
}

type DockerImagesDeployer struct {
	binary string
	args   []string
	images []resources.Image
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
func runDockerImage(binary string, image string, log logger.ActionLogger, args ...string) (err error) {
	log.Info("Running image: %s ...", image)

	var runParams []string
	runParams = append(runParams, "run")

	runParams = append(runParams, args...)

	log.Debug("run params: %s", runParams)
	cmd := exec.Command(binary, runParams...)
	//cmd.Dir = image.Dir()

	err = command.RunLogging(cmd, log)
	if err != nil {
		log.Flush()
		return fmt.Errorf("Error running image %s : %w", image, err)
	}
	log.Info("Run finished for image: %s .", image)
	return
}

func runImage(binary string, image resources.Image) (err error) {
	d := display.Service()
	log := d.BufferedActionLogger("run", image.FullName())

	var runParams []string
	runParams = append(runParams, "run")

	config, errors := resources.MergedConfig(image)
	if errors != nil {
		return errors
	}
	// Add envVars
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runParams = append(runParams, envArg)
	}

	// Add image full name
	runParams = append(runParams, image.FullName())

	// Add runParams
	for _, argValue := range config.RunArgs {
		runParams = append(runParams, argValue)
	}

	err = runDockerImage(binary, image.FullName(), log, runParams...)

	return
}

type DockerComposeProjectsDeployer struct {
	binary   string
	args     []string
	projects []resources.Project
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

	runComposeOnDockerArgs = append(runComposeOnDockerArgs, composeImage)

	config, errors := resources.MergedConfig(project)
	if errors != nil {
		return errors
	}
	// Supply environment at docker run level
	for argKey, argValue := range config.Environment {
		var envArg string = "-e=" + argKey + "=" + argValue
		runComposeOnDockerArgs = append(runComposeOnDockerArgs, envArg)
	}

	// Compose level congig
	var composeParams []string
	composeParams = append(composeParams, args...)

	// compose interesting Options
	// --project-name
	// --profile (to deploy part of compose file) ?
	// --context (to deploy on different hosts) ?
	// --verbose
	// --log-level
	// --host (daemon socket to connect to)
	// --env-file ?

	// Set project name
	composeParams = append(composeParams, "--project-name", project.Name())
	//composeParams = append(composeParams, "--verbose")

	// Up level config
	var upParams []string
	upParams = append(upParams, "up")
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
	upParams = append(upParams, "--detach", "--no-build", "--remove-orphans", "--force-recreate")

	allParams := append(runComposeOnDockerArgs, composeParams...)
	allParams = append(allParams, upParams...)
	err = runDockerImage(binary, composeImage, log, allParams...)

	if err != nil {
		log.Flush()
		return fmt.Errorf("Error upping project %s : %w", project.Name(), err)
	}
	log.Info("Up finished for project: %s .", project.Name())
	return
}
