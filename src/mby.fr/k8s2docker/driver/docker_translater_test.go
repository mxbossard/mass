package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/cmdz"
)

//"mby.fr/utils/promise"

func Test_forgeLabelMetadataArgs(t *testing.T) {
	expectedPorts := []corev1.ContainerPort{
		{Name: "https", HostPort: 8443, ContainerPort: 443, Protocol: corev1.ProtocolTCP},
		{Name: "http", HostPort: 8000, ContainerPort: 80, Protocol: corev1.ProtocolTCP},
	}
	expectedLabelKey := "ports"
	expectedLabelExpr := MetadataLabelKeyPrefix + "." + expectedLabelKey + "=[{\"name\":\"https\",\"hostPort\":8443,\"containerPort\":443,\"protocol\":\"TCP\"},{\"name\":\"http\",\"hostPort\":8000,\"containerPort\":80,\"protocol\":\"TCP\"}]"
	args, err := forgeLabelMetadataArgs("ports", expectedPorts)
	require.NoError(t, err)
	require.Len(t, args, 2)
	assert.Equal(t, "--label", args[0])
	assert.Equal(t, expectedLabelExpr, args[1])
}

func Test_loadLabelMetadata(t *testing.T) {
	expectedMetadataLabelKey := "ports"
	labelKey := forgeLabelMetadataKey(expectedMetadataLabelKey)
	labelVal := "[{\"name\":\"https\",\"hostPort\":8443,\"containerPort\":443,\"protocol\":\"TCP\"},{\"name\":\"http\",\"hostPort\":8000,\"containerPort\":80,\"protocol\":\"TCP\"}]"
	labelsMap := map[string]any{
		"foo":    "pif",
		"bar":    "paf",
		labelKey: labelVal,
	}
	expectedPorts := []corev1.ContainerPort{
		{Name: "https", HostPort: 8443, ContainerPort: 443, Protocol: corev1.ProtocolTCP},
		{Name: "http", HostPort: 8000, ContainerPort: 80, Protocol: corev1.ProtocolTCP},
	}
	ports, err := loadLabelMetadata[[]corev1.ContainerPort](labelsMap, expectedMetadataLabelKey)
	require.NoError(t, err)
	require.NotNil(t, ports)
	assert.Equal(t, expectedPorts, ports)
}

func Test_podContainerNameFilter(t *testing.T) {
	var f string

	f = podContainerNameFilter("", "", "", false)
	assert.Equal(t, "name=^.+__.+__.+$", f)

	f = podContainerNameFilter("", "", "", true)
	assert.Equal(t, "name=^.+__.+--root$", f)

	f = podContainerNameFilter("ns1", "", "", false)
	assert.Equal(t, "name=^ns1__.+__.+$", f)

	f = podContainerNameFilter("ns1", "", "", true)
	assert.Equal(t, "name=^ns1__.+--root$", f)

	f = podContainerNameFilter("ns1", "pod1", "", false)
	assert.Equal(t, "name=^ns1__pod1__.+$", f)

	f = podContainerNameFilter("ns1", "pod1", "", true)
	assert.Equal(t, "name=^ns1__pod1--root$", f)

	f = podContainerNameFilter("ns1", "pod1", "ct1", false)
	assert.Equal(t, "name=^ns1__pod1__ct1$", f)

	f = podContainerNameFilter("ns1", "pod1", "ct1", true)
	assert.Equal(t, "name=^ns1__pod1--root$", f)
}

func TestCreateNamespace(t *testing.T) {
	// TODO
}

func TestUpdateNamespace(t *testing.T) {
	// TODO
}

func TestDeleteNamespace(t *testing.T) {
	// TODO
}

func TestListNamespaceNames(t *testing.T) {
	expectedNs1 := "ns1"
	rootCt1 := forgePodRootContainerName(expectedNs1, "pod1")
	expectedNs2 := "ns2"
	rootCt2 := forgePodRootContainerName(expectedNs2, "pod2")
	expectedNs3 := "ns3"
	rootCt3 := forgePodRootContainerName(expectedNs3, "pod3")
	expectedBinary := "foo"
	dt := DockerTranslater{expectedBinary}
	f := dt.ListNamespaceNames()

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "", "")
	r, e := f.Format()
	cmdz.StopMock()

	require.NoError(t, e)
	assert.Len(t, r, 0)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = rootCt1 + "\n" + rootCt2 + "\n" + rootCt3
		return
	})
	r, e = f.Format()
	cmdz.StopMock()

	require.NoError(t, e)
	require.Len(t, r, 3)
	assert.Contains(t, r, expectedNs1)
	assert.Contains(t, r, expectedNs2)
	assert.Contains(t, r, expectedNs3)

	// Error response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	r, e = f.Format()
	cmdz.StopMock()

	require.Error(t, e)
	assert.Len(t, r, 0)
}

