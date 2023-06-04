package pod

import (
	k8s "k8s.io/api"
	apimachinery "k8s.io/apimachinery/pkg"
)

type Executor {
	config []string
	translator Translator
}

func (e Executor) exec(commands []string, retries int) (status int, err error) {

}

func (e Executor) Create(namespace string, resource any) (err error) {
	switch res := resource.(type) {
	case k8s.core.v1.Pod:
		return e.createPod(namespace, res)
	}
}

func (e Executor) createPod(namespace string, pod k8s.core.v1.Pod) (err error) {
	cmds, err := translator.createNetworkOwnerContainer(pod)
	if err != nil {
		return err
	}
	status, err := e.exec(cmds, -1)

	for _, volume := pod.Spec.Volumes {
		cmds, err := translator.createVolume(namespace, volume)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}

	for _, container := pod.Spec.InitContainers {
		cmds, err := translator.createContainer(namespace, container)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}

	for _, container := pod.Spec.Containers {
		cmds, err := translator.createContainer(namespace, container)
		if err != nil {
			return err
		}
		status, err := e.exec(cmds, -1)
	}
}

type Translator {
	binary string
}

func (t Translator) forgeResName(namespace string, resource any) (name string, err error) {
	var resName string
	switch res := resource.(type) {
	case k8s.core.v1.Pod:
		resName = res.ObjectMeta.Name
	case k8s.core.v1.Volume:
		resName = res.Name
	default:
		err = fmt.Errorf("Cannot forge a name for unknown type: %T !", resource)
	}

	name = fmt.Sprintf("%s_%s", namespace, resName)
}

func (t Translator) createNetworkOwnerContainer(namespace string, pod k8s.core.v1.Pod) (cmds []string, err error) {
	podName := pod.ObjectMeta.Name
	podMainCtCpus := "0.05"
	podMainCtMemory := "64m"
	networkName := ""
	addHostRules := ""
	pauseImage := ""
	cmd := fmt.Sprintf(	"docker run -d --name '%s' --cpus=%s --memory=%s --memory-swap=%s --memory-swappiness=0" + 
						" --network '%s' %s --restart=always '%s' /bin/sleep inf",
						podName, podMainCtCpus, podMainCtMemory, podMainCtMemory, networkName, addHostRules, pauseImage
	)
	append(cmds, cmd)
}

func (t Translator) createVolume(namespace string, vol k8s.core.v1.Volume) (cmds []string, err error) {
	if vol.VolumeSource.HostPath != nil {
		return createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return createEmptyDirPodVolume(namespace, vol)
	}
	err = fmt.Errorf("Not supported volume type for volume: %s !", vol.Name)
}

func (t Translator) createHostPathPodVolume(namespace string, vol k8s.core.v1.Volume) (cmds []string, err error) {
	hostPathType := vol.HostPathVolumeSource.HostPathType
	if hostPathType != "" {
		err = fmt.Errorf("Not supported HostPathType: %s for volume: %s !", hostPathType, vol.Name)
	}

	name := t.forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPathVolumeSource.Path
	cmd := fmt.Sprintf("docker volume create --driver local -o o=bind -o type=none -o device=%s %s", path, name)
	append(cmds, cmd)
}

func (t Translator) createEmptyDirPodVolume(namespace string, vol k8s.core.v1.Volume) (cmds []string, err error) {
	if vol.
	name := t.forgeResName(namespace, vol)
	cmd := fmt.Sprintf("docker volume create --driver local %s", name)
	append(cmds, cmd)
}

