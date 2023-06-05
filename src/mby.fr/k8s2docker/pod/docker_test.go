package pod

import (
	"testing"

	k8s "k8s.io/api"
	apimachinery "k8s.io/apimachinery/pkg"
)

func TestCreateVolume(t *testing.T) {
	translator := Translator{"docker"}
	translator.createVolume(expectedNamespace, volume)
}
