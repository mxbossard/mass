package driver

import (
	"fmt"
	"testing"

	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedBinary0    = "docker0"
	expectedBinary     = "docker"
	expectedNamespace1 = "ns1"
	expectedNamespace2 = "ns2"

	volume1Name             = "vol1"
	expectedVolume1EmptyDir = k8sv1.EmptyDirVolumeSource{}
	volume1                 = k8sv1.Volume{
		Name:         volume1Name,
		VolumeSource: k8sv1.VolumeSource{EmptyDir: &expectedVolume1EmptyDir},
	}

	volume2Name             = "vol2"
	volume2Type             = k8sv1.HostPathUnset
	volume2Path             = "/tmp/foo"
	expectedVolume2HostPath = k8sv1.HostPathVolumeSource{
		Path: volume2Path,
		Type: &volume2Type,
	}
	volume2 = k8sv1.Volume{
		Name:         volume2Name,
		VolumeSource: k8sv1.VolumeSource{HostPath: &expectedVolume2HostPath},
	}

	container1Name  = "ct1"
	container1Image = "busybox:1.36.1-musl"
	container1      = forgeContainer(container1Name, container1Image, k8sv1.PullAlways, false, false, false, []string{"/bin/sleep"}, []string{"0.1"})
	container2Name  = "ct2"
	container2Image = "alpine:3.18.2"
	container2      = forgeContainer(container2Name, container2Image, k8sv1.PullIfNotPresent, true, false, false, []string{"/bin/sleep"}, []string{"10"})
	container3Name  = "ct3"
	container3Image = "busybox:1.35.0-musl"
	container3      = forgeContainer(container3Name, container3Image, k8sv1.PullAlways, false, true, false, []string{"/bin/sleep", "0.1"}, []string{})
	container4Name  = "ct4"
	container4Image = "alpine:3.17.4"
	container4      = forgeContainer(container4Name, container4Image, k8sv1.PullIfNotPresent, false, false, true, []string{"/bin/sleep", "10"}, []string{})
	container5Name  = "ct5"
	container5Image = "busybox:1.34.1-musl"
	container5      = forgeContainer(container5Name, container5Image, k8sv1.PullNever, true, true, false, []string{"/bin/sh"}, []string{"-c", "sleep 0.1"})
	container6Name  = "ct6"
	container6Image = "busybox:1.34.1-musl"
	container6      = forgeContainer(container6Name, container6Image, k8sv1.PullIfNotPresent, true, true, false, []string{"/bin/sh", "-c"}, []string{"/bin/sleep 0.1"})

	pod1Name                   = "pod1"
	expectedPod1Volumes        = []k8sv1.Volume{volume1, volume2}
	expectedPod1InitContainers = []k8sv1.Container{container1, container2, container6}
	expectedPod1Containers     = []k8sv1.Container{container3, container4}
	expectedPod1RestartPolicy  = k8sv1.RestartPolicyAlways
	pod1                       = forgePod(pod1Name, expectedPod1Volumes, expectedPod1InitContainers, expectedPod1Containers, expectedPod1RestartPolicy, false)

	pod2Name                   = "pod2"
	expectedPod2Volumes        = []k8sv1.Volume{volume1}
	expectedPod2InitContainers = []k8sv1.Container{}
	expectedPod2Containers     = []k8sv1.Container{container3}
	expectedPod2RestartPolicy  = k8sv1.RestartPolicyOnFailure
	pod2                       = forgePod(pod2Name, expectedPod2Volumes, expectedPod2InitContainers, expectedPod2Containers, expectedPod2RestartPolicy, true)

	pod3Name                   = "pod3"
	expectedPod3Volumes        = []k8sv1.Volume{}
	expectedPod3InitContainers = []k8sv1.Container{container1, container5}
	expectedPod3Containers     = []k8sv1.Container{container3, container4}
	expectedPod3RestartPolicy  = k8sv1.RestartPolicyNever
	pod3                       = forgePod(pod3Name, expectedPod3Volumes, expectedPod3InitContainers, expectedPod3Containers, expectedPod3RestartPolicy, false)
)

