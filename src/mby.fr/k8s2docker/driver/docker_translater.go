package driver

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/cmdz"
	"mby.fr/utils/collections"

	//"mby.fr/utils/promise"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	//k8syaml "sigs.k8s.io/yaml"
	"encoding/json"
)

/*
const (
	//ContainerName_NamespaceSeparator = "__"
	//ContainerName_NameSeparator      = "__"
	ContainerName_Separator   = "__"
	ContainerName_PodRootFlag = "--root"
)
*/

const (
	MetadataLabelKeyPrefix = "k8s2docker.mby.fr"
	ContainerMetadataKey   = "descriptor.container"
	PodMetadataKey         = "descriptor.pod"
)

func forgePodVolumeName(namespace, podName, volName string) (name string) {
	return namespace + ContainerName_Separator + podName + ContainerName_Separator + volName
}

func forgePodContainerName(namespace, podName, ctName string) (name string) {
	return namespace + ContainerName_Separator + podName + ContainerName_Separator + ctName
}

func forgePodRootContainerName(namespace string, podName string) string {
	ctName := fmt.Sprintf("%s%s%s%s", namespace, ContainerName_Separator, podName, ContainerName_PodRootFlag)
	return ctName
}

// TODO: remonter tout ce qui concerne les pods dans executor ne garder que les concepts à la maille docker dans le translator :
// - Containers et Namespaces
// TODO comment gerer les init containers ?
// TODO comment heberger les pod phase et container status ?
// TODO retirer les commits ?

type DockerTranslater struct {
	binary string
}

func (t DockerTranslater) CreateNamespace(ns corev1.Namespace) cmdz.Executer {
	// Nothing to do
	return cmdz.Cmd("true")
}

func (t DockerTranslater) UpdateNamespace(ns corev1.Namespace) cmdz.Executer {
	// Nothing to do
	return cmdz.Cmd("true")
}

func (t DockerTranslater) DeleteNamsepace(ns string) cmdz.Executer {
	// Nothing to do
	return cmdz.Cmd("true")
}

func (t DockerTranslater) ListNamespaceNames() cmdz.Formatter[[]string, error] {
	allNsAllRootContainersFilter := podContainerNameFilter("", "", "", true)
	formatter := func(rc int, stdout, stderr []byte, inErr error) (res []string, err error) {
		if inErr != nil {
			return nil, inErr
		}
		res = strings.Split(string(stdout), "\n")
		res = collections.Map(res, func(in string) string {
			return strings.Split(in, ContainerName_Separator)[0]
		})
		res = collections.Filter(res, func(in string) bool {
			return in != ""
		})
		return res, nil
	}
	return cmdz.FormattedCmd[[]string, error](formatter, t.binary, "ps", "-a", "--format", "{{ .Names }}", "-f", allNsAllRootContainersFilter).ErrorOnFailure(true)
}

func (t DockerTranslater) SetupPod(namespace string, pod corev1.Pod) cmdz.Executer {
	// Do we need to create the network ?
	/*
		networkName := networkName(namespace, pod)
		networkArgs := []string{"network", "create", networkName}
		exec := cmdz.Cmd(t.binary, networkArgs...).ErrorOnFailure(true)
	*/

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
	return exec
}

func (t DockerTranslater) DeletePod(namespace, name string) cmdz.Executer {
	allContainersFilter := podContainerNameFilter(namespace, name, "", false)
	return cmdz.Sh(t.binary, "rm -f $(", t.binary, "ps -q -f", allContainersFilter, ")")
	/*
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
		return rmExec
	*/
}

func (t DockerTranslater) DescribePod(namepsace, podName string) cmdz.Formatter[corev1.Pod, error] {
	// TODO
	return nil
}

func (t DockerTranslater) CreateVolume(namespace, podName string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.HostPath != nil {
		return t.createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return t.createEmptyDirPodVolume(namespace, vol)
	}
	return nil, fmt.Errorf("not supported volume type for volume: %s", vol.Name)
}

func (t DockerTranslater) DeleteVolume(namespace, podName, name string) cmdz.Executer {
	volumeFilter := podContainerNameFilter(namespace, podName, name, false)
	return cmdz.Sh(t.binary, "volume rm -f $(", t.binary, "volume ls -q -f", volumeFilter, ")")
}

