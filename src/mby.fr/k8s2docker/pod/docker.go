package pod

import (
	"fmt"
	"log"

	"strings"

	"mby.fr/utils/cmdz"
	//"mby.fr/utils/promise"

	k8sv1 "k8s.io/api/core/v1"
)

type Executor struct {
	translator Translator
	config     []string
	forkCount  int
}

func (e Executor) exec(execs ...*cmdz.Exec) error {
	statuses, err := cmdz.ParallelRunAll(e.forkCount, execs...)

	if err != nil {
		return err
	}
	if cmdz.Failed(statuses...) {
		errorMessages := ""
		for idx, status := range statuses {
			if status != 0 {
				//stdout := execs[idx].StdoutRecord.String()
				stderr := execs[idx].StderrRecord.String()
				errorMessages += fmt.Sprintf("\nRC=%d ERR> %s", status, strings.TrimSpace(stderr))
			}
		}
		return fmt.Errorf("Encountered some execution failure: %s", errorMessages)
	}
	return nil
}

func (e Executor) Create(namespace string, resource any) (err error) {
	switch res := resource.(type) {
	case k8sv1.Pod:
		err = e.runPod(namespace, res)
	default:
		err = fmt.Errorf("Cannot create not supported type %T", res)
	}
	return
}

func (e Executor) runPod(namespace string, pod k8sv1.Pod) (err error) {
	podName := e.translator.forgeResName(namespace, pod)

	execs, err := e.translator.createNetworkOwnerContainer(namespace, pod)
	if err != nil {
		return err
	}
	err = e.exec(execs...)
	if err != nil {
		return fmt.Errorf("Unable to create Network owner container for pod: %s. Caused by: %w", podName, err)
	}

	for _, volume := range pod.Spec.Volumes {
		execs, err := e.translator.createVolume(namespace, volume)
		if err != nil {
			return err
		}
		err = e.exec(execs...)
		if err != nil {
			volName := e.translator.forgeResName(podName, volume)
			return fmt.Errorf("Unable to create volume %s. Caused by: %w", volName, err)
		}
	}

	for _, container := range pod.Spec.InitContainers {
		execs, err := e.translator.createContainer(namespace, pod, container)
		if err != nil {
			return err
		}
		err = e.exec(execs...)
		if err != nil {
			return fmt.Errorf("Unable to run Init Containers for pod %s. Caused by: %w", podName, err)
		}
	}

	for _, container := range pod.Spec.Containers {
		execs, err := e.translator.createContainer(namespace, pod, container)
		if err != nil {
			return err
		}
		err = e.exec(execs...)
		if err != nil {
			return fmt.Errorf("Unable to run Containers for pod %s. Caused by: %w", podName, err)
		}
	}
	return
}

type Translator struct {
	binary string
}

func (t Translator) forgeResName(prefix string, resource any) (name string) {
	var resName string
	switch res := resource.(type) {
	case k8sv1.Pod:
		resName = res.ObjectMeta.Name
	case k8sv1.Volume:
		resName = res.Name
	case k8sv1.VolumeMount:
		resName = res.Name
	case k8sv1.Container:
		resName = res.Name
	default:
		log.Fatalf("Cannot forge a name for unknown type: %T !", resource)
	}

	name = fmt.Sprintf("%s_%s", prefix, resName)
	return
}

func (t Translator) createNetworkOwnerContainer(namespace string, pod k8sv1.Pod) (cmds []*cmdz.Exec, err error) {
	podName := t.forgeResName(namespace, pod)
	ctName := fmt.Sprintf("%s_root", podName)
	cpusArgs := "--cpus=0.05"
	memoryArgs := "--memory=64m"
	//swapArgs := "--memory-swap=128m"
	networkName := fmt.Sprintf("%s_net", podName)
	addHostRules := ""
	pauseImage := "alpine:3.17.3"

	networkArgs := []string{"network", "create", networkName}

	runArgs := []string{"run", "-d", "--name", ctName, "--restart=always", "--network", networkName,
		cpusArgs, memoryArgs} //"--memory-swappiness=0"

	if addHostRules != "" {
		runArgs = append(runArgs, addHostRules)
	}

	runArgs = append(runArgs, pauseImage)
	runArgs = append(runArgs, "/bin/sleep", "inf")

	cmdNetwork := cmdz.Execution(t.binary, networkArgs...)
	cmdRun := cmdz.Execution(t.binary, runArgs...)
	cmds = append(cmds, cmdNetwork, cmdRun)
	return
}

func (t Translator) createVolume(namespace string, vol k8sv1.Volume) (cmds []*cmdz.Exec, err error) {
	if vol.VolumeSource.HostPath != nil {
		return t.createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return t.createEmptyDirPodVolume(namespace, vol)
	}
	err = fmt.Errorf("Not supported volume type for volume: %s !", vol.Name)
	return
}

func (t Translator) createHostPathPodVolume(namespace string, vol k8sv1.Volume) (cmds []*cmdz.Exec, err error) {
	hostPathType := *vol.VolumeSource.HostPath.Type
	if hostPathType != k8sv1.HostPathUnset {
		err = fmt.Errorf("Not supported HostPathType: %s for volume: %s !", hostPathType, vol.Name)
	}

	name := t.forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPath.Path
	cmd := cmdz.Execution(t.binary, "volume", "create", "--driver", "local")
	cmd.AddArgs("-o", "o=bind", "-o", "type=none", "-o", "device="+path)
	cmd.AddArgs(name)
	cmds = append(cmds, cmd)
	return
}

