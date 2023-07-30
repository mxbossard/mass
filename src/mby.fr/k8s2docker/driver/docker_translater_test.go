package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mby.fr/utils/cmdz"
)

//"mby.fr/utils/promise"

func TestListNamespaceNames(t *testing.T) {
	expectedBinary := "foo"
	expectedNs1 := "ns1"
	rootCt1 := forgePodRootContainerName(expectedNs1, "pod1")
	expectedNs2 := "ns2"
	rootCt2 := forgePodRootContainerName(expectedNs2, "pod2")
	expectedNs3 := "ns3"
	rootCt3 := forgePodRootContainerName(expectedNs3, "pod3")
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