func forgeContainer(name, image string, pp k8sv1.PullPolicy, privileged, roRootFs, tty bool, cmd, args []string) k8sv1.Container {
	containerName := name
	expectedContainerImage := image
	expectedContainerCommand := cmd
	expectedContainerArgs := args
	expectedContainerWorkingDir := "/tmp/" + name + "_workdir1"
	expectedContainerPorts := []k8sv1.ContainerPort{
		k8sv1.ContainerPort{Name: name + "-port1", HostPort: 8080, ContainerPort: 80, HostIP: "1.2.3.4"},
	}
	expectedContainerEnv := []k8sv1.EnvVar{
		k8sv1.EnvVar{Name: name + "envKey1", Value: name + "envVal1"},
		k8sv1.EnvVar{Name: name + "envKey2", Value: name + "envVal2"},
		k8sv1.EnvVar{Name: name + "envKey3", Value: name + "envVal3"},
	}
	expectedContainerVolumeMounts := []k8sv1.VolumeMount{
		k8sv1.VolumeMount{Name: name + "-vol1", ReadOnly: true, MountPath: "/foo/bar"},
		k8sv1.VolumeMount{Name: name + "-vol2", ReadOnly: true, MountPath: "/foo/baz"},
	}
	var uid int64 = 1001
	var gid int64 = 1001
	nonRoot := true
	expectedContainerSecurityContext := k8sv1.SecurityContext{
		Privileged:             &privileged,
		RunAsUser:              &uid,
		RunAsGroup:             &gid,
		RunAsNonRoot:           &nonRoot,
		ReadOnlyRootFilesystem: &roRootFs,
	}

	container := k8sv1.Container{
		Name:            containerName,
		Image:           expectedContainerImage,
		Command:         expectedContainerCommand,
		Args:            expectedContainerArgs,
		WorkingDir:      expectedContainerWorkingDir,
		Ports:           expectedContainerPorts,
		Env:             expectedContainerEnv,
		VolumeMounts:    expectedContainerVolumeMounts,
		ImagePullPolicy: pp,
		SecurityContext: &expectedContainerSecurityContext,
		TTY:             tty,
	}
	return container
}

func forgePod(name string, volumes []k8sv1.Volume, initContainers, containers []k8sv1.Container, rp k8sv1.RestartPolicy, hostNetwork bool) k8sv1.Pod {
	podName := name
	expectedPod1Labels := map[string]string{
		name + "_labelKey1": name + "_labelVal1",
		name + "_labelKey1": name + "_labelVal2",
	}
	expectedPod1Annotations := map[string]string{
		name + "_annotationKey1":  name + "_annotationVal1",
		name + "_annotationKkey2": name + "_annotationVal2",
	}

	expectedPodHostname := name + "-hostname"
	expectedPodSubdomain := name + "-subdomain"
	expectedPodHostNetwork := hostNetwork
	var uid int64 = 1001
	var gid int64 = 1001
	nonRoot := true
	expectedPodSecurityContext := k8sv1.PodSecurityContext{
		RunAsUser:          &uid,
		RunAsGroup:         &gid,
		RunAsNonRoot:       &nonRoot,
		SupplementalGroups: []int64{1, 2, 3},
	}

	expectedPodHostAliases := []k8sv1.HostAlias{
		k8sv1.HostAlias{IP: "1.1.1.1", Hostnames: []string{name + "-foo", name + "-bar"}},
		k8sv1.HostAlias{IP: "2.2.2.2", Hostnames: []string{name + "-baz"}},
	}
	expectedPodDnsConfig := k8sv1.PodDNSConfig{
		Nameservers: []string{"8.8.8.8", "8.8.9.9"},
		Searches:    []string{name + ".foo.local"},
		Options:     []k8sv1.PodDNSConfigOption{},
	}
	pod := k8sv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:        podName,
			Labels:      expectedPod1Labels,
			Annotations: expectedPod1Annotations,
		},
		Spec: k8sv1.PodSpec{
			Volumes:         volumes,
			InitContainers:  initContainers,
			Containers:      containers,
			RestartPolicy:   rp,
			HostNetwork:     expectedPodHostNetwork,
			SecurityContext: &expectedPodSecurityContext,
			Hostname:        expectedPodHostname,
			Subdomain:       expectedPodSubdomain,
			HostAliases:     expectedPodHostAliases,
			DNSConfig:       &expectedPodDnsConfig,
		},
	}

	return pod
}

func TestForgeResName(t *testing.T) {
	prefix := "pre-foo"
	expectedVolume1Name := prefix + "_" + volume1Name
	volumeName := forgeResName(prefix, volume1)
	assert.Equal(t, expectedVolume1Name, volumeName, "Bad Volume name !")

	expectedPod1Name := prefix + "_" + pod1Name
	podName := forgeResName(prefix, pod1)
	assert.Equal(t, expectedPod1Name, podName, "Bad Pod name !")
}