func (t Translator) createEmptyDirPodVolume(namespace string, vol k8sv1.Volume) (cmds []*cmdz.Exec, err error) {
	if vol.VolumeSource.EmptyDir == nil {
		err = fmt.Errorf("Bad EmptyDirVolume !")
		return
	}
	name := t.forgeResName(namespace, vol)
	cmd := cmdz.Execution(t.binary, "volume", "create", "--driver", "local", name)
	cmds = append(cmds, cmd)
	return
}

func (t Translator) createContainer(namespace string, pod k8sv1.Pod, container k8sv1.Container) (cmds []*cmdz.Exec, err error) {
	podName := t.forgeResName(namespace, pod)
	ctName := t.forgeResName(podName, container)
	image := container.Image
	privileged := *container.SecurityContext.Privileged
	tty := container.TTY
	workingDir := container.WorkingDir
	restartPolicy := pod.Spec.RestartPolicy
	pullPolicy := container.ImagePullPolicy
	volumeMounts := container.VolumeMounts
	env := container.Env
	entrypoint := container.Command
	args := container.Args
	cpuLimitInMilli := container.Resources.Limits.Cpu().MilliValue()
	//memroryRequestInMega, _ := container.Resources.Requests.Memory().AsScale(k8smachineryresource.Mega)
	//memoryLimitInMega, _ := container.Resources.Limits.Memory().AsScale(k8smachineryresource.Mega)
	memroryRequestInMega, _ := container.Resources.Requests.Memory().AsInt64()
	memoryLimitInMega, _ := container.Resources.Limits.Memory().AsInt64()
	labels := pod.ObjectMeta.Labels

	var runArgs []string
	var resourcesArgs []string
	var envArgs []string
	var cmdArgs []string
	var labelArgs []string

	runArgs = append(runArgs, "--rm")
	runArgs = append(runArgs, "--name")
	runArgs = append(runArgs, ctName)

	if privileged {
		runArgs = append(runArgs, "--privileged")
	}
	if tty {
		runArgs = append(runArgs, "-t")
	}
	if workingDir != "" {
		runArgs = append(runArgs, fmt.Sprintf("--workdir=%s", workingDir))
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
	runArgs = append(runArgs, fmt.Sprintf("--restart=%s", dockerRestartPolicy))

	dockerPullPolicy := ""
	switch pullPolicy {
	case k8sv1.PullAlways:
		dockerPullPolicy = "always"
	case k8sv1.PullNever:
		dockerPullPolicy = "never"
	case k8sv1.PullIfNotPresent:
		dockerPullPolicy = "missing"
	default:
		err = fmt.Errorf("No supported pull policy: %s in container: %s !", pullPolicy, ctName)
		return
	}
	runArgs = append(runArgs, fmt.Sprintf("--pull=%s", dockerPullPolicy))

	if len(volumeMounts) > 0 {
		for _, volMount := range volumeMounts {
			volumeName := t.forgeResName(namespace, volMount)
			mountPath := volMount.MountPath
			readOnly := volMount.ReadOnly
			mountOpts := "rw"
			if readOnly {
				mountOpts = "ro"
			}
			runArgs = append(runArgs, "-v")
			runArgs = append(runArgs, fmt.Sprintf("%s:%s:%s", volumeName, mountPath, mountOpts))
		}
	}

	// Pass first entrypoint item as docker run uniq entrypoint
	if len(entrypoint) > 0 {
		runArgs = append(runArgs, "--entrypoint")
		runArgs = append(runArgs, entrypoint[0])
	}

	// Pass folowing entrypoint items as docker run command args
	if len(entrypoint) > 1 {
		cmdArgs = append(cmdArgs, entrypoint[1:]...)
	}

	if len(args) > 0 {
		cmdArgs = append(cmdArgs, args...)
	}

	if cpuLimitInMilli > 0 {
		resourcesArgs = append(resourcesArgs, "--cpu-period=100000")
		resourcesArgs = append(resourcesArgs, fmt.Sprintf("--cpu-quota=%d00", cpuLimitInMilli))
	}

	//append(resourcesArgs, "--memory-swappiness=0")

	if memoryLimitInMega > 0 {
		resourcesArgs = append(resourcesArgs, fmt.Sprintf("--memory=%dm", memoryLimitInMega))
		resourcesArgs = append(resourcesArgs, "--memory-swap=-1")
	}

	if memroryRequestInMega > 0 {
		resourcesArgs = append(resourcesArgs, fmt.Sprintf("--memory-reservation=%dm", memroryRequestInMega))
	}

	if len(env) > 0 {
		for _, envVar := range env {
			envArgs = append(envArgs, "-e")
			envArgs = append(envArgs, fmt.Sprintf("\"%s=%s\"", envVar.Name, envVar.Value))
		}
	}

	for key, value := range labels {
		labelArgs = append(labelArgs, "-l")
		labelArgs = append(labelArgs, fmt.Sprintf("%s=%s", key, value))
	}

	// TODO: add annotations ?

	cmd := cmdz.Execution(t.binary, "run")
	cmd.AddArgs(runArgs...)
	cmd.AddArgs(resourcesArgs...)
	cmd.AddArgs(envArgs...)
	cmd.AddArgs(labelArgs...)
	cmd.AddArgs(image)
	cmd.AddArgs(cmdArgs...)
	cmds = append(cmds, cmd)

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

	return
}