func TestCreateVolume(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e1, err := dt.CreateVolume(expectedNamespace1, pod1Name, volume1)
	require.NoError(t, err)
	require.NotNil(t, e1)
	expectedCmd1 := fmt.Sprintf("%s volume create --driver local %s__%s", expectedBinary0, expectedNamespace1, volume1Name)
	assert.Equal(t, expectedCmd1, e1.String())

	e2, err := dt.CreateVolume(expectedNamespace1, pod1Name, volume2)
	require.NoError(t, err, "should not error")
	expectedCmd2 := fmt.Sprintf("%s volume create --driver local -o o=bind -o type=none -o device=%s %s__%s", expectedBinary0, volume2Path, expectedNamespace1, volume2Name)
	assert.Equal(t, expectedCmd2, e2.String())
}

func TestDeleteVolume(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e1 := dt.DeleteVolume(expectedNamespace1, pod1Name, volume1Name)
	require.NotNil(t, e1)
	expectedFilter := "^" + expectedNamespace1 + "__" + pod1Name + "__" + volume1Name + "$"
	expectedCmd1 := fmt.Sprintf(`sh -c %[1]s volume rm -f $( %[1]s volume ls -q -f name=%[2]s )`, expectedBinary0, expectedFilter)
	assert.Equal(t, expectedCmd1, e1.String())
}

func TestDescribeVolume(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.DescribeVolume(expectedNamespace1, pod1Name, volume1Name)
	require.NotNil(t, f)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "[]", "")
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 0)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = `
[
	{
		"CreatedAt": "2023-06-27T00:32:18+02:00",
		"Driver": "local",
		"Labels": null,
		"Mountpoint": "/var/lib/docker/volumes/ns1_ct5-vol1/_data",
		"Name": "ns1_ct5-vol1",
		"Options": null,
		"Scope": "local"
	}
]`
		return
	})
	res, err = f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Contains(t, res[0].Name, "ns1_ct5-vol1")

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f.Format()
	cmdz.StopMock()
	require.Error(t, err)
	assert.Nil(t, res)

}

func TestListVolumeNames(t *testing.T) {
	expectedNs1 := "ns1"
	vol1Name := forgePodVolumeName(expectedNs1, "pod1", "vol1")
	vol2Name := forgePodVolumeName(expectedNs1, "pod1", "vol2")

	dt := DockerTranslater{expectedBinary0}
	f := dt.ListVolumeNames(expectedNs1, "pod1")
	require.NotNil(t, f)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "", "")
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 0)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = vol1Name + "\n" + vol2Name + "\n"
		return
	})
	res, err = f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Contains(t, res, "vol1")
	assert.Contains(t, res, "vol2")

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f.Format()
	cmdz.StopMock()
	require.Error(t, err)
	require.Nil(t, res)
}

