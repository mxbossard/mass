package pod

import (
	"fmt"
	"testing"

	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedBinary     = "docker0"
	expectedNamespace1 = "namespace1"

	volume1Name             = "vol1"
	expectedVolume1EmptyDir = k8sv1.EmptyDirVolumeSource{}
	volume1                 = k8sv1.Volume{
		Name:         volume1Name,
		VolumeSource: k8sv1.VolumeSource{EmptyDir: &expectedVolume1EmptyDir},
	}

	volume2Name             = "vol2"
	volume2Type             = k8sv1.HostPathUnset
	expectedVolume2HostPath = k8sv1.HostPathVolumeSource{
		Path: "/tmp/foo",
		Type: &volume2Type,
	}
	volume2 = k8sv1.Volume{
		Name:         volume2Name,
		VolumeSource: k8sv1.VolumeSource{HostPath: &expectedVolume2HostPath},
	}

	container1Name = "ct1"
	container1     = forgeContainer(container1Name, k8sv1.PullAlways, false, false, false)
	container2Name = "ct2"
	container2     = forgeContainer(container2Name, k8sv1.PullNever, true, false, false)
	container3Name = "ct3"
	container3     = forgeContainer(container3Name, k8sv1.PullIfNotPresent, false, true, false)
	container4Name = "ct4"
	container4     = forgeContainer(container4Name, k8sv1.PullAlways, false, false, true)
	container5Name = "ct5"
	container5     = forgeContainer(container5Name, k8sv1.PullAlways, true, true, true)

	pod1Name                   = "pod1"
	expectedPod1Volumes        = []k8sv1.Volume{volume1, volume2}
	expectedPod1InitContainers = []k8sv1.Container{container1, container2}
	expectedPod1Containers     = []k8sv1.Container{container3, container4, container5}
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
	expectedPod3InitContainers = []k8sv1.Container{container1}
	expectedPod3Containers     = []k8sv1.Container{container3, container4}
	expectedPod3RestartPolicy  = k8sv1.RestartPolicyNever
	pod3                       = forgePod(pod3Name, expectedPod3Volumes, expectedPod3InitContainers, expectedPod3Containers, expectedPod3RestartPolicy, false)
)

func forgeContainer(name string, pp k8sv1.PullPolicy, privileged, roRootFs, tty bool) k8sv1.Container {
	containerName := name
	expectedContainer1Image := name + "-image"
	expectedContainer1Command := []string{name + "_cmd1", name + "_cmd2"}
	expectedContainer1Args := []string{name + "_arg1", name + "_arg2", name + "_arg3"}
	expectedContainer1WorkingDir := name + "_workdir1"
	expectedContainer1Ports := []k8sv1.ContainerPort{
		k8sv1.ContainerPort{Name: name + "-port1", HostPort: 8080, ContainerPort: 80, HostIP: "1.2.3.4"},
	}
	expectedContainer1Env := []k8sv1.EnvVar{
		k8sv1.EnvVar{Name: name + "envKey1", Value: name + "envVal1"},
		k8sv1.EnvVar{Name: name + "envKey2", Value: name + "envVal2"},
		k8sv1.EnvVar{Name: name + "envKey3", Value: name + "envVal3"},
	}
	expectedContainer1VolumeMounts := []k8sv1.VolumeMount{
		k8sv1.VolumeMount{Name: name + "-vol1", ReadOnly: true, MountPath: "/foo/bar"},
		k8sv1.VolumeMount{Name: name + "-vol2", ReadOnly: true, MountPath: "/foo/baz"},
	}
	var uid int64 = 1001
	var gid int64 = 1001
	nonRoot := true
	expectedContainer1SecurityContext := k8sv1.SecurityContext{
		Privileged:             &privileged,
		RunAsUser:              &uid,
		RunAsGroup:             &gid,
		RunAsNonRoot:           &nonRoot,
		ReadOnlyRootFilesystem: &roRootFs,
	}

	container := k8sv1.Container{
		Name:            containerName,
		Image:           expectedContainer1Image,
		Command:         expectedContainer1Command,
		Args:            expectedContainer1Args,
		WorkingDir:      expectedContainer1WorkingDir,
		Ports:           expectedContainer1Ports,
		Env:             expectedContainer1Env,
		VolumeMounts:    expectedContainer1VolumeMounts,
		ImagePullPolicy: pp,
		SecurityContext: &expectedContainer1SecurityContext,
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

	expectedPod1Hostname := name + "-hostname"
	expectedPod1Subdomain := name + "-subdomain"
	expectedPod1HostNetwork := hostNetwork
	var uid int64 = 1001
	var gid int64 = 1001
	nonRoot := true
	expectedPod1SecurityContext1 := k8sv1.PodSecurityContext{
		RunAsUser:          &uid,
		RunAsGroup:         &gid,
		RunAsNonRoot:       &nonRoot,
		SupplementalGroups: []int64{1, 2, 3},
	}

	expectedPod1HostAliases := []k8sv1.HostAlias{
		k8sv1.HostAlias{IP: "1.1.1.1", Hostnames: []string{name + "-foo", name + "-bar"}},
		k8sv1.HostAlias{IP: "2.2.2.2", Hostnames: []string{name + "-baz"}},
	}
	expectedPod1DnsConfig := k8sv1.PodDNSConfig{
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
			Volumes:         expectedPod1Volumes,
			InitContainers:  expectedPod1InitContainers,
			Containers:      expectedPod1Containers,
			RestartPolicy:   rp,
			HostNetwork:     expectedPod1HostNetwork,
			SecurityContext: &expectedPod1SecurityContext1,
			Hostname:        expectedPod1Hostname,
			Subdomain:       expectedPod1Subdomain,
			HostAliases:     expectedPod1HostAliases,
			DNSConfig:       &expectedPod1DnsConfig,
		},
	}

	return pod
}

func TestForgeResName(t *testing.T) {
	translator := Translator{expectedBinary}
	prefix := "pre-foo"
	expectedVolume1Name := prefix + "_" + volume1Name
	volumeName, err := translator.forgeResName(prefix, volume1)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedVolume1Name, volumeName, "Bad Volume name !")

	expectedPod1Name := prefix + "_" + pod1Name
	podName, err := translator.forgeResName(prefix, pod1)
	require.NoError(t, err, "should not error")
	assert.Equal(t, expectedPod1Name, podName, "Bad Pod name !")
}

func TestCreateVolume(t *testing.T) {
	translator := Translator{expectedBinary}
	cmds, err := translator.createVolume(expectedNamespace1, volume1)
	require.NoError(t, err, "should not error")
	assert.Len(t, cmds, 1)
	expectedCmd := fmt.Sprintf("%s volume create --driver local %s_%s", expectedBinary, expectedNamespace1, volume1Name)
	assert.Equal(t, expectedCmd, cmds[0])
}
