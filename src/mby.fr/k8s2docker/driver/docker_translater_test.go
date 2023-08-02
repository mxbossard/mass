package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/cmdz"
)

//"mby.fr/utils/promise"

func TestPodContainerNameFilter(t *testing.T) {
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

func TestDescribePod2(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e1 := dt.DescribePod(expectedNamespace1, pod1Name)
	require.NotNil(t, e1)
	// TODO
}
func TestCreateVolume2(t *testing.T) {
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

func TestCreatePodContainer2_Init(t *testing.T) {
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

func TestCreatePodContainer2(t *testing.T) {
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

func TestDescribePodContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.DescribePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, f)

	// Empty response from docker
	cmdz.StartSimpleMock(t, 0, "[]", "")
	res, err := f.Format()
	cmdz.StopMock()
	require.NoError(t, err)
	require.Len(t, res, 0)

	// OK response from docker
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (rc int, stdout, stderr string) {
		stdout = container1_docker_inspect
		return
	})
	res, err = f.Format()
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
	// TODO add ports
	expectedPorts := []corev1.ContainerPort{}
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

	// Err response from docker
	cmdz.StartSimpleMock(t, 1, "", "some error")
	res, err = f.Format()
	cmdz.StopMock()
	require.Error(t, err)
	require.Nil(t, res)
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
