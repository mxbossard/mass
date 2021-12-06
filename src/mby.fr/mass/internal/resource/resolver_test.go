package resource

import(
	"testing"

	"mby.fr/mass/internal/workspace"
)

func TestResolvePath(t *testing.T) {
	wksPath := workspace.TestInitTempWorkspace(t)
	_ = wksPath

}