func TestSetupPod(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}

	/*
		cmds1, err := translator.createPodRootContainer(expectedNamespace1, pod1)
		require.NoError(t, err, "should not error")
		expectedCmd11 := fmt.Sprintf(`%s run --detach --name %[2]s__%[3]s--root --restart=always --network %[2]s__%[3]s--net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod1.Name)
		assert.Equal(t, expectedCmd11, cmds1.String())

		cmds2, err := translator.createPodRootContainer(expectedNamespace1, pod2)
		require.NoError(t, err, "should not error")
		expectedCmd21 := fmt.Sprintf(`%s run --detach --name %[2]s__%[3]s--root --restart=always --network %[2]s__%[3]s--net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod2.Name)
		assert.Equal(t, expectedCmd21, cmds2.String())
	*/

	e1 := dt.SetupPod(expectedNamespace1, pod1)
	require.NotNil(t, e1)
	expectedCmd11 := fmt.Sprintf(`%s run --detach --name %[2]s__%[3]s--root --restart=always --network %[2]s__%[3]s--net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod1.Name)
	assert.Equal(t, expectedCmd11, e1.String())

	e2 := dt.SetupPod(expectedNamespace1, pod2)
	require.NotNil(t, e2)
	expectedCmd21 := fmt.Sprintf(`%s run --detach --name %[2]s__%[3]s--root --restart=always --network %[2]s__%[3]s--net --cpus=0.05 --memory=64m alpine:3.17.3 /bin/sleep inf`, expectedBinary0, expectedNamespace1, pod2.Name)
	assert.Equal(t, expectedCmd21, e2.String())
}

func TestDeletePod(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e1 := dt.DeletePod(expectedNamespace1, pod1Name)
	require.NotNil(t, e1)
	expectedFilter := "^" + expectedNamespace1 + "__" + pod1Name + "__.+$"
	expectedCmd1 := fmt.Sprintf(`sh -c %[1]s rm -f $( %[1]s ps -q -f name=%[2]s )`, expectedBinary0, expectedFilter)
	assert.Equal(t, expectedCmd1, e1.String())
}

func Test_buildPodContainer_empty(t *testing.T) {
	var jsonResults []map[string]any
	err := yaml.Unmarshal([]byte("[]"), &jsonResults)
	assert.NoError(t, err)
	assert.Empty(t, jsonResults)
	// TODO
}

func Test_buildPodContainer_container1(t *testing.T) {
	var jsonResults []map[string]any
	err := yaml.Unmarshal([]byte(container1_docker_inspect), &jsonResults)
	assert.NoError(t, err)
	assert.NotNil(t, jsonResults)
	assert.NotEmpty(t, jsonResults)
	assert.Len(t, jsonResults, 1)

	ct, err := buildPodContainer(jsonResults[0])
	assert.NoError(t, err)
	assert.NotNil(t, ct)
	assert.Equal(t, "quizzical_hodgkin", ct.Name)
	// TODO
}

func TestDescribePod_empty(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f1 := dt.DescribePod(expectedNamespace1, pod1Name)
	require.NotNil(t, f1)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "[]", "")
	res, err := f1.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Nil(t, res, 0)

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f1.Format()
	cmdz.StopMock()
	require.Error(t, err)
	require.Nil(t, res)

}

func TestDescribePod_pod2(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f1 := dt.DescribePod(expectedNamespace1, pod2Name)
	require.NotNil(t, f1)

	// OK response from docker
	cmdz.StartSimpleMock(t, 0, pod2_docker_inspect, "")
	res, err := f1.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, pod2Name, res.Name)
	assert.Equal(t, expectedNamespace1, res.Namespace)

	require.NotNil(t, res.ObjectMeta)
	expectedLabels := map[string]string(nil)
	assert.Equal(t, expectedLabels, res.ObjectMeta.Labels)

	require.NotNil(t, res.Spec)
	expectedRestartPolicy := corev1.RestartPolicyAlways
	assert.Equal(t, expectedRestartPolicy, res.Spec.RestartPolicy)

	/*
		expectedPodSecurityContext := corev1.PodSecurityContext{
			Privileged:             boolPtr(false),
			ReadOnlyRootFilesystem: boolPtr(false),
			RunAsNonRoot:           boolPtr(false),
			RunAsUser:              nil,
			RunAsGroup:             nil,
		}
		assert.Equal(t, expectedSecurityContext, res.Spec.SecurityContext)
	*/

}

func TestCreatePodContainer_initContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}

	ct1 := pod1.Spec.InitContainers[0]
	e1, err := dt.CreatePodContainer(expectedNamespace1, pod1, ct1, true)
	require.NoError(t, err)
	require.NotNil(t, e1)
	expectedCmd1 := fmt.Sprintf(`%[1]s run --detach --rm --name %[2]s__%s__%s --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=always -v %[2]s__%[4]s-vol1:/foo/bar:ro -v %[2]s__%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 0.1`, expectedBinary0, expectedNamespace1, pod1.Name, ct1.Name, ct1.Image)
	assert.Equal(t, expectedCmd1, e1.String())

	ct2 := pod1.Spec.InitContainers[1]
	e2, err := dt.CreatePodContainer(expectedNamespace1, pod1, ct2, true)
	require.NoError(t, err, "should not error")
	expectedCmd2 := fmt.Sprintf(`%[1]s run --detach --rm --privileged --name %[2]s__%s__%s --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=missing -v %[2]s__%[4]s-vol1:/foo/bar:ro -v %[2]s__%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 10`, expectedBinary0, expectedNamespace1, pod1.Name, ct2.Name, ct2.Image)
	assert.Equal(t, expectedCmd2, e2.String())

	ct5 := pod3.Spec.InitContainers[1]
	e3, err := dt.CreatePodContainer(expectedNamespace2, pod3, ct5, true)
	require.NoError(t, err, "should not error")
	expectedCmd3 := fmt.Sprintf(`%s run --detach --rm --privileged --read-only --name %s__%s__%s --workdir=/tmp/%[4]s_workdir1 --restart=no --pull=never -v %[2]s__%[4]s-vol1:/foo/bar:ro -v %[2]s__%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sh -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s -c sleep 0.1`, expectedBinary0, expectedNamespace2, pod3.Name, ct5.Name, ct5.Image)
	assert.Equal(t, expectedCmd3, e3.String())
}

func TestCreatePodContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}

	ct3 := pod1.Spec.Containers[0]
	e1, err := dt.CreatePodContainer(expectedNamespace1, pod1, ct3, false)
	require.NoError(t, err)
	require.NotNil(t, e1)
	expectedCmd1 := fmt.Sprintf(`%[1]s run --detach --rm --read-only --name %[2]s__%s__%s --workdir=/tmp/%[4]s_workdir1 --restart=always --pull=always -v %[2]s__%[4]s-vol1:/foo/bar:ro -v %[2]s__%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 0.1`, expectedBinary0, expectedNamespace1, pod1.Name, ct3.Name, ct3.Image)
	assert.Equal(t, expectedCmd1, e1.String())

	ct4 := pod1.Spec.Containers[1]
	e2, err := dt.CreatePodContainer(expectedNamespace1, pod1, ct4, false)
	require.NoError(t, err, "should not error")
	expectedCmd2 := fmt.Sprintf(`%[1]s run --detach --rm -t --name %[2]s__%s__%s --workdir=/tmp/%[4]s_workdir1 --restart=always --pull=missing -v %[2]s__%[4]s-vol1:/foo/bar:ro -v %[2]s__%[4]s-vol2:/foo/baz:ro --entrypoint /bin/sleep -e "%[4]senvKey1=%[4]senvVal1" -e "%[4]senvKey2=%[4]senvVal2" -e "%[4]senvKey3=%[4]senvVal3" -l %[3]s_labelKey1=%[3]s_labelVal2 %[5]s 10`, expectedBinary0, expectedNamespace1, pod1.Name, ct4.Name, ct4.Image)
	assert.Equal(t, expectedCmd2, e2.String())
}

func TestUpdatePodContainer(t *testing.T) {
	t.Skip()
	dt := DockerTranslater{expectedBinary0}
	e, err := dt.UpdatePodContainer(expectedNamespace1, pod1, container1)
	require.NoError(t, err)
	require.NotNil(t, e)

	// TODO
}

func TestDeletePodContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e1 := dt.DeletePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, e1)
	expectedCmd1 := fmt.Sprintf(`%[1]s rm -f %[2]s__%[3]s__%[4]s`, expectedBinary0, expectedNamespace1, pod1Name, container1Name)
	assert.Equal(t, expectedCmd1, e1.String())
}

func TestDescribePodContainer_empty(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.DescribePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, f)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "[]", "")
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 0)

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f.Format()
	cmdz.StopMock()
	require.Error(t, err)
	require.Nil(t, res)
}

func TestDescribePodContainer_container1(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.DescribePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, f)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = container1_docker_inspect
		return
	})
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "quizzical_hodgkin", res[0].Name)
	assert.Equal(t, "busybox", res[0].Image)
	assert.Equal(t, "", res[0].WorkingDir)
	assert.Equal(t, false, res[0].TTY)
	assert.Equal(t, []string(nil), res[0].Command)
	assert.Equal(t, []string{"ls", "/world"}, res[0].Args)
	expectedEnv := []corev1.EnvVar{
		{
			Name:  "PATH",
			Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		},
	}
	assert.Equal(t, expectedEnv, res[0].Env)
	expectedPorts := []corev1.ContainerPort(nil)
	assert.Equal(t, expectedPorts, res[0].Ports)
	expectedSecurityContext := corev1.SecurityContext{
		Privileged:             boolPtr(false),
		ReadOnlyRootFilesystem: boolPtr(false),
		RunAsNonRoot:           boolPtr(false),
	}
	assert.Equal(t, &expectedSecurityContext, res[0].SecurityContext)
	expectedVolumeMounts := []corev1.VolumeMount{
		descriptor.BuildVolumeMount("hello", "/world"),
	}
	assert.Equal(t, expectedVolumeMounts, res[0].VolumeMounts)
}

func TestDescribePodContainer_container2(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.DescribePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, f)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = container2_docker_inspect
		return
	})
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "wonderful_faraday", res[0].Name)
	assert.Equal(t, "nginx", res[0].Image)
	assert.Equal(t, "/tmp", res[0].WorkingDir)
	assert.Equal(t, false, res[0].TTY)
	assert.Equal(t, []string{"/docker-entrypoint.sh"}, res[0].Command)
	assert.Equal(t, []string{"nginx", "-g", "daemon off;"}, res[0].Args)
	expectedEnv := []corev1.EnvVar{
		{Name: "PATH", Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		{Name: "NGINX_VERSION", Value: "1.25.1"},
		{Name: "NJS_VERSION", Value: "0.7.12"},
		{Name: "PKG_RELEASE", Value: "1~bookworm"},
	}
	assert.Equal(t, expectedEnv, res[0].Env)
	// TODO must put port names in labels ?
	expectedPorts := []corev1.ContainerPort{
		{Name: "https", HostPort: 8443, ContainerPort: 443, Protocol: corev1.ProtocolTCP},
		{Name: "http", HostPort: 8000, ContainerPort: 80, Protocol: corev1.ProtocolTCP},
	}
	assert.Equal(t, expectedPorts, res[0].Ports)
	expectedSecurityContext := corev1.SecurityContext{
		Privileged:             boolPtr(false),
		ReadOnlyRootFilesystem: boolPtr(false),
		RunAsNonRoot:           boolPtr(false),
		RunAsUser:              int64Ptr(10),
		RunAsGroup:             int64Ptr(12),
	}
	assert.Equal(t, &expectedSecurityContext, res[0].SecurityContext)
	expectedVolumeMounts := []corev1.VolumeMount{
		descriptor.BuildVolumeMount("foo", "/tmp/foo"),
		descriptor.BuildVolumeMount("bar", "/tmp/bar"),
	}
	assert.Equal(t, expectedVolumeMounts, res[0].VolumeMounts)
	expectedCpuLimit := resource.MustParse("0.2")
	expectedMemoryLimit := resource.MustParse("64Mi")
	assert.Equal(t, expectedCpuLimit.AsApproximateFloat64(), res[0].Resources.Limits.Cpu().AsApproximateFloat64())
	assert.Equal(t, expectedMemoryLimit.AsDec(), res[0].Resources.Limits.Memory().AsDec())
}

func TestListPodContainerNames(t *testing.T) {
	expectedNs1 := "ns1"
	expectedPod1 := "pod1"
	expectedCt1 := "ct1"
	expectedCt2 := "ct2"
	expectedCt3 := "ct3"
	ct1Name := forgePodContainerName(expectedNs1, expectedPod1, expectedCt1)
	ct2Name := forgePodContainerName(expectedNs1, expectedPod1, expectedCt2)
	ct3Name := forgePodContainerName(expectedNs1, expectedPod1, expectedCt3)

	dt := DockerTranslater{expectedBinary0}
	f := dt.ListPodContainerNames(expectedNs1, expectedPod1)
	require.NotNil(t, f)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "", "")
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 0)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = ct1Name + "\n" + ct2Name + "\n" + ct3Name + "\n"
		return
	})
	res, err = f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	assert.Len(t, res, 3)
	assert.Contains(t, res, expectedCt1)
	assert.Contains(t, res, expectedCt2)
	assert.Contains(t, res, expectedCt3)

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f.Format()
	cmdz.StopMock()
	require.Error(t, err)
	require.Nil(t, res)
}