func (t DockerTranslater) DescribeVolume(namespace, podName, name string) cmdz.Formatter[[]corev1.Volume, error] {
	volName := forgePodVolumeName(namespace, podName, name)
	formatter := func(rc int, stdout, stderr []byte, inErr error) (res []corev1.Volume, err error) {
		if inErr != nil {
			return nil, inErr
		}
		var jsonResults []map[string]any
		err = yaml.Unmarshal(stdout, &jsonResults)
		// FIXME: bind mount are not listed in volumes !!!
		for _, jr := range jsonResults {
			// FIXME: manage cast errors
			name := jr["Name"].(string)
			vol := descriptor.BuildDefaultEmptyDirVolume(name)
			res = append(res, vol)
		}
		return
	}
	return cmdz.FormattedCmd[[]corev1.Volume, error](formatter, t.binary, "inspect", "volume", volName).ErrorOnFailure(true)
}

func (t DockerTranslater) ListVolumeNames(namespace, podName string) cmdz.Formatter[[]string, error] {
	allContainersFilter := podContainerNameFilter(namespace, podName, "", false)
	formatter := func(rc int, stdout, stderr []byte, inErr error) (res []string, err error) {
		if inErr != nil {
			return nil, inErr
		}
		res = strings.Split(string(stdout), "\n")
		res = collections.Map(res, func(in string) string {
			s := strings.Split(in, ContainerName_Separator)
			if len(s) == 3 {
				return s[2]
			}
			return ""
		})
		res = collections.Filter(res, func(in string) bool {
			return in != ""
		})
		return
	}
	return cmdz.FormattedCmd[[]string, error](formatter, t.binary, "volume", "ls", "--format", "{{ .Names }}", "-f", allContainersFilter).ErrorOnFailure(true)

}

