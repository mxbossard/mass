package driver

import (
	"fmt"
	"log"

	"strings"

	"mby.fr/utils/cmdz"
	//"mby.fr/utils/promise"

	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/stringz"

	corev1 "k8s.io/api/core/v1"
)

const (
	//ContainerName_NamespaceSeparator = "__"
	//ContainerName_NameSeparator      = "__"
	ContainerName_Separator   = "__"
	ContainerName_PodRootFlag = "--root"
)

// TODO: remonter tout ce qui concerne les pods dans executor ne garder que les concepts à la maille docker dans le translator :
// - Containers et Namespaces
// TODO comment gerer les init containers ?
// TODO comment heberger les pod phase et container status ?
// TODO retirer les commits ?

type Translator struct {
	binary string
}

func (t Translator) describeNamespace(name string) (corev1.Namespace, error) {
	ns := descriptor.BuildNamespace(name)
	descriptor.CompleteK8sResourceDefaults(&ns)
	return ns, nil
}

func (t Translator) listNamespaces() (namespaces []corev1.Namespace, err error) {
	allNsAllRootContainersFilter := podContainerNameFilter("", "", "", true)
	exec := cmdz.Cmd(t.binary, "ps", "-a", "--format", "{{ .Names }}", "-f", allNsAllRootContainersFilter).ErrorOnFailure(true)
	_, err = exec.BlockRun()
	if err != nil {
		return nil, err
	}
	stdOut := strings.TrimSpace(exec.StdoutRecord())
	podRootCtNames, _ := stringz.SplitByRegexp(stdOut, `\s`)
	namespaceNames := map[string]bool{}
	for _, podRootCtName := range podRootCtNames {
		namespace := getNamespaceNameFromContainerName(podRootCtName)
		namespaceNames[namespace] = true
	}
	for nsName, _ := range namespaceNames {
		ns, err := t.describeNamespace(nsName)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}

	return
}

func (t Translator) listPodContainers(namespace, podName string) (containers map[string]corev1.Container, err error) {
	allContainersFilter := podContainerNameFilter(namespace, podName, "", false)
	//log.Printf("docker ps filter: %s", allRootContainersFilter)
	exec := cmdz.Cmd(t.binary, "ps", "-a", "--format", "{{ .Names }} {{ .Image }}", "-f", allContainersFilter).ErrorOnFailure(true)
	_, err = exec.BlockRun()
	if err != nil {
		return nil, err
	}
	//log.Printf("listPods: RC=%d, stdout=%s", rc, exec.StdoutRecord())
	stdOut := strings.TrimSpace(exec.StdoutRecord())
	podCtInfos, _ := stringz.SplitByRegexp(stdOut, `\n`)
	containers = make(map[string]corev1.Container, len(podCtInfos))
	for _, ctInfos := range podCtInfos {
		splitted := strings.Split(ctInfos, " ")
		ctName := splitted[0]
		podCtName := getContainerNameFromContainerName(ctName)
		ctImage := splitted[1]
		ct := descriptor.BuildDefaultContainer(podCtName, ctImage)
		containers[ctName] = ct
	}
	//log.Printf("list %s/%s containers => %v", namespace, podName, containers)
	return
}

// TODO: remove listPods not needed. listPodContainers is enough
/*
func (t Translator) listPods(namespace string) (pods []corev1.Pod, err error) {
	allRootContainersFilter := podContainerNameFilter(namespace, "", "", true)
	//log.Printf("docker ps filter: %s", allRootContainersFilter)
	exec := cmdz.Cmd(t.binary, "ps", "-a", "--format", "{{ .Names }}", "-f", allRootContainersFilter).ErrorOnFailure(true)
	_, err = exec.BlockRun()
	if err != nil {
		return nil, err
	}
	//log.Printf("listPods: RC=%d, stdout=%s", rc, exec.StdoutRecord())
	stdOut := strings.TrimSpace(exec.StdoutRecord())
	podRootCtNames, _ := stringz.SplitByRegexp(stdOut, `\s`)
	podNames := []string{}
	//log.Printf("list of pods: %v", podRootCtNames)
	for _, podRootCtName := range podRootCtNames {
		podName := getPodNameFromContainerName(podRootCtName)
		podNames = append(podNames, podName)
	}
	for _, podName := range podNames {
		pod, err := t.describePod(namespace, podName)
		if err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}
	log.Printf("list pods in ns: %s => %v", namespace, pods)
	return
}
*/

