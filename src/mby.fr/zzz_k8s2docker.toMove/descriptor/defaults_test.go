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

func TestCompletePodSpecDefaults_Empty(t *testing.T) {
	spec := corev1.PodSpec{}
	CompletePodSpecDefaults(&spec)

	assert.Equal(t, Default_PodSpec_RestartPolicy, spec.RestartPolicy)
	assert.Equal(t, Default_PodSpec_DNSPolicy, spec.DNSPolicy)
	assert.Equal(t, Default_PodSpec_SchedulerName, spec.SchedulerName)
	assert.Equal(t, Default_PodSpec_PriorityClassName, spec.PriorityClassName)

	assert.Equal(t, Default_PodSpec_TerminationGracePeriodSeconds, spec.TerminationGracePeriodSeconds)
	assert.Equal(t, Default_PodSpec_PreemptionPolicy(), spec.PreemptionPolicy)
}

func TestCompletePodSpecDefaults_Container(t *testing.T) {
	specWithContainer := corev1.PodSpec{Containers: []corev1.Container{
		corev1.Container{
			Name:  "foo",
			Image: "alpine",
		},
	}}
	CompletePodSpecDefaults(&specWithContainer)

	assert.Equal(t, Default_PodSpec_RestartPolicy, specWithContainer.RestartPolicy)
	assert.Equal(t, Default_PodSpec_DNSPolicy, specWithContainer.DNSPolicy)
	assert.Equal(t, Default_PodSpec_SchedulerName, specWithContainer.SchedulerName)
	assert.Equal(t, Default_PodSpec_PriorityClassName, specWithContainer.PriorityClassName)

	assert.Equal(t, Default_PodSpec_TerminationGracePeriodSeconds, specWithContainer.TerminationGracePeriodSeconds)
	assert.Equal(t, Default_PodSpec_PreemptionPolicy(), specWithContainer.PreemptionPolicy)

	assert.Equal(t, Default_Container_TerminationMessagePath, specWithContainer.Containers[0].TerminationMessagePath)
	assert.Equal(t, Default_Container_TerminationMessagePolicy, specWithContainer.Containers[0].TerminationMessagePolicy)
	assert.Equal(t, Default_Container_ImagePullPolicy, specWithContainer.Containers[0].ImagePullPolicy)
}