func (t DockerTranslater) CreatePodContainer(namespace string, pod corev1.Pod, container corev1.Container, init bool) (cmdz.Executer, error) {
	ctName := podContainerName(namespace, pod, container)
	image := container.Image
	privileged := false
	readOnlyRootFs := false
	if container.SecurityContext != nil {
		if container.SecurityContext.Privileged != nil {
			privileged = *container.SecurityContext.Privileged
			readOnlyRootFs = *container.SecurityContext.ReadOnlyRootFilesystem
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

	if init {
		restartPolicy = corev1.RestartPolicyNever
	}

	var runArgs []string
	var resourcesArgs []string
	var envArgs []string
	var cmdArgs []string
	var labelArgs []string

	runArgs = append(runArgs, "--detach")
	runArgs = append(runArgs, "--rm")

	if privileged {
		runArgs = append(runArgs, "--privileged")
	}
	if readOnlyRootFs {
		runArgs = append(runArgs, "--read-only")
	}
	if tty {
		runArgs = append(runArgs, "-t")
	}

	runArgs = append(runArgs, "--name")
	runArgs = append(runArgs, ctName)

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
		err := fmt.Errorf("no supported restart policy: %s in container: %s", restartPolicy, ctName)
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
		err := fmt.Errorf("no supported pull policy: %s in container: %s", pullPolicy, ctName)
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

func (t DockerTranslater) UpdatePodContainer(namespace string, pod corev1.Pod, ct corev1.Container) (cmdz.Executer, error) {
	return nil, fmt.Errorf("DockerTranslater.CreatePodContainer() not implemented yet")
}

func (t DockerTranslater) DeletePodContainer(namespace, podName, name string) cmdz.Executer {
	ctName := forgePodContainerName(namespace, podName, name)
	return cmdz.Cmd(t.binary, "rm", "-f", ctName)
}

func (t DockerTranslater) DescribePodContainer(namespace, podName, name string) cmdz.Formatter[[]corev1.Container, error] {
	ctName := forgePodContainerName(namespace, podName, name)
	formatter := func(rc int, stdout, stderr []byte, inErr error) (res []corev1.Container, err error) {
		if inErr != nil {
			return nil, inErr
		}
		var jsonResults []map[string]any
		err = yaml.Unmarshal(stdout, &jsonResults)
		for _, jr := range jsonResults {
			// FIXME: manage errors of casts and path with eplorer lib

			config := jr["Config"].(map[string]any)
			hostConfig := jr["HostConfig"].(map[string]any)

			labels := config["Labels"].(map[string]any)
			originalContainer, err := loadLabelMetadata[corev1.Container](labels, ContainerMetadataKey)
			if err != nil {
				return nil, err
			}

			name := jr["Name"].(string)
			name = strings.Split(name, "/")[1]

			privileged := hostConfig["Privileged"].(bool)
			readOnlyRootFs := hostConfig["Privileged"].(bool)

			image := config["Image"].(string)
			workingDir := config["WorkingDir"].(string)
			tty := config["Tty"].(bool)
			user := config["User"].(string)
			userSplit := strings.Split(user, ":")
			var uid, gid *int64
			if len(userSplit) > 1 {
				u, err := strconv.ParseInt(userSplit[0], 10, 64)
				if err != nil {
					return nil, err
				}
				uid = &u
				g, err := strconv.ParseInt(userSplit[1], 10, 64)
				if err != nil {
					return nil, err
				}
				gid = &g
			} else if len(user) > 0 {
				u, err := strconv.ParseInt(user, 10, 64)
				if err != nil {
					return nil, err
				}
				uid = &u
				gid = nil
			}

			// FIXME: how to set securityContext.RunAsNonRoot ? a label ???
			// To get from labels :
			// - ImagePullPolicy
			// - resources & limits
			// - port names

			cmd := []any{}
			if config["Cmd"] != nil {
				cmd = config["Cmd"].([]any)
			}
			entrypoint := []any{}
			if config["Entrypoint"] != nil {
				entrypoint = config["Entrypoint"].([]any)
			}
			env := []any{}
			if config["Env"] != nil {
				env = config["Env"].([]any)
			}

			var securityContext corev1.SecurityContext
			originalSecurityContext := originalContainer.SecurityContext
			if originalSecurityContext != nil {
				securityContext = descriptor.BuildDefaultSecurityContext(uid, gid)
				securityContext.Privileged = &privileged
				securityContext.ReadOnlyRootFilesystem = &readOnlyRootFs
				// Read RunAsNonRoot from originalContainer
				securityContext.RunAsNonRoot = originalSecurityContext.RunAsNonRoot
			}

			// Read probes config from originalContainer
			originalStartupProbe := originalContainer.StartupProbe
			originalLivenessProbe := originalContainer.LivenessProbe
			originalReadinessProbe := originalContainer.ReadinessProbe

			// Read ports config from originalContainer
			originalPorts := originalContainer.Ports

			var containerPorts []corev1.ContainerPort
			if hostConfig["PortBindings"] != nil {
				portBindings := hostConfig["PortBindings"].(map[string]any)
				for key, value := range portBindings {
					bindings := value.([]any)
					split := strings.Split(key, "/")
					containerPort := split[0]
					containerPortNum, _ := strconv.ParseInt(containerPort, 10, 32)
					var protocol corev1.Protocol
					switch split[1] {
					case "tcp":
						protocol = corev1.ProtocolTCP
					case "udp":
						protocol = corev1.ProtocolUDP
					default:
						err = fmt.Errorf("not supported protocol: %s", split[1])
						return nil, err
					}

					matchingOriginalPorts := collections.Filter[corev1.ContainerPort](originalPorts, func(cp corev1.ContainerPort) bool {
						return cp.ContainerPort == int32(containerPortNum)
					})
					portName := "nonamed"
					if len(matchingOriginalPorts) == 0 {
						// New port absent from original config !
						portName = key

					} else if len(matchingOriginalPorts) > 1 {
						// Multiple ports with same container port !
						err = fmt.Errorf("problem with original ports config ! Duplicate container port config found for port: %d", containerPortNum)
						return nil, err
					} else {
						portName = matchingOriginalPorts[0].Name
					}

					for _, binding := range bindings {
						b := binding.(map[string]any)
						hostPort := b["HostPort"].(string)
						hostPortNum, _ := strconv.ParseInt(hostPort, 10, 32)
						hostIp := b["HostIp"].(string)
						ctPort := corev1.ContainerPort{Name: portName, Protocol: protocol, ContainerPort: int32(containerPortNum), HostPort: int32(hostPortNum), HostIP: hostIp}
						containerPorts = append(containerPorts, ctPort)
					}
				}
			}

			// Read volumeMounts config from originalContainer
			originalVolumeMounts := originalContainer.VolumeMounts

			var volumeMounts []corev1.VolumeMount
			if jr["Mounts"] != nil {
				mounts := jr["Mounts"].([]any)
				for _, mount := range mounts {
					m := mount.(map[string]any)
					path := m["Destination"].(string)
					readOnly := !m["RW"].(bool)
					name := "nonamed"
					if m["Name"] != nil {
						name = m["Name"].(string)
					} else {
						matchingMounts := collections.Filter[corev1.VolumeMount](originalVolumeMounts, func(vm corev1.VolumeMount) bool {
							return vm.MountPath == path
						})
						if len(matchingMounts) == 0 {
							// New mount absent from original config !
							name = path
						} else if len(matchingMounts) > 1 {
							// Multiple mounts with same container port !
							err = fmt.Errorf("problem with original volumeMounts config ! Duplicate volumeMount config found for mount path: %s", path)
							return nil, err
						} else {
							name = matchingMounts[0].Name
						}
					}

					volMount := corev1.VolumeMount{Name: name, MountPath: path, ReadOnly: readOnly}
					volumeMounts = append(volumeMounts, volMount)
				}
			}

			nanoCpus := hostConfig["NanoCpus"].(int)
			memory := hostConfig["Memory"].(int)
			resources := corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewScaledQuantity(int64(nanoCpus), resource.Nano),
					corev1.ResourceMemory: *resource.NewQuantity(int64(memory), resource.BinarySI),
				},
			}

			ct := descriptor.BuildDefaultContainer(name, image)
			ct.WorkingDir = workingDir
			ct.SecurityContext = &securityContext
			for _, e := range entrypoint {
				ct.Command = append(ct.Command, e.(string))
			}
			for _, c := range cmd {
				ct.Args = append(ct.Args, c.(string))
			}
			ct.TTY = tty
			for _, kv := range env {
				s := strings.Split(kv.(string), "=")
				ct.Env = append(ct.Env, corev1.EnvVar{Name: s[0], Value: s[1]})
			}
			ct.VolumeMounts = append(ct.VolumeMounts, volumeMounts...)
			ct.StartupProbe = originalStartupProbe
			ct.LivenessProbe = originalLivenessProbe
			ct.ReadinessProbe = originalReadinessProbe
			ct.Ports = append(ct.Ports, containerPorts...)

			ct.Resources = resources

			res = append(res, ct)
		}
		return
	}
	return cmdz.FormattedCmd[[]corev1.Container, error](formatter, t.binary, "inspect", ctName).ErrorOnFailure(true)
}

func (t DockerTranslater) ListPodContainerNames(namespace, podName string) cmdz.Formatter[[]string, error] {
	allContainersFilter := podContainerNameFilter(namespace, podName, "", false)
	formatter := func(rc int, stdout, stderr []byte, inErr error) (res []string, err error) {
		if inErr != nil {
			return nil, inErr
		}
		res = strings.Split(string(stdout), "\n")
		res = collections.Map(res, func(in string) string {
			s := strings.Split(in, ContainerName_Separator)
			if len(s) == 3 {
				return s[2]
			}
			return ""
		})
		res = collections.Filter(res, func(in string) bool {
			return in != ""
		})
		return
	}
	return cmdz.FormattedCmd[[]string, error](formatter, t.binary, "ps", "-a", "--format", "{{ .Names }}", "-f", allContainersFilter).ErrorOnFailure(true)
}

func (t DockerTranslater) createHostPathPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	hostPathType := *vol.VolumeSource.HostPath.Type
	if hostPathType != corev1.HostPathUnset {
		err := fmt.Errorf("not supported HostPathType: %s for volume: %s", hostPathType, vol.Name)
		return nil, err
	}

	name := forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPath.Path
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local")
	exec.AddArgs("-o", "o=bind", "-o", "type=none", "-o", "device="+path)
	exec.AddArgs(name)
	return exec, nil
}

