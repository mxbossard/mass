package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mby.fr/utils/cmdz"
)

//"mby.fr/utils/promise"

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

	// Error response from docker
	cmdz.StartSimpleMock(t, 1, "", "")
	r, e = f.Format()
	cmdz.StopMock()

	require.Error(t, e)
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
	e := dt.DeletePod(expectedNamespace1, pod1Name)
	require.NotNil(t, e)
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
	e := dt.DeletePod(expectedNamespace1, pod1Name)
	require.NotNil(t, e)
}

func TestInspectVolume(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.InspectVolume(expectedNamespace1, pod1Name, volume1Name)
	require.NotNil(t, f)
}

func TestListVolumeNames(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	f := dt.ListVolumeNames(expectedNamespace1, pod1Name)
	require.NotNil(t, f)
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
}

func TestDeletePodContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e := dt.DeletePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, e)
}

func TestInspectPodContainer(t *testing.T) {
	dt := DockerTranslater{expectedBinary0}
	e := dt.DeletePodContainer(expectedNamespace1, pod1Name, container1Name)
	require.NotNil(t, e)
}

func TestListPodContainerNames(t *testing.T) {

}
