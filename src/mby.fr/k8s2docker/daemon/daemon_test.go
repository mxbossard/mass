package daemon

import (
	"os"
	"testing"
	"time"

	"mby.fr/k8s2docker/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemon(t *testing.T) {
	repo.InitDb("./myTestDb")
	defer os.RemoveAll("./myTestDb")

	var err error
	//_, err = repo.Put("ns1", "Pod", "pod1", "")
	require.NoError(t, err)

	Start()
	//defer Stop()

	time.Sleep(300 * time.Millisecond)

	BlockingStop()
	assert.Fail(t, "daemon stopped")
}
