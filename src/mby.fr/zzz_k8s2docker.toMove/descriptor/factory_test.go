package descriptor

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestCompare_Pods(t *testing.T) {
	p1a := BuildDefaultPod("ns1", "p1")
	p1b := BuildDefaultPod("ns1", "p1")
	p2 := BuildDefaultPod("ns1", "p2")

	assert.Equal(t, p1a, p1b)
	assert.NotEqual(t, p1a, p2)
}

func TestCompare_Containers(t *testing.T) {
	ct1a := BuildDefaultContainer("ct1", "image1")
	ct1b := BuildDefaultContainer("ct1", "image1")
	ct1c := BuildDefaultContainer("ct1", "image2")
	ct2 := BuildDefaultContainer("ct2", "image1")

	assert.Equal(t, ct1a, ct1b)
	assert.NotEqual(t, ct1a, ct1c)
	assert.NotEqual(t, ct1a, ct2)
	assert.NotEqual(t, ct1c, ct2)

}

func TestCompare_PodsWithContainers(t *testing.T) {
	p1a := BuildDefaultPod("ns1", "p1")
	ct1a := BuildDefaultContainer("ct1", "image1")
	p1a.Spec.Containers = append(p1a.Spec.Containers, ct1a)
	p1b := BuildDefaultPod("ns1", "p1")
	ct1b := BuildDefaultContainer("ct1", "image1")
	p1b.Spec.Containers = append(p1b.Spec.Containers, ct1b)
	p1c := BuildDefaultPod("ns1", "p1")
	ct1c := BuildDefaultContainer("ct1", "image2")
	p1c.Spec.Containers = append(p1c.Spec.Containers, ct1c)

	assert.Equal(t, p1a, p1b)
	assert.NotEqual(t, p1a, p1c)

	p1a.Spec.Containers = []corev1.Container{ct1a}
	p1b.Spec.Containers = []corev1.Container{ct1b}
	assert.Equal(t, p1a, p1b)

}
