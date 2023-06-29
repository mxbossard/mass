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

func (e Executor) Create(namespace string, resource any) (err error) {
	switch res := resource.(type) {
	case k8sv1.Pod:
		err = e.createPod(namespace, res)
	default:
		err = fmt.Errorf("Cannot create not supported type %T", res)
	}
	return
}

func (e Executor) createPod(namespace string, pod k8sv1.Pod) (err error) {
	podName := forgeResName(namespace, pod)

	netId, err := e.translator.podNetworkId(namespace, pod)
	if err != nil {
		return err
	}
	if netId == "" {
		execs, err := e.translator.createPodNetwork(namespace, pod)
		if err != nil {
			return err
		}
		err = cmdz.Sequential(execs...)
		if err != nil {
			return fmt.Errorf("Unable to create Network for pod: %s. Caused by: %w", podName, err)
		}
	}

	ctId, err := e.translator.podRootContainerId(namespace, pod)
	if err != nil {
		return err
	}
	if ctId == "" {
		execs, err := e.translator.createPodRootContainer(namespace, pod)
		if err != nil {
			return err
		}
		err = cmdz.Sequential(execs...)
		if err != nil {
			return fmt.Errorf("Unable to create Root Container for pod: %s. Caused by: %w", podName, err)
		}
	}

	for _, volume := range pod.Spec.Volumes {
		volId, err := e.translator.volumeId(namespace, volume)
		if err != nil {
			return err
		}
		if volId == "" {
			execs, err := e.translator.createVolume(namespace, volume)
			if err != nil {
				return err
			}
			err = cmdz.Parallel(e.forkCount, execs...)
			if err != nil {
				volName := forgeResName(podName, volume)
				return fmt.Errorf("Unable to create volume %s. Caused by: %w", volName, err)
			}
		}
	}

	var ictExecs []*cmdz.Exec
	for _, container := range pod.Spec.InitContainers {
		ctId, err := e.translator.podContainerId(namespace, pod, container)
		if err != nil {
			return err
		}
		if ctId != "" {
			ctName := podContainerName(namespace, pod, container)
			err = fmt.Errorf("Init Container %s already exists !", ctName)
			return err
		} else {
			execs, err := e.translator.createPodContainer(namespace, pod, container, true)
			if err != nil {
				return err
			}
			ictExecs = append(ictExecs, execs...)
		}
	}
	if len(ictExecs) > 0 {
		err = cmdz.Parallel(e.forkCount, ictExecs...)
		if err != nil {
			return fmt.Errorf("Unable to run Init Containers for pod %s. Caused by: %w", podName, err)
		}
	}

	var ctExecs []*cmdz.Exec
	for _, container := range pod.Spec.Containers {
		ctId, err := e.translator.podContainerId(namespace, pod, container)
		if err != nil {
			return err
		}
		if ctId == "" {
			execs, err := e.translator.createPodContainer(namespace, pod, container, false)
			if err != nil {
				return err
			}
			ctExecs = append(ctExecs, execs...)
		}
	}
	if len(ctExecs) > 0 {
		err = cmdz.Parallel(e.forkCount, ctExecs...)
		if err != nil {
			return fmt.Errorf("Unable to run Containers for pod %s. Caused by: %w", podName, err)
		}
	}

	return
}

func networkName(namespace string, pod k8sv1.Pod) string {
	podName := forgeResName(namespace, pod)
	networkName := fmt.Sprintf("%s_net", podName)
	return networkName
}

func podRootContainerName(namespace string, pod k8sv1.Pod) string {
	podName := forgeResName(namespace, pod)
	ctName := fmt.Sprintf("%s_root", podName)
	return ctName
}

func podContainerName(namespace string, pod k8sv1.Pod, container k8sv1.Container) string {
	podName := forgeResName(namespace, pod)
	ctName := forgeResName(podName, container)
	return ctName
}