func (t Translator) deletePod(namespace, name string) (cmdz.Executer, error) {
	allContainersFilter := podContainerNameFilter(namespace, name, "", false)
	log.Printf("allContainersFilter: %s", allContainersFilter)
	psExec := cmdz.Cmd(t.binary, "ps", "--format", "{{ .Names }}", "-f", allContainersFilter).ErrorOnFailure(true)
	_, err := psExec.BlockRun()
	if err != nil {
		return nil, err
	}

	stdOut := strings.TrimSpace(psExec.StdoutRecord())
	podCtNames, _ := stringz.SplitByRegexp(stdOut, `\s`)
	rootContainerName := podRootContainerName0(namespace, name)
	podCtNames = append(podCtNames, rootContainerName)
	var rmExec cmdz.Executer
	if len(podCtNames) > 0 {
		log.Printf("deletePod: Deleting containers: %v", podCtNames)
		rmArgs := []string{"rm", "-f"}
		rmArgs = append(rmArgs, podCtNames...)
		rmExec = cmdz.Cmd(t.binary, rmArgs...).ErrorOnFailure(true)
	} else {
		rmExec = cmdz.Cmd("true")
	}
	return rmExec, nil
}

/*
func (t Translator) commitPod(namespace string, pod corev1.Pod) (cmdz.Executer, error) {
	b, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", fmt.Sprintf("echo '%s' > /podConfig.txt", string(b))}
	exec := cmdz.Cmd(t.binary, execArgs...)
	return exec, nil
}

func (t Translator) getCommitedRootPod(rootPodName string) (*corev1.Pod, error) {
	execArgs := []string{"exec", rootPodName, "/bin/sh", "-c", "cat /podConfig.txt"}
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

func (t Translator) getCommitedPod(namespace string, pod corev1.Pod) (*corev1.Pod, error) {
	ctName := podRootContainerName(namespace, pod)
	return t.getCommitedRootPod(ctName)
}
*/

func (t Translator) commitPodPhase(namespace string, pod corev1.Pod, phase corev1.PodPhase) cmdz.Executer {
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", fmt.Sprintf("echo %s > /podPhase.txt", phase)}
	exec := cmdz.Cmd(t.binary, execArgs...)
	return exec
}

func (t Translator) getCommitedPodPhase(namespace string, pod corev1.Pod) (*corev1.PodPhase, error) {
	ctName := podRootContainerName(namespace, pod)
	execArgs := []string{"exec", ctName, "/bin/sh", "-c", "cat /podPhase.txt"}
	exec := cmdz.Cmd(t.binary, execArgs...).ErrorOnFailure(true)
	_, err := exec.BlockRun()
	if err != nil {
		return nil, err
	}
	p := corev1.PodPhase(exec.StdoutRecord())
	return &p, nil
}

func (t Translator) containerState(namespace string, pod corev1.Pod, container corev1.Container) (p corev1.ContainerState, e error) {
	return
}