func (t DockerTranslater) createEmptyDirPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.EmptyDir == nil {
		err := fmt.Errorf("bad EmptyDirVolume")
		return nil, err
	}
	name := forgeResName(namespace, vol)
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local", name)
	return exec, nil
}

func forgeLabelDataArgs(key, value string) []string {
	return []string{"--label", fmt.Sprintf("%s=%s", key, value)}
}

func forgeLabelMetadataKey(key string) string {
	return fmt.Sprintf("%s.%s", MetadataLabelKeyPrefix, key)
}

func forgeLabelMetadataArgs(key string, value any) ([]string, error) {
	key = forgeLabelMetadataKey(key)
	labelValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return forgeLabelDataArgs(key, string(labelValue)), nil
}

func appendLabelMetadataArgs(labels map[string]any, key string, value any) error {
	key = forgeLabelMetadataKey(key)
	labelValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	labels[key] = string(labelValue)
	return nil
}

func loadLabelMetadata[T any](labels map[string]any, key string) (value T, err error) {
	key = forgeLabelMetadataKey(key)
	labelVal, ok := labels[key]
	if !ok {
		err = fmt.Errorf("no metada label correspond to key: %s", key)
		return value, err
	}
	labelStr, ok := labelVal.(string)
	if !ok {
		err = fmt.Errorf("bad type for label value with key: %s", key)
		return value, err
	}

	err = json.Unmarshal([]byte(labelStr), &value)
	return
}

