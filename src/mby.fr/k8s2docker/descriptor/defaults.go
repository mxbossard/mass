package descriptor

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	Default_Container_TerminationMessagePath   = "/dev/termination-log"
	Default_Container_TerminationMessagePolicy = corev1.TerminationMessageReadFile
	Default_Container_ImagePullPolicy          = corev1.PullAlways

	Default_PodSpec_RestartPolicy = corev1.RestartPolicyAlways
	Default_PodSpec_DNSPolicy     = corev1.DNSClusterFirst

	Default_PodSpec_SchedulerName     = "scheduler"
	Default_PodSpec_PriorityClassName = "default"
)

var (
	Default_PodSpec_TerminationGracePeriodSeconds = int64Ptr(int64(30))
)

func int64Ptr(in int64) *int64 {
	return &in
}

func Default_PodSpec_PreemptionPolicy() *corev1.PreemptionPolicy {
	v := corev1.PreemptLowerPriority
	return &v
}

func CompleteContainerDefaults(res *corev1.Container) {
	if res.TerminationMessagePath == "" {
		res.TerminationMessagePath = Default_Container_TerminationMessagePath
	}
	if res.TerminationMessagePolicy == "" {
		res.TerminationMessagePolicy = Default_Container_TerminationMessagePolicy
	}
	if res.ImagePullPolicy == "" {
		res.ImagePullPolicy = Default_Container_ImagePullPolicy
	}
}

func CompletePodSpecDefaults(res *corev1.PodSpec) {
	if res.RestartPolicy == "" {
		res.RestartPolicy = Default_PodSpec_RestartPolicy
	}
	if res.TerminationGracePeriodSeconds == nil {
		res.TerminationGracePeriodSeconds = Default_PodSpec_TerminationGracePeriodSeconds
	}
	if res.DNSPolicy == "" {
		res.DNSPolicy = Default_PodSpec_DNSPolicy
	}
	if res.SchedulerName == "" {
		res.SchedulerName = Default_PodSpec_SchedulerName
	}
	if res.PriorityClassName == "" {
		res.PriorityClassName = Default_PodSpec_PriorityClassName
	}
	if res.PreemptionPolicy == nil {
		res.PreemptionPolicy = Default_PodSpec_PreemptionPolicy()
	}

	for _, ct := range res.InitContainers {
		CompleteContainerDefaults(&ct)
	}
	for _, ct := range res.Containers {
		CompleteContainerDefaults(&ct)
	}
}

func CompleteK8sResourceDefaults[T any](input *T) (err error) {
	switch res := any(input).(type) {
	case *corev1.Namespace:
		// Nothing to do
	case *corev1.Pod:
		CompletePodSpecDefaults(&res.Spec)
	default:
		err = fmt.Errorf("Cannot complete resource of type: %T ! Not supported yet !", res)
	}
	return
}