func (t Translator) podNetworkId(namespace string, pod corev1.Pod) (string, error) {
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

func (t Translator) createPodNetwork(namespace string, pod corev1.Pod) (cmdz.Executer, error) {
	networkName := networkName(namespace, pod)
	networkArgs := []string{"network", "create", networkName}
	exec := cmdz.Cmd(t.binary, networkArgs...).ErrorOnFailure(true)
	return exec, nil
}

func (t Translator) podRootContainerId(namespace string, pod corev1.Pod) (string, error) {
	ctName := podRootContainerName(namespace, pod)
	exec := cmdz.Cmd(t.binary, "ps", "-q", "-f", fmt.Sprintf("Name=%s", ctName)).ErrorOnFailure(true)

	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t Translator) createPodRootContainer(namespace string, pod corev1.Pod) (cmdz.Executer, error) {
	ctName := podRootContainerName(namespace, pod)
	cpusArgs := "--cpus=0.05"
	memoryArgs := "--memory=64m"
	networkName := networkName(namespace, pod)
	addHostRules := ""
	pauseImage := "alpine:3.17.3"

	runArgs := []string{"run", "--detach", "--name", ctName, "--restart=always", "--network", networkName,
		cpusArgs, memoryArgs} //"--memory-swappiness=0"

	if addHostRules != "" {
		runArgs = append(runArgs, addHostRules)
	}

	runArgs = append(runArgs, pauseImage)
	runArgs = append(runArgs, "/bin/sleep", "inf")

	exec := cmdz.Cmd(t.binary, runArgs...).ErrorOnFailure(true)
	return exec, nil
}

func (t Translator) volumeId(namespace string, vol corev1.Volume) (string, error) {
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

func (t Translator) createVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.HostPath != nil {
		return t.createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return t.createEmptyDirPodVolume(namespace, vol)
	}
	err := fmt.Errorf("Not supported volume type for volume: %s !", vol.Name)
	return nil, err
}

func (t Translator) createHostPathPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	hostPathType := *vol.VolumeSource.HostPath.Type
	if hostPathType != corev1.HostPathUnset {
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

func (t Translator) createEmptyDirPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.EmptyDir == nil {
		err := fmt.Errorf("Bad EmptyDirVolume !")
		return nil, err
	}
	name := forgeResName(namespace, vol)
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local", name)
	return exec, nil
}

func (t Translator) podContainerId(namespace string, pod corev1.Pod, container corev1.Container) (string, error) {
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

func (t Translator) createPodContainer(namespace string, pod corev1.Pod, container corev1.Container, init bool) (cmdz.Executer, error) {
	ctName := podContainerName(namespace, pod, container)
	image := container.Image
	privileged := false
	if container.SecurityContext != nil {
		if container.SecurityContext.Privileged != nil {
			privileged = *container.SecurityContext.Privileged
		}
	}
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

	runArgs = append(runArgs, "--detach")

	if init {
		restartPolicy = corev1.RestartPolicyNever
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
	case corev1.PullAlways:
		dockerPullPolicy = "always"
	case corev1.PullNever:
		dockerPullPolicy = "never"
	case corev1.PullIfNotPresent:
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

func networkName(namespace string, pod corev1.Pod) string {
	podName := forgeResName(namespace, pod)
	networkName := fmt.Sprintf("%s--net", podName)
	return networkName
}

func podRootContainerName0(namespace string, podName string) string {
	ctName := fmt.Sprintf("%s%s%s%s", namespace, ContainerName_Separator, podName, ContainerName_PodRootFlag)
	return ctName
}

func podRootContainerName(namespace string, pod corev1.Pod) string {
	return podRootContainerName0(namespace, pod.ObjectMeta.Name)
}

func podContainerName(namespace string, pod corev1.Pod, container corev1.Container) string {
	podName := forgeResName(namespace, pod)
	ctName := forgeResName(podName, container)
	return ctName
}

func getNamespaceNameFromContainerName(name string) string {
	splitted := strings.Split(name, ContainerName_Separator)
	return splitted[0]
}

func getPodNameFromContainerName(name string) string {
	splitted := strings.Split(name, ContainerName_Separator)
	podName := splitted[1]
	return strings.Split(podName, "--")[0]
}

func getContainerNameFromContainerName(name string) string {
	splitted := strings.Split(name, ContainerName_Separator)
	ctName := splitted[2]
	return ctName
}

// Filter for docker ps -f arg
func podContainerNameFilter(namespace, podName, containerName string, isRoot bool) string {
	sb := strings.Builder{}
	sb.WriteString("name=^")
	if namespace != "" {
		sb.WriteString(namespace)
	} else {
		sb.WriteString(".+")
	}
	sb.WriteString(ContainerName_Separator)
	if podName != "" {
		sb.WriteString(podName)
	} else {
		sb.WriteString(".+")
	}

	if isRoot {
		sb.WriteString(ContainerName_PodRootFlag)
	} else {
		sb.WriteString(ContainerName_Separator)
		if containerName != "" {
			sb.WriteString(containerName)
		} else {
			sb.WriteString(".+")
		}
	}

	sb.WriteString("$")
	return sb.String()
}

func forgeResName(prefix string, resource any) (name string) {
	var resName string
	switch res := resource.(type) {
	case corev1.Pod:
		resName = res.ObjectMeta.Name
	case corev1.Volume:
		resName = res.Name
	case corev1.VolumeMount:
		resName = res.Name
	case corev1.Container:
		resName = res.Name
	default:
		log.Fatalf("Cannot forge a name for unknown type: %T !", resource)
	}

	name = fmt.Sprintf("%s%s%s", prefix, ContainerName_Separator, resName)
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
