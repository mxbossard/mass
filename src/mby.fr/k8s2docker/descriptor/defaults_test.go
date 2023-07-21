package descriptor

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestCompleteContainerDefaults(t *testing.T) {
	ct := corev1.Container{}
	CompleteContainerDefaults(&ct)

	assert.Equal(t, Default_Container_TerminationMessagePath, ct.TerminationMessagePath)
	assert.Equal(t, Default_Container_TerminationMessagePolicy, ct.TerminationMessagePolicy)
	assert.Equal(t, Default_Container_ImagePullPolicy, ct.ImagePullPolicy)
}

func TestCompletePodSpecDefaults(t *testing.T) {
	spec := corev1.PodSpec{}
	CompletePodSpecDefaults(&spec)

	assert.Equal(t, Default_PodSpec_RestartPolicy, spec.RestartPolicy)
	assert.Equal(t, Default_PodSpec_DNSPolicy, spec.DNSPolicy)
	assert.Equal(t, Default_PodSpec_SchedulerName, spec.SchedulerName)
	assert.Equal(t, Default_PodSpec_PriorityClassName, spec.PriorityClassName)

	assert.Equal(t, Default_PodSpec_TerminationGracePeriodSeconds, spec.TerminationGracePeriodSeconds)
	assert.Equal(t, Default_PodSpec_PreemptionPolicy(), spec.PreemptionPolicy)
}
