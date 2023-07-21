package descriptor

import (

	//"gopkg.in/yaml.v3"
	//"encoding/json"

	"fmt"
	"log"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"mby.fr/utils/collections"
	"mby.fr/utils/errorz"
)

func assertNotEmpty(errors *errorz.Aggregated, element any, field, kind, name string) {
	//log.Printf("assertNotEmpty field: %s of kind: %s => val: [%v] of type: %T", field, kind, element, element)
	var invalid bool
	switch t := element.(type) {
	case string:
		invalid = t == ""
	case int, int8, int16, int32, int64, float32, float64:
		invalid = t == 0
	default:
		log.Printf("assertingNotEmty default for field: %s of kind: %s", field, kind)
		invalid = t == nil
	}
	if invalid {
		var err error
		if name != "" {
			err = fmt.Errorf("Missing field %s in resource: %s{name=%s} !", field, kind, name)
		} else {
			err = fmt.Errorf("Missing field %s in resource of kind: %s !", field, kind)
		}
		errors.Add(err)
	}
}

func assertIn[T any](errors *errorz.Aggregated, element T, field, kind, name string, possibilities ...T) {
	if !collections.ContainsAny[T](&possibilities, element) {
		var err error
		if name != "" {
			err = fmt.Errorf("Invalid field %s value: %v not in expected list: %v in resource: %s{name=%s} !", field, element, possibilities, kind, name)
		} else {
			err = fmt.Errorf("Invalid field %s value: %v not in expected list: %v in resource: of kind %s !", field, element, possibilities, kind)
		}
		errors.Add(err)
	}
}

func ValidateNamespace(res corev1.Namespace, name string) error {
	errors := errorz.Aggregated{}
	errors.Concat(ValidateTypeMeta(res.TypeMeta, "Namespace", name))
	errors.Concat(ValidateObjectMeta(res.ObjectMeta, "Namespace", name))
	return errors.Return()
}

func ValidatePod(res corev1.Pod, name string) error {
	errors := errorz.Aggregated{}
	errors.Concat(ValidateTypeMeta(res.TypeMeta, "Pod", name))
	errors.Concat(ValidateObjectMeta(res.ObjectMeta, "Pod", name))
	errors.Concat(ValidatePodSpec(res.Spec, "spec", "Pod", name))
	return errors.Return()
}

func ValidateTypeMeta(res metav1.TypeMeta, kind, name string) (errors errorz.Aggregated) {
	assertNotEmpty(&errors, res.Kind, "kind", kind, name)
	assertNotEmpty(&errors, res.APIVersion, "apiVersion", kind, name)
	return
}

func ValidateObjectMeta(res metav1.ObjectMeta, kind, name string) (errors errorz.Aggregated) {
	assertNotEmpty(&errors, res.Name, "metadata.name", kind, name)
	return
}

func ValidatePodSpec(res corev1.PodSpec, field, kind, name string) (errors errorz.Aggregated) {
	//assertNotEmpty(&errors, res.RestartPolicy, field+".restartPolicy", kind, name)
	assertIn(&errors, res.RestartPolicy, field+".restartPolicy", kind, name, corev1.RestartPolicyAlways, corev1.RestartPolicyOnFailure, corev1.RestartPolicyNever)
	for id, ct := range res.InitContainers {
		f := fmt.Sprintf("%s.initContainers[%d]", field, id)
		errors.Concat(ValidateContainer(ct, f, kind, name))
	}
	for id, ct := range res.Containers {
		f := fmt.Sprintf("%s.containers[%d]", field, id)
		errors.Concat(ValidateContainer(ct, f, kind, name))
	}
	for id, vol := range res.Volumes {
		f := fmt.Sprintf("%s.volumes[%d]", field, id)
		errors.Concat(ValidateVolume(vol, f, kind, name))
	}
	return
}

func ValidateContainer(res corev1.Container, field, kind, name string) (errors errorz.Aggregated) {
	assertNotEmpty(&errors, res.Name, field+".name", kind, name)
	assertNotEmpty(&errors, res.Image, field+".image", kind, name)
	//assertNotEmpty(&errors, res.ImagePullPolicy, field+".imagePullPolicy", kind, name)
	assertIn(&errors, res.ImagePullPolicy, field+".imagePullPolicy", kind, name, corev1.PullAlways, corev1.PullNever, corev1.PullIfNotPresent)

	errors.Concat(ValidateSecurityContext(res.SecurityContext, field+".securityContext", kind, name))
	for id, vol := range res.VolumeMounts {
		f := fmt.Sprintf("%s.volumeMounts[%d]", field, id)
		errors.Concat(ValidateVolumeMount(vol, f, kind, name))
	}
	errors.Concat(ValidateProbe(res.LivenessProbe, field+".livenessProbe", kind, name))
	errors.Concat(ValidateProbe(res.ReadinessProbe, field+".readynessProbe", kind, name))
	errors.Concat(ValidateProbe(res.StartupProbe, field+".startupProbe", kind, name))
	return
}

func ValidateVolume(res corev1.Volume, field, kind, name string) (errors errorz.Aggregated) {
	assertNotEmpty(&errors, res.Name, field+".name", kind, name)
	return
}

func ValidateVolumeMount(res corev1.VolumeMount, field, kind, name string) (errors errorz.Aggregated) {
	assertNotEmpty(&errors, res.Name, field+".name", kind, name)
	assertNotEmpty(&errors, res.MountPath, field+".mountPath", kind, name)
	return
}

func ValidateProbe(res *corev1.Probe, field, kind, name string) (errors errorz.Aggregated) {
	if res != nil {
		// TODO
	}
	return
}
func ValidateSecurityContext(res *corev1.SecurityContext, field, kind, name string) (errors errorz.Aggregated) {
	if res != nil {
		// TODO
	}
	return
}

func ValidateSerializedK8sResource(input []byte, kind, name string) (err error) {
	return
}

func ValidateMappedK8sResource(input map[string]any, kind, name string) (err error) {
	return
}

func ValidateK8sResource[T any](input T, name string) (validated T, err error) {
	err = CompleteK8sResourceDefaults(&input)
	if err != nil {
		return
	}

	switch res := any(input).(type) {
	case corev1.Namespace:
		err = ValidateNamespace(res, name)
		validated = any(input).(T)
	case corev1.Pod:
		err = ValidatePod(res, name)
		validated = any(input).(T)
	default:
		err = fmt.Errorf("Cannot validate resource of type: %T ! Not supported yet !", res)
	}
	return
}