func TestCreateVolume(t *testing.T) {
	translator := Translator{expectedBinary0}

	cmds1, err := translator.createVolume(expectedNamespace1, volume1)
	require.NoError(t, err, "should not error")
	expectedCmd1 := fmt.Sprintf("%s volume create --driver local %s_%s", expectedBinary0, expectedNamespace1, volume1Name)
	assert.Equal(t, expectedCmd1, cmds1.String())

	cmds2, err := translator.createVolume(expectedNamespace1, volume2)
	require.NoError(t, err, "should not error")
	expectedCmd2 := fmt.Sprintf("%s volume create --driver local -o o=bind -o type=none -o device=%s %s_%s", expectedBinary0, volume2Path, expectedNamespace1, volume2Name)
	assert.Equal(t, expectedCmd2, cmds2.String())
}

func TestCreatePodContainer(t *testing.T) {
	translator := Translator{expectedBinary0}

	ct1 := pod1.Spec.InitContainers[0]
	cmds1, err := translator.createPodContainer(expectedNamespace1, pod1, ct1, true)
	require.NoError(t, err, "should not error")
	expectedCmd1 := fmt.Sprintf(`%[1]s run --rm --name %[2]s_%s_%s --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=always -v %[2]s_%[4]s-vol1:/foo/bar:ro -v %[2]s_%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 0.1`, expectedBinary0, expectedNamespace1, pod1.Name, ct1.Name, ct1.Image)
	assert.Equal(t, expectedCmd1, cmds1.String())

	ct2 := pod1.Spec.InitContainers[1]
	cmds2, err := translator.createPodContainer(expectedNamespace1, pod1, ct2, true)
	require.NoError(t, err, "should not error")
	expectedCmd2 := fmt.Sprintf(`%[1]s run --rm --name %[2]s_%s_%s --privileged --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=missing -v %[2]s_%[4]s-vol1:/foo/bar:ro -v %[2]s_%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 10`, expectedBinary0, expectedNamespace1, pod1.Name, ct2.Name, ct2.Image)
	assert.Equal(t, expectedCmd2, cmds2.String())

	ct5 := pod3.Spec.InitContainers[1]
	cmds3, err := translator.createPodContainer(expectedNamespace2, pod3, ct5, true)
	require.NoError(t, err, "should not error")
	expectedCmd3 := fmt.Sprintf(`%s run --rm --name %s_%s_%s --privileged --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=never -v %[2]s_%[4]s-vol1:/foo/bar:ro -v %[2]s_%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sh -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s -c sleep 0.1`, expectedBinary0, expectedNamespace2, pod3.Name, ct5.Name, ct5.Image)
	assert.Equal(t, expectedCmd3, cmds3.String())
}

func TestCreatePodNetwork(t *testing.T) {
	translator := Translator{expectedBinary0}

	cmds1, err := translator.createPodNetwork(expectedNamespace1, pod1)
	require.NoError(t, err, "should not error")
	expectedCmd10 := fmt.Sprintf(`%s network create %s_%s_net`, expectedBinary0, expectedNamespace1, pod1.Name)
	assert.Equal(t, expectedCmd10, cmds1.String())

	cmds2, err := translator.createPodNetwork(expectedNamespace1, pod2)
	require.NoError(t, err, "should not error")
	expectedCmd20 := fmt.Sprintf(`%s network create %s_%s_net`, expectedBinary0, expectedNamespace1, pod2.Name)
	assert.Equal(t, expectedCmd20, cmds2.String())
}

func TestCreatePodRootContainer(t *testing.T) {
	translator := Translator{expectedBinary0}

	cmds1, err := translator.createPodRootContainer(expectedNamespace1, pod1)
	require.NoError(t, err, "should not error")
	expectedCmd11 := fmt.Sprintf(`%s run -d --name %[2]s_%[3]s_root --restart=always --network %[2]s_%[3]s_net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod1.Name)
	assert.Equal(t, expectedCmd11, cmds1.String())

	cmds2, err := translator.createPodRootContainer(expectedNamespace1, pod2)
	require.NoError(t, err, "should not error")
	expectedCmd21 := fmt.Sprintf(`%s run -d --name %[2]s_%[3]s_root --restart=always --network %[2]s_%[3]s_net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod2.Name)
	assert.Equal(t, expectedCmd21, cmds2.String())
}

func TestApplyPod(t *testing.T) {
	translator := Translator{expectedBinary}
	executor := Executor{translator: translator}
	err := executor.applyPod(expectedNamespace1, pod1)
	require.NoError(t, err)

}
