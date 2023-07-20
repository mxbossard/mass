package daemon

import (
	"log"
	"os"
	"testing"
	"time"

	"mby.fr/k8s2docker/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemon_1(t *testing.T) {
	log.Printf("----- TestDaemon_1 -----")
	os.RemoveAll("./myTestDb")
	repo.InitDb("./myTestDb")
	//defer os.RemoveAll("./myTestDb")

	var err error
	_, err = repo.Put("ns1", "Pod", "pod1", "")
	require.NoError(t, err)

	Start(90*time.Millisecond, 180*time.Millisecond)
	//defer Stop()

	time.Sleep(300 * time.Millisecond)

	BlockingStop()
	assert.Fail(t, "daemon stopped")
}

func TestDaemon_2(t *testing.T) {
	log.Printf("----- TestDaemon_2 -----")
	repo.InitDb("./myTestDb")
	//defer os.RemoveAll("./myTestDb")

	var err error
	_, err = repo.Delete("ns1", "Pod", "pod1", "")
	require.NoError(t, err)

	Start(90*time.Millisecond, 180*time.Millisecond)
	//defer Stop()

	time.Sleep(300 * time.Millisecond)

	BlockingStop()
	assert.Fail(t, "daemon stopped")
}
