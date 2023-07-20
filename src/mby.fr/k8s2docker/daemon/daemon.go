package daemon

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	k8sYaml "k8s.io/apimachinery/pkg/util/yaml"

	"mby.fr/k8s2docker/driver"
	"mby.fr/k8s2docker/repo"
	"mby.fr/utils/collections"
	"mby.fr/utils/errorz"
)

var dockerDriver = driver.DockerExecutor()
var stopped = true
var resourceDaemonRunning, probeDaemonRunning bool

func Start(resourcePeriod, probePeriod time.Duration) {
	if !stopped {
		return
	}
	stopped = false
	resourceDaemonRunning = true
	probeDaemonRunning = false
	go runResourcesDaemon(resourcePeriod)
	go runProbeDaemon(probePeriod)
}

func Stop() {
	stopped = true
}

func BlockingStop() {
	Stop()
	for resourceDaemonRunning || probeDaemonRunning {
		time.Sleep(10 * time.Millisecond)
	}
}

func runResourcesDaemon(period time.Duration) {
	for !stopped {
		log.Printf("### Start processing resources ...")
		err := processResources()
		if err.GotError() {
			log.Printf("ResourcesDaemon ERROR: %s", err)
		}
		log.Printf("### Finished processing resources.")
		time.Sleep(period)
	}
	resourceDaemonRunning = false
}

func runProbeDaemon(period time.Duration) {
	for !stopped {
		time.Sleep(period)
		log.Printf("--- Start probing ...")
		err := probeContainers()
		if err != nil {
			log.Printf("ProbeDaemon ERROR: %s", err)
		}
		log.Printf("--- Finished probing.")
	}
	probeDaemonRunning = false
}

func processResources() errorz.Aggregated {
	existingNamespaces, err := dockerDriver.ListNamespaces()
	if err != nil {
		return errorz.NewAggregated(err)
	}
	existingNamespaceNames := collections.Map[corev1.Namespace, string](existingNamespaces, func(ns corev1.Namespace) string {
		return ns.ObjectMeta.Name
	})
	declaredNamespaceNames, err := repo.ListNamespaceNames()
	if err != nil {
		return errorz.NewAggregated(err)
	}

	/*
		namespaceNames := collections.Deduplicate(&existingNamespaceNames, &declaredNamespaceNames)
		log.Printf("Existing namespaces: %v", existingNamespaceNames)
		log.Printf("Declared namespaces: %v", declaredNamespaceNames)
		log.Printf("Deduplicated namespaces: %v", namespaceNames)
	*/

	errs := errorz.Aggregated{}

	deletedNamespaceNames := collections.KeepLeft(&existingNamespaceNames, &declaredNamespaceNames)
	for _, ns := range deletedNamespaceNames {
		log.Printf("Deleting namespace: %s ...", ns)
	}

	for _, ns := range declaredNamespaceNames {
		existingPods, err := dockerDriver.ListPods(ns)
		if err != nil {
			return errorz.NewAggregated(err)
		}
		log.Printf("Existing pods in ns %s: %v", ns, existingPods)

		declaredResources, err := repo.Get(ns, "Pod", "", "")
		if err != nil {
			return errorz.NewAggregated(err)
		}
		declaredPods, err := convertMapSliceToK8sSlice[corev1.Pod](declaredResources)
		if err != nil {
			return errorz.NewAggregated(err)
		}

		agerr := processPods(ns, &existingPods, &declaredPods)
		errs.Concat(agerr)
	}

	return errs
}

func processPods(ns string, existingPods *[]corev1.Pod, declaredPods *[]corev1.Pod) (errs errorz.Aggregated) {
	// Apply each declared pods
	for _, pod := range *declaredPods {
		err := applyPod(ns, existingPods, pod)
		errs.Add(err)
		// TODO: store config
	}

	// Delete not declared anymore pods
	existingPodNames := collections.Map[corev1.Pod, string](*existingPods, func(pod corev1.Pod) string {
		return pod.ObjectMeta.Name
	})
	declaredPodNames := collections.Map[corev1.Pod, string](*declaredPods, func(pod corev1.Pod) string {
		return pod.ObjectMeta.Name
	})

	deletedPodNames := collections.KeepLeft[string](&existingPodNames, &declaredPodNames)
	for _, deletedPodName := range deletedPodNames {
		log.Printf("Pod %s was removed.", deletedPodName)
		err := dockerDriver.Delete(ns, "Pod", deletedPodName)
		errs.Add(err)
		// TODO: clear config
	}
	return
}

func applyPod(ns string, existingPods *[]corev1.Pod, declaredPod corev1.Pod) (err error) {
	declaredPodName := declaredPod.ObjectMeta.Name
	log.Printf("Processing pod: %s", declaredPodName)
	if collections.ContainsAny[corev1.Pod](existingPods, declaredPod) {
		// Found identical spec already deployed
		log.Printf("Pod %s already exists and is untouched.", declaredPodName)
	} else {
		match := collections.Filter(*existingPods, func(p corev1.Pod) bool {
			return p.ObjectMeta.Name == declaredPodName
		})
		if len(match) > 0 {
			// Found a matching different spec (same name)
			log.Printf("Pod %s already exists but was modified ([%v] => [%v]).", declaredPodName, match, declaredPod)
			err = dockerDriver.Apply(ns, declaredPod)
			if err != nil {
				return
			}
		} else {
			// New pod
			log.Printf("Pod %s does not exists yet.", declaredPodName)
			err = dockerDriver.Apply(ns, declaredPod)
			if err != nil {
				return
			}
		}
	}
	return
}

func probeContainers() (err error) {
	return
}

func convertMapToK8sRes[T any](res map[string]any) (resource T, err error) {
	buffer, err := yaml.Marshal(res)
	if err != nil {
		return
	}
	err = k8sYaml.Unmarshal(buffer, &resource)
	return
}

func convertMapSliceToK8sSlice[T any](slice []map[string]any) (resources []T, err error) {
	for _, item := range slice {
		res, err := convertMapToK8sRes[T](item)
		if err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return
}

func convertAnySliceToK8sSlice[T any](slice []any) (resources []T, err error) {
	for _, item := range slice {
		if res, ok := item.(T); ok {
			resources = append(resources, res)
		} else {
			err = fmt.Errorf("Unable to convert res slice: [%s] into type: %T", slice, resources)
			return
		}
	}
	return
}
