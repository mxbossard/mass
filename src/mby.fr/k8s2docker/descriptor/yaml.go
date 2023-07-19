package descriptor

import (

	//"gopkg.in/yaml.v3"
	//"encoding/json"
	"bytes"
	"log"
	"strings"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	//_ "k8s.io/client-go/pkg/api/install"
	//_ "k8s.io/client-go/pkg/apis/extensions/install"

	//"k8s.io/client-go/kubernetes/scheme"
	k8sYaml "k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	//kubernetesCoreV1 "k8s.io/kubernetes/pkg/apis/core/v1"
	//apiv1 "k8s.io/kubernetes/staging/src/k8s.io/kubernetes/pkg/apis/core/v1"
)

func ValidateResource(input []byte) (err error) {
	return
}

func LoadPod(input []byte) (pod k8sv1.Pod, err error) {
	err = ValidateResource(input)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(input, &pod)
	if err != nil {
		return
	}
	return
}

func LoadPodWithDefaults(data []byte) (pod k8sv1.Pod, err error) {
	//apiv1.SetDefaults_Pod(&pod)

	decoder := scheme.Codecs.UniversalDeserializer()
	obj, gKV, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return pod, err
	}
	if gKV.Kind == "Pod" {
		pod = *obj.(*corev1.Pod)
	}

	/*
		err = kubernetesCoreV1.Convert_v1_Service_To_core_Service(service, kubernetesService, nil)
		if err != nil {
			return pod, err
		}
		err = kubernetesCoreValidation.ValidateService(kubernetesService)
		if err != nil {
			return pod, err
		}
	*/

	return
}

func LoadPodWithDefaults2(data []byte) (k8sv1.Pod, error) {
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDecoder()
	object := k8sv1.Pod{}
	err := runtime.DecodeInto(decoder, data, &object)
	if err != nil {
		return object, err
	}
	return object, nil
}

func LoadPodWithDefaults0(input []byte) (pod k8sv1.Pod, err error) {
	pod, err = LoadPod(input)
	if err != nil {
		return
	}

	if k8sv1.RestartPolicy("") == pod.Spec.RestartPolicy {
		pod.Spec.RestartPolicy = k8sv1.RestartPolicyAlways
	}
	return
}

func LoadPodWithDefaults1(input []byte) (pod k8sv1.Pod, err error) {
	// Create a runtime.Decoder from the Codecs field within
	// k8s.io/client-go that's pre-loaded with the schemas for all
	// the standard Kubernetes resource types.
	//decoder := scheme.Codecs.UniversalDeserializer()

	//scheme := runtime.NewScheme()
	//codecs := serializer.NewCodecFactory(scheme)
	//decoder := codecs.UniversalDeserializer()
	//defaults := FromAPIVersionAndKind("v1", "Pod")

	for _, resourceYAML := range strings.Split(string(input), "---") {
		// skip empty documents, `Decode` will fail on them
		if len(resourceYAML) == 0 {
			continue
		}

		// - obj is the API object (e.g., Deployment)
		// - groupVersionKind is a generic object that allows
		//   detecting the API type we are dealing with, for
		//   accurate type casting later.

		/*
			obj, groupVersionKind, err := decoder.Decode(
				[]byte(resourceYAML),
				nil,
				nil)
		*/

		decoder := k8sYaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(resourceYAML)), 1000)
		err = decoder.Decode(&pod)

		if err != nil {
			log.Print(err)
			return pod, err
			continue
		}

		// Figure out from `Kind` the resource type, and attempt
		// to cast appropriately.
		/*
			if groupVersionKind.Version == "v1" && //groupVersionKind.Group == "apps" &&
				groupVersionKind.Kind == "Pod" {
				if p, ok := obj.(*k8sv1.Pod); ok {
					//log.Print(p)
					return *p, err
				}
			}
		*/
	}
	return
}