func boolPtr(in bool) *bool {
	return &in
}

func int64Ptr(in int64) *int64 {
	return &in
}

/*
func (t DockerTranslater) podNetworkId(namespace string, pod corev1.Pod) (string, error) {
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

func (t DockerTranslater) createPodNetwork(namespace string, pod corev1.Pod) (cmdz.Executer, error) {
	networkName := networkName(namespace, pod)
	networkArgs := []string{"network", "create", networkName}
	exec := cmdz.Cmd(t.binary, networkArgs...).ErrorOnFailure(true)
	return exec, nil
}

func (t DockerTranslater) podRootContainerId(namespace string, pod corev1.Pod) (string, error) {
	ctName := podRootContainerName(namespace, pod)
	exec := cmdz.Cmd(t.binary, "ps", "-q", "-f", fmt.Sprintf("Name=%s", ctName)).ErrorOnFailure(true)

	_, err := exec.BlockRun()
	if err != nil {
		return "", err
	}
	id := exec.StdoutRecord()
	return id, nil
}

func (t DockerTranslater) volumeId(namespace string, vol corev1.Volume) (string, error) {
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

func (t DockerTranslater) createVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.HostPath != nil {
		return t.createHostPathPodVolume(namespace, vol)
	} else if vol.VolumeSource.EmptyDir != nil {
		return t.createEmptyDirPodVolume(namespace, vol)
	}
	err := fmt.Errorf("not supported volume type for volume: %s", vol.Name)
	return nil, err
}

func (t DockerTranslater) createHostPathPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	hostPathType := *vol.VolumeSource.HostPath.Type
	if hostPathType != corev1.HostPathUnset {
		err := fmt.Errorf("not supported HostPathType: %s for volume: %s", hostPathType, vol.Name)
		return nil, err
	}

	name := forgeResName(namespace, vol)
	path := vol.VolumeSource.HostPath.Path
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local")
	exec.AddArgs("-o", "o=bind", "-o", "type=none", "-o", "device="+path)
	exec.AddArgs(name)
	return exec, nil
}

func (t DockerTranslater) createEmptyDirPodVolume(namespace string, vol corev1.Volume) (cmdz.Executer, error) {
	if vol.VolumeSource.EmptyDir == nil {
		err := fmt.Errorf("bad EmptyDirVolume")
		return nil, err
	}
	name := forgeResName(namespace, vol)
	exec := cmdz.Cmd(t.binary, "volume", "create", "--driver", "local", name)
	return exec, nil
}

func (t DockerTranslater) podContainerId(namespace string, pod corev1.Pod, container corev1.Container) (string, error) {
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

func (t DockerTranslater) createPodContainer(namespace string, pod corev1.Pod, container corev1.Container, init bool) (cmdz.Executer, error) {
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
		err := fmt.Errorf("no supported restart policy: %s in container: %s", restartPolicy, ctName)
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
		err := fmt.Errorf("no supported pull policy: %s in container: %s", pullPolicy, ctName)
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

	exec := cmdz.Cmd(t.binary, "run")
	//cmd.Retries = 2
	exec.AddArgs(runArgs...)
	exec.AddArgs(resourcesArgs...)
	exec.AddArgs(envArgs...)
	exec.AddArgs(labelArgs...)
	exec.AddArgs(image)
	exec.AddArgs(cmdArgs...)
	//cmds = append(cmds, cmd)

	//    # Test if ct already started or start it if excited or create it
	//    cmd="docker ps --format '{{ .Names }}' | grep -w '$ctName' || docker ps -f 'status=created' -f 'status=exited' --format '{{ .Names }}' | grep -w '$ctName' && docker start '$ctName' || docker run -d --name \"$ctName\" $resourcesArgs --network 'container:$podName' $runArgs $envArgs \"$image\" $entrypointArgs $cmdArgs"
	//    ! $DEBUG && >&2 echo "- $cmd"
	//    ! $DEBUG || >&2 echo "Running container $podName:$name ..."
	//    echo "$cmd"

	//    # FIXME: should write containers.txt once on pod creation
	//    #cmd="docker exec -u=0 '$podName' /bin/sh -c 'echo $ctName >> /containers.txt'"
	//    #! $DEBUG && >&2 echo "- $cmd"
	//    #echo "$cmd"

	return exec, nil
}
*/