func (t Translator) createContainer(podName string, restartPolicy k8s.core.v1.RestartPolicy, container k8s.core.v1.Container) (cmds []string, err error) {
	ctName := fmt.Sprintf("%s-%s", podName, container.Name)
	image := container.Image
	privileged := container.SecurityContext.Privileged
	tty != container.TTY
	workingDir := container.WorkingDir
	pullPolicy := container.ImagePullPolicy
	volumeMounts := container.VolumeMounts
	env := container.Env
	entrypoint := container.Command
	args := container.Args
	cpuLimitInMilli := container.Resources.Limits.Cpu().MilliValue()
	memroryRequestInMega := container.Resources.Requests.Memory().AsScale(apimachinery.api.resource.Mega)
	memoryLimitInMega := container.Resources.Limits.Memory().AsScale(apimachinery.api.resource.Mega)

	var runArgs []string
	var resourcesArgs []string
	var envArgs []string
	var cmdArgs []string

	if privileged {
		append(runArgs, "--privileged")
	}
	if tty {
		append(runArgs, "-t")
	}
	if workingDir != "" {
		append(runArgs, fmt.Sprintf("--workdir=%s", workingDir))
	}
	
	if cpuLimitInMilli > 0 {
		append(resourcesArgs, "--cpu-period=100000")
		append(resourcesArgs, fmt.Sprintf("--cpu-quota=%d00", cpuLimitInMilli))
	}

	//append(resourcesArgs, "--memory-swappiness=0")

	if memoryLimitInMega > 0 {
		append(resourcesArgs, fmt.Sprintf("--memory=%dm", memoryLimitInMega))
		append(resourcesArgs, "--memory-swap=-1")
	}

	if memroryRequestInMega > 0 {
		append(resourcesArgs, fmt.Sprintf("--memory-reservation=%dm", memroryRequestInMega))
	}

	dockerRestartPolicy := ""
	switch restartPolicy {
	case "Always":
		dockerRestartPolicy = "always"
	case "Never":
		dockerRestartPolicy = "no"
	case "OnFailure":
		dockerRestartPolicy = "on-failure"
	default:
		err = fmt.Errorf("No supported restart policy: %s in container: %s !", restartPolicy, ctName)
		return
	}

	append(runArgs, "--restart=$dockerRestartPolicy")

	if len(volumeMounts) > 0 {
		for _, volMount := range volumeMounts {
			volumeName := forgeResName(namespace, volMount)
			mountPath := volMount.MountPath
			readOnly := volMount.ReadOnly
			mountOpts := "rw"
			if readOnly {
				mountOpts = "ro"
			}
			append(runArgs, "-v")
			append(runArgs, fmt.Sprintf("%s:%s:%s", volumeName, mountPath, mountOpts))
		}
	}

	// Pass first entrypoint item as docker run uniq entrypoint
	if len(entrypoint) > 0 {
		append(runArgs, "--entrypoint")
		append(runArgs, entrypoint[0])
	}

	// Pass folowing entrypoint items as docker run command args
	if len(entrypoint) > 1 {
		append(cmdArgs, entrypoint[1:])
	}

	if len(args) > 0 [
		append(cmdArgs, args)
	]

	if len(env) > 0 {
		for _, envVar := range env {
			append(envArgs, "-e")
			append(envArgs, fmt.Sprintf("\"%s=%s\"", envVar.Name, envVar.Value))
		}
	}

	/*
        # Test if ct already started or start it if excited or create it
        cmd="docker ps --format '{{ .Names }}' | grep -w '$ctName' || docker ps -f 'status=created' -f 'status=exited' --format '{{ .Names }}' | grep -w '$ctName' && docker start '$ctName' || docker run -d --name \"$ctName\" $resourcesArgs --network 'container:$podName' $runArgs $envArgs \"$image\" $entrypointArgs $cmdArgs"
        ! $DEBUG && >&2 echo "- $cmd"
        ! $DEBUG || >&2 echo "Running container $podName:$name ..."
        echo "$cmd"

        # FIXME: should write containers.txt once on pod creation
        #cmd="docker exec -u=0 '$podName' /bin/sh -c 'echo $ctName >> /containers.txt'"
        #! $DEBUG && >&2 echo "- $cmd"
        #echo "$cmd"
	*/

}
