package pod

import (
	"fmt"
	"log"

	"encoding/json"
	"strings"

	"mby.fr/utils/cmdz"
	//"mby.fr/utils/promise"
	"mby.fr/k8s2docker/compare"
	"mby.fr/k8s2docker/descriptor"

	k8sv1 "k8s.io/api/core/v1"
)

type Executor struct {
	translator Translator
	config     []string
	forkCount  int
}

func (e Executor) Apply(namespace string, resource any) (err error) {
	switch res := resource.(type) {
	case k8sv1.Pod:
		err = e.createPod(namespace, res)
	default:
		err = fmt.Errorf("Cannot create not supported type %T", res)
	}
	return

}

func (e Executor) deletePod(namespace string, pod k8sv1.Pod) (err error) {
	//TODO
	return
}

func (e Executor) createPod(namespace string, pod k8sv1.Pod) (err error) {
	podName := forgeResName(namespace, pod)

	netId, err := e.translator.podNetworkId(namespace, pod)
	if err != nil {
		return err
	}
	if netId == "" {
		exec, err := e.translator.createPodNetwork(namespace, pod)
		if err != nil {
			return err
		}
		_, err = exec.BlockRun()
		if err != nil {
			return fmt.Errorf("Unable to create Network for pod: %s. Caused by: %w", podName, err)
		}
	}

	execs, err := e.translator.createPodRootContainer(namespace, pod)
	if err != nil {
		return err
	}
	_, err = execs.BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create Root Container for pod: %s. Caused by: %w", podName, err)
	}

	volExec := cmdz.Parallel().ErrorOnFailure(true)
	for _, volume := range pod.Spec.Volumes {
		volId, err := e.translator.volumeId(namespace, volume)
		if err != nil {
			return err
		}
		if volId == "" {
			exec, err := e.translator.createVolume(namespace, volume)
			if err != nil {
				return err
			}
			volExec.Add(exec)
		}
	}
	_, err = volExec.BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create a volume ! Caused by: %w", err)
	}

	ictExec := cmdz.Parallel().ErrorOnFailure(true)
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
			exec, err := e.translator.createPodContainer(namespace, pod, container, true)
			if err != nil {
				return err
			}
			ictExec.Add(exec)
		}
	}
	/*
		if len(ictExecs) > 0 {
			err = cmdz.BlockParallel(e.forkCount, ictExecs...)
			if err != nil {
				return fmt.Errorf("Unable to run Init Containers for pod %s. Caused by: %w", podName, err)
			}
		}*/
	_, err = ictExec.BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create init containers ! Caused by: %w", err)
	}

	ctExec := cmdz.Parallel().ErrorOnFailure(true)
	for _, container := range pod.Spec.Containers {
		ctId, err := e.translator.podContainerId(namespace, pod, container)
		if err != nil {
			return err
		}
		if ctId == "" {
			exec, err := e.translator.createPodContainer(namespace, pod, container, false)
			if err != nil {
				return err
			}
			ctExec.Add(exec)
		}
	}
	/*
		if len(ctExecs) > 0 {
			err = cmdz.BlockParallel(e.forkCount, ctExecs...)
			if err != nil {
				return fmt.Errorf("Unable to run Containers for pod %s. Caused by: %w", podName, err)
			}
		}*/
	_, err = ctExec.BlockRun()
	if err != nil {
		return fmt.Errorf("Unable to create containers ! Caused by: %w", err)
	}

	return
}