func forgeResName(prefix string, resource any) (name string) {
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

func escapeCmdArg(arg string) string {
	if ! strings.Contains(arg, " ") {
		return arg
	}
	
	// If spaces in arg
	if ! strings.Contains(arg, `"`) {
		return `"` + arg + `"`
	} else if ! strings.Contains(arg, "'") {
		return "'" + arg + "'"
	}
	escapedArg := strings.Replace(arg, `"`, `\"`, -1)
	return `"` + escapedArg + `"`
}

type Translator struct {
	binary string
}

func (t Translator) podPhase(namespace string, pod k8sv1.Pod) (k8sv1.PodePhase, error) { 
 	// TODO
}

func (t Translator) containerState(namespace string, pod k8sv1.Pod, container k8sv1.Container) (k8sv1.ContainerState, error) { 
	// TODO
}

func (t Translator) podNetworkId(namespace string, pod k8sv1.Pod) (string, error) { 
	networkName := networkName(namespace, pod)
	script := fmt.Sprintf("%s network ls -f 'Name=%s' -q", t.binary, networkName)
	exec := cmdz.Execution("/bin/sh", "-c", script)
	stdout := &strings.Builder{}
	exec.RecordingOutputs(stdout, nil)
	exec.Retries = 2

	err := cmdz.Sequential(exec)
	if err != nil {
		return "", err
	}
	id := stdout.String()
	return id, nil
}

func (t Translator) createPodNetwork(namespace string, pod k8sv1.Pod) (cmds []*cmdz.Exec, err error) {
	networkName := networkName(namespace, pod)

	networkArgs := []string{"network", "create", networkName}

	cmd := cmdz.Execution(t.binary, networkArgs...)
	cmd.Retries = 2
	cmds = append(cmds, cmd)
	return
}

func (t Translator) podRootContainerId(namespace string, pod k8sv1.Pod) (string, error) { 
	ctName := podRootContainerName(namespace, pod)
	script := fmt.Sprintf("%s ps -f 'Name=%s' -q", t.binary, ctName)
	exec := cmdz.Execution("/bin/sh", "-c", script)
	stdout := &strings.Builder{}
	exec.RecordingOutputs(stdout, nil)
	exec.Retries = 2

	err := cmdz.Sequential(exec)
	if err != nil {
		return "", err
	}
	id := stdout.String()
	return id, nil
}

func (t Translator) createPodRootContainer(namespace string, pod k8sv1.Pod) (cmds []*cmdz.Exec, err error) {
	ctName := podRootContainerName(namespace, pod)
	cpusArgs := "--cpus=0.05"
	memoryArgs := "--memory=64m"
	networkName := networkName(namespace, pod)
	addHostRules := ""
	pauseImage := "alpine:3.17.3"

	runArgs := []string{"run", "-d", "--name", ctName, "--restart=always", "--network", networkName,
		cpusArgs, memoryArgs} //"--memory-swappiness=0"

	if addHostRules != "" {
		runArgs = append(runArgs, addHostRules)
	}

	runArgs = append(runArgs, pauseImage)
	runArgs = append(runArgs, "/bin/sleep", "inf")

	cmd := cmdz.Execution(t.binary, runArgs...)
	cmd.Retries = 2
	cmds = append(cmds, cmd)
	return
}

func (t Translator) volumeId(namespace string, vol k8sv1.Volume) (string, error) { 
	volName := forgeResName(namespace, vol)
	script := fmt.Sprintf("%s volume ls -f 'Name=%s' -q", t.binary, volName)
	exec := cmdz.Execution("/bin/sh", "-c", script)
	stdout := &strings.Builder{}
	exec.RecordingOutputs(stdout, nil)
	exec.Retries = 2

	err := cmdz.Sequential(exec)
	if err != nil {
		return "", err
	}
	id := stdout.String()
	return id, nil
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

	name := forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPath.Path
	cmd := cmdz.Execution(t.binary, "volume", "create", "--driver", "local")
	cmd.Retries = 2
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
	name := forgeResName(namespace, vol)
	cmd := cmdz.Execution(t.binary, "volume", "create", "--driver", "local", name)
	cmd.Retries = 2
	cmds = append(cmds, cmd)
	return
}

func (t Translator) podContainerId(namespace string, pod k8sv1.Pod, container k8sv1.Container) (string, error) { 
	ctName := podContainerName(namespace, pod, container)
	script := fmt.Sprintf("%s ps -f 'Name=%s' -q", t.binary, ctName)
	exec := cmdz.Execution("/bin/sh", "-c", script)
	stdout := &strings.Builder{}
	exec.RecordingOutputs(stdout, nil)
	exec.Retries = 2

	err := cmdz.Sequential(exec)
	if err != nil {
		return "", err
	}
	id := stdout.String()
	return id, nil
}

func (t Translator) createPodContainer(namespace string, pod k8sv1.Pod, container k8sv1.Container, init bool) (cmds []*cmdz.Exec, err error) {
	ctName := podContainerName(namespace, pod, container)
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

	if init {
		restartPolicy = k8sv1.RestartPolicyNever
		runArgs = append(runArgs, "--rm")
	}

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
			volumeName := forgeResName(namespace, volMount)
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
		for _, arg := range(entrypoint[1:]) {
			cmdArgs = append(cmdArgs, arg)
		} 
	}

	if len(args) > 0 {
		for _, arg := range(args) {
			cmdArgs = append(cmdArgs, arg)
		} 
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
	cmd.Retries = 2
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
