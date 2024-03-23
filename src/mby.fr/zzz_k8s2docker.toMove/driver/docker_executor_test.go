package driver

import (
	//"fmt"
	//"strings"
	"testing"

	//k8sv1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/cmdz"
	//"mby.fr/utils/collections"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyPod(t *testing.T) {
	translator := Translator{expectedBinary}
	executor := Executor{translator: translator}
	err := executor.updatePod(expectedNamespace1, pod1)
	require.NoError(t, err)

}

func TestDescribePod(t *testing.T) {
	defer cmdz.StopMock()
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (int, string, string) {
		//ns1__p2__bar-ct bar-image2\nns2__p1__bar-ct bar-image3\n
		return 0, "ns1__p1__foo-ct foo-image\nns1__p1__bar-ct bar-image\n", ""
	})

	expectedFooCt := descriptor.BuildDefaultContainer("foo-ct", "foo-image")
	expectedBarCt := descriptor.BuildDefaultContainer("bar-ct", "bar-image")
	expectedP1Pod := descriptor.BuildDefaultPod("ns1", "p1")
	expectedP1Pod.Spec.Containers = append(expectedP1Pod.Spec.Containers, expectedBarCt, expectedFooCt)

	translator := Translator{expectedBinary}
	executor := Executor{translator: translator}

	pod, err := executor.describePod("ns1", "p1")
	require.NoError(t, err)
	assert.Equal(t, expectedP1Pod, pod)
}

func TestListPods(t *testing.T) {
	defer cmdz.StopMock()
	cmdz.StartStringMock(t, func(c cmdz.Vcmd) (int, string, string) {
		if cmdz.Contains(c, "ns1", "p1") {
			return 0, "ns1__p1__foo-ct foo-image\nns1__p1__bar-ct bar-image\n", ""
		} else if cmdz.Contains(c, "ns1", "p2") {
			return 0, "ns1__p2__bar-ct bar-image2\n", ""
		} else if cmdz.Contains(c, "ns1") {
			return 0, "ns1__p1__foo-ct foo-image\nns1__p1__bar-ct bar-image\nns1__p2__bar-ct bar-image2\n", ""
		} else if cmdz.Contains(c, "ns2") {
			return 0, "ns2__p1__bar-ct bar-image3\n", ""
		}
		require.Fail(t, "Bad mocking !")
		return -1, "", ""
	})

	expectedFooCt := descriptor.BuildDefaultContainer("foo-ct", "foo-image")
	expectedBarCt := descriptor.BuildDefaultContainer("bar-ct", "bar-image")
	expectedBarCt2 := descriptor.BuildDefaultContainer("bar-ct", "bar-image2")
	expectedBarCt3 := descriptor.BuildDefaultContainer("bar-ct", "bar-image3")
	expectedP1Pod := descriptor.BuildDefaultPod("ns1", "p1")
	expectedP1Pod.Spec.Containers = append(expectedP1Pod.Spec.Containers, expectedBarCt, expectedFooCt)
	expectedP2Pod := descriptor.BuildDefaultPod("ns1", "p2")
	expectedP2Pod.Spec.Containers = append(expectedP2Pod.Spec.Containers, expectedBarCt2)
	expectedP3Pod := descriptor.BuildDefaultPod("ns2", "p1")
	expectedP3Pod.Spec.Containers = append(expectedP3Pod.Spec.Containers, expectedBarCt3)

	translator := Translator{expectedBinary}
	executor := Executor{translator: translator}

	pods, err := executor.ListPods("ns1")
	require.NoError(t, err)
	require.Len(t, pods, 2)
	assert.Contains(t, pods, expectedP1Pod)
	assert.Contains(t, pods, expectedP2Pod)

	pods, err = executor.ListPods("ns2")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	assert.Contains(t, pods, expectedP3Pod)
}