func (e Executor) applyPod(namespace string, pod k8sv1.Pod) (err error) {
	// TODO: how to use a kubelet ?
	// Kubelet should be responsible for create / update / delete

	// TODO : check for pod existance. If it already exists and phase is ok check for "updateness"
	// THEN update pod or delete/create pod

	ctId, err := e.translator.podRootContainerId(namespace, pod)
	if err != nil {
		return err
	}
	if ctId == "" {
		// Create pod
	} else {
		podPhase, err := e.translator.getCommitedPodPhase(namespace, pod)
		if err != nil {
			log.Printf("Swallowed error: %s", err)
			// TODO recreate pod
		}
		// TODO what to do with podPhase ?
		_ = podPhase

		commitedPod, err := e.translator.getCommitedPod(namespace, pod)
		if err != nil {
			log.Printf("Swallowed error: %s", err)
			// TODO recreate pod
		}
		if podPhase != nil && commitedPod != nil {
			diff := compare.ComparePods(*commitedPod, pod)
			if diff.IsUpdatable() {
				// TODO update pod
			}
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
	if !strings.Contains(arg, " ") {
		return arg
	}

	// If spaces in arg
	if !strings.Contains(arg, `"`) {
		return `"` + arg + `"`
	} else if !strings.Contains(arg, "'") {
		return "'" + arg + "'"
	}
	escapedArg := strings.Replace(arg, `"`, `\"`, -1)
	return `"` + escapedArg + `"`
}

type Translator struct {
	binary string
}

func (t Translator) commitPod(namespace string, pod k8sv1.Pod) (cmdz.Executer, error) {
	b, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", fmt.Sprintf("echo '%s' > /podConfig.txt", string(b))}
	exec := cmdz.Cmd(t.binary, execArgs...)
	return exec, nil
}

func (t Translator) getCommitedPod(namespace string, pod k8sv1.Pod) (*k8sv1.Pod, error) {
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", "cat /podConfig.txt"}
	exec := cmdz.Cmd(t.binary, execArgs...).ErrorOnFailure(true)
	_, err := exec.BlockRun()
	if err != nil {
		return nil, err
	}
	jsonBlob := exec.StdoutRecord()
	p, err := descriptor.LoadPod([]byte(jsonBlob))
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (t Translator) commitPodPhase(namespace string, pod k8sv1.Pod, phase k8sv1.PodPhase) cmdz.Executer {
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", fmt.Sprintf("echo %s > /podPhase.txt", phase)}
	exec := cmdz.Cmd(t.binary, execArgs...)
	return exec
}

func (t Translator) getCommitedPodPhase(namespace string, pod k8sv1.Pod) (*k8sv1.PodPhase, error) {
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", "cat /podPhase.txt"}
	exec := cmdz.Cmd(t.binary, execArgs...).ErrorOnFailure(true)
	_, err := exec.BlockRun()
	if err != nil {
		return nil, err
	}
	p := k8sv1.PodPhase(exec.StdoutRecord())
	return &p, nil
}

func (t Translator) containerState(namespace string, pod k8sv1.Pod, container k8sv1.Container) (p k8sv1.ContainerState, e error) {
	return
}

func (t Translator) podNetworkId(namespace string, pod k8sv1.Pod) (string, error) {
	networkName := networkName(namespace, pod)
	networkArgs := []string{"network", "ls", "-q", "-f", fmt.Sprintf("Name=%s", networkName)}
	exec := cmdz.Cmd(t.binary, networkArgs...).ErrorOnFailure(true)
	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t Translator) createPodNetwork(namespace string, pod k8sv1.Pod) (cmdz.Executer, error) {
	networkName := networkName(namespace, pod)
	networkArgs := []string{"network", "create", networkName}
	exec := cmdz.Cmd(t.binary, networkArgs...).ErrorOnFailure(true)
	return exec, nil
}

func (t Translator) podRootContainerId(namespace string, pod k8sv1.Pod) (string, error) {
	ctName := podRootContainerName(namespace, pod)
	exec := cmdz.Cmd(t.binary, "ps", "-q", "-f", fmt.Sprintf("Name=%s", ctName)).ErrorOnFailure(true)

	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t Translator) createPodRootContainer(namespace string, pod k8sv1.Pod) (cmdz.Executer, error) {
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

	exec := cmdz.Cmd(t.binary, runArgs...).ErrorOnFailure(true)
	return exec, nil
}

func (t Translator) volumeId(namespace string, vol k8sv1.Volume) (string, error) {
	volName := forgeResName(namespace, vol)
	volumeArgs := []string{"volume", "ls", "-q", "-f", fmt.Sprintf("Name=%s", volName)}
	exec := cmdz.Cmd(t.binary, volumeArgs...).ErrorOnFailure(true)

	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t Translator) createVolume(namespace string, vol k8sv1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.HostPath != nil {
		return t.createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return t.createEmptyDirPodVolume(namespace, vol)
	}
	err := fmt.Errorf("Not supported volume type for volume: %s !", vol.Name)
	return nil, err
}

func (t Translator) createHostPathPodVolume(namespace string, vol k8sv1.Volume) (cmdz.Executer, error) {
	hostPathType := *vol.VolumeSource.HostPath.Type
	if hostPathType != k8sv1.HostPathUnset {
		err := fmt.Errorf("Not supported HostPathType: %s for volume: %s !", hostPathType, vol.Name)
		return nil, err
	}

	name := forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPath.Path
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local")
	exec.AddArgs("-o", "o=bind", "-o", "type=none", "-o", "device="+path)
	exec.AddArgs(name)
	return exec, nil
}

func (t Translator) createEmptyDirPodVolume(namespace string, vol k8sv1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.EmptyDir == nil {
		err := fmt.Errorf("Bad EmptyDirVolume !")
		return nil, err
	}
	name := forgeResName(namespace, vol)
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local", name)
	return exec, nil
}

func (t Translator) podContainerId(namespace string, pod k8sv1.Pod, container k8sv1.Container) (string, error) {
	ctName := podContainerName(namespace, pod, container)
	psArgs := []string{"ps", "-q", "-f", fmt.Sprintf("Name=%s", ctName)}
	exec := cmdz.Cmd(t.binary, psArgs...).ErrorOnFailure(true)

	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t Translator) createPodContainer(namespace string, pod k8sv1.Pod, container k8sv1.Container, init bool) (cmdz.Executer, error) {
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
		err := fmt.Errorf("No supported restart policy: %s in container: %s !", restartPolicy, ctName)
		return nil, err
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
		err := fmt.Errorf("No supported pull policy: %s in container: %s !", pullPolicy, ctName)
		return nil, err
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
		for _, arg := range entrypoint[1:] {
			cmdArgs = append(cmdArgs, arg)
		}
	}

	if len(args) > 0 {
		for _, arg := range args {
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

	exec := cmdz.Cmd(t.binary, "run")
	//cmd.Retries = 2
	exec.AddArgs(runArgs...)
	exec.AddArgs(resourcesArgs...)
	exec.AddArgs(envArgs...)
	exec.AddArgs(labelArgs...)
	exec.AddArgs(image)
	exec.AddArgs(cmdArgs...)
	//cmds = append(cmds, cmd)

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

	return exec, nil
}
