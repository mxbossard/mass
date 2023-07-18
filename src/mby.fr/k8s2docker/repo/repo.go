package repo

import (
	"encoding/json"
	_ "errors"
	"fmt"
	"log"
	_ "os"
	"strings"

	"gopkg.in/yaml.v3"
	//corev1 "k8s.io/api/core/v1"

	"mby.fr/utils/collections"
	"mby.fr/utils/serializ"

	scribble "mby.fr/scribble"
)

var (
	db                 *scribble.Driver
	meta_ns_collection = "__meta"
	//dbDir = "mydb"
)

func init() {
	//initDb()
}

func initDb(dbDirPath string) {
	if db != nil {
		return
	}
	var err error
	db, err = scribble.New(dbDirPath, nil)
	if err != nil {
		log.Fatalf("Error initializing DB: %s", err)
	}
}

func Interract(namespace, kind, name, jsonIn, method string) (jsonOut string, err error) {
	if kind != "Pod" && kind != "Namespace" {
		// Not supported resource kind
		err = fmt.Errorf("Not supported kind: %s !", kind)
	}

	var resources []map[string]any
	if method == "GET" {
		resources, err = Get(namespace, kind, name, jsonIn)
	} else if method == "POST" {
		resources, err = Post(namespace, kind, name, jsonIn)
	} else if method == "PUT" {
		resources, err = Put(namespace, kind, name, jsonIn)
	} else if method == "PATCH" {
		resources, err = Patch(namespace, kind, name, jsonIn)
	} else if method == "DELETE" {
		resources, err = Delete(namespace, kind, name, jsonIn)
	} else {
		err = fmt.Errorf("Not supported method: %s !", method)
	}
	if err != nil {
		return "", err
	}

	return mappedResourcesToJson(resources)
}

func mappedResourcesToJson(resources []map[string]any) (string, error) {
	var mappingError error
	jsonResources := collections.Map(resources, func(i map[string]any) string {
		outBytes, err := json.Marshal(i)
		if err != nil {
			mappingError = err
			return ""
		}
		return string(outBytes)
	})
	if mappingError != nil {
		return "", fmt.Errorf("Unable to map resources ! Caused by: %w", mappingError)
	}
	out := strings.Join(jsonResources, "\n---\n")
	return out, mappingError
}

func forgeNamespace(name string) map[string]any {
	ns := map[string]any{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]any{
			"name": name,
		},
	}
	return ns
}

func forgeResource(namespace, kind, name string) map[string]any {
	if kind == "" {
		return forgeNamespace(namespace)
	}
	ns := map[string]any{
		"apiVersion": "v1",
		"kind":       kind,
		"metadata": map[string]any{
			"name":      name,
			"namespace": namespace,
		},
	}
	return ns
}

// Forge resource collection name
func forgeResCollection(namespace, kind string) string {
	return fmt.Sprintf("%s___%s", namespace, kind)
}

func loadResource(namespace, kind, name string) (map[string]any, error) {
	if kind == "Namespace" {
		namespace = meta_ns_collection
	}
	collection := forgeResCollection(namespace, kind)
	return scribble.Read[map[string]any](db, collection, name)
}

func storeResource(namespace string, resourceTree map[string]any) ([]map[string]any, error) {
	var storedResources []map[string]any
	// Read kind from tree
	kind, err := serializ.ResolveJsonMap[string](resourceTree, "/kind")
	if err != nil {
		return nil, err
	}
	// Read name from tree
	name, err := serializ.ResolveJsonMap[string](resourceTree, "/metadata/name")
	if err != nil {
		return nil, err
	}
	// swallow error namespace is optionnal in tree
	jsonNamespace, _ := serializ.ResolveJsonMap[string](resourceTree, "/metadata/namespace")
	if jsonNamespace != "" {
		namespace = jsonNamespace
	}

	if kind == "Namespace" {
		namespace = meta_ns_collection
	} else {
		// Ensure resources collection is referenced
		recordResourcesCollection(namespace, kind)
		// Ensure namespace exists, creating it
		nsJson := forgeNamespace(namespace)
		res, err := storeResource(meta_ns_collection, nsJson)
		if err != nil {
			return nil, err
		}
		storedResources = append(storedResources, res...)
	}

	collection := forgeResCollection(namespace, kind)
	log.Printf("Storing %s %s: %s\n", collection, name, resourceTree)
	err = scribble.Write(db, collection, name, resourceTree)
	if err != nil {
		return storedResources, err
	}
	storedResources = append(storedResources, resourceTree)
	return storedResources, err
}

// Record a resources collection for further listing
func recordResourcesCollection(namespace, kind string) (err error) {
	collectionsCollection := forgeResCollection(namespace, "collections")
	resCollection := forgeResCollection(namespace, kind)
	log.Printf("Recording resource collection in ns %s: %s ...\n", namespace, resCollection)
	err = scribble.Write(db, collectionsCollection, kind, resCollection)
	return
}

// List all resource collections in a namespace
func listResourcesCollections(namespace string) (collections []string, err error) {
	collectionsCollection := forgeResCollection(namespace, "collections")
	allCollections, err := scribble.ReadAllOrEmpty[string](db, collectionsCollection)
	if err != nil {
		return nil, fmt.Errorf("Unable to read resources collections ! Caused by: %w", err)
	}
	for _, value := range allCollections {
		collections = append(collections, value)
	}
	log.Printf("List of resource collection of ns %s: %s ...\n", namespace, collections)
	return
}

func ListNamespaces() (namespaces []string, err error) {
	// Browse all NS
	allNs, err := listResourcesAsMap("", "", "")
	if err != nil {
		return nil, err
	}
	for _, ns := range allNs {
		name, err := serializ.ResolveJsonMap[string](ns, "/metadata/name")
		if err != nil {
			//err = fmt.Errorf("Bad NS format: %s !", ns)
			return nil, err
		}
		namespaces = append(namespaces, name)
	}
	return
}

// List namespace names corresponding to input.
func developNamespaceNames(namespaceIn string) (namespaces []string, err error) {
	if namespaceIn == "" {
		return
	} else if namespaceIn == "all" {
		// Browse all NS
		allNs, err := ListNamespaces()
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, allNs...)
	} else {
		namespaces = append(namespaces, namespaceIn)
	}
	return
}

func listResourcesAsMap(namespace, kind, name string) ([]map[string]any, error) {
	if namespace == "" {
		// Liste all namespaces
		return listResourcesAsMap(meta_ns_collection, "Namespace", "")
	}
	var resources []map[string]any

	namespaces, err := developNamespaceNames(namespace)
	if err != nil {
		return nil, err
	}

	if kind == "" {
		for _, ns := range namespaces {
			// List all resources in namespace
			collections, err := listResourcesCollections(ns)
			if err != nil {
				return nil, err
			}
			for _, collection := range collections {
				allRes, err := scribble.ReadAllOrEmpty[map[string]any](db, collection)
				if err != nil {
					return nil, fmt.Errorf("Unable to read all resources in namespace: %s ! Caused by: %w", ns, err)
				}
				resources = append(resources, allRes...)
			}
		}
		return resources, nil
	}

	for _, ns := range namespaces {
		kindCollection := forgeResCollection(ns, kind)
		records, err := scribble.ReadAllOrEmpty[map[string]any](db, kindCollection)
		if err != nil {
			return nil, fmt.Errorf("Unable to read resources of kind: %s in namespace: %s ! Caused by: %w", kind, ns, err)
		}
		resources = append(resources, records...)
	}

	if name == "" {
		// List all resources in namespace of kind
		return resources, nil
	}

	// Return one resource in namespace of kind with name
	var mappingError error
	filteredResources := collections.Filter(resources, func(i map[string]any) bool {
		if metadata, ok := i["metadata"].(map[string]any); ok {
			return metadata["name"] == name
		}
		mappingError = fmt.Errorf("Bad metadata in resource: %s", i)
		return false
	})

	return filteredResources, mappingError
}

func consolidateMetadata(namespace, kind, name, jsonIn string) (string, string, string) {
	var k, n, ns string
	if jsonIn != "" {
		// Read namespace from jsonIn (swallow error because optional)
		ns, _ = serializ.ResolveJsonString[string](jsonIn, "/metadata/namespace")
		// Read kind from jsonIn (swallow error because optional)
		k, _ = serializ.ResolveJsonString[string](jsonIn, "/kind")
		// Read name from jsonIn (swallow error because optional)
		n, _ = serializ.ResolveJsonString[string](jsonIn, "/metadata/name")
	}
	if ns == "" {
		ns = namespace
	}
	if k == "" {
		k = kind
	}
	if n == "" {
		n = name
	}
	/*
		if k == "" || kind == "Namespace" {
			if name == "" {
				name = namespace
			}
			kind = "Namespace"
		}
	*/
	return ns, k, n
}

func completeJsonInput(namespace, kind, name, jsonIn string) (map[string]any, error) {
	resourceTree := forgeResource(namespace, kind, name)
	log.Printf("completing JsonInput: %s with forged %s", jsonIn, resourceTree)
	if jsonIn != "" {
		err := yaml.Unmarshal([]byte(jsonIn), &resourceTree)
		if err != nil {
			return nil, err
		}
	}
	var err error
	resourceTree, err = serializ.PatcherMap(resourceTree).
		Default("/kind", kind).
		Default("/metadata/name", name).
		Default("/metadata/namespace", namespace).
		Test("/kind", "Namespace").SwallowError().Then(serializ.OpRemove("/metadata/namespace")).
		ResolveMap()
	if err != nil {
		log.Printf("Found error: %s", err)
		return nil, err
	}
	log.Printf("completed Json: %s", resourceTree)
	return resourceTree, nil
}

// Get resources list (get json description of resources)
func Get(namespace, kind, name, jsonIn string) ([]map[string]any, error) {
	namespace, kind, name = consolidateMetadata(namespace, kind, name, jsonIn)
	mappedResources, err := listResourcesAsMap(namespace, kind, name)
	return mappedResources, err
}

// Create resources (do not overwrite nor update)
func Post(namespace, kind, name, jsonIn string) (resources []map[string]any, err error) {
	namespace, kind, name = consolidateMetadata(namespace, kind, name, jsonIn)
	if kind == "" || kind == "Namespace" {
		// Special case, we need to verify if NS exists
		allNs, err := ListNamespaces()
		if err != nil {
			return nil, err
		}
		if collections.Contains(&allNs, namespace) {
			// namespace already exists
			err = fmt.Errorf("Namespace %s already exists !", namespace)
			return nil, err
		}
	} else {
		// Verify if resource already exists
		out, err := Get(namespace, kind, name, jsonIn)
		if err != nil {
			return nil, err
		}
		if len(out) > 0 {
			// Should not overwrite a resource in Post
			namespace, kind, name = consolidateMetadata(namespace, kind, name, jsonIn)
			err = fmt.Errorf("Resource %s/%s already exists in namespace: %s !", kind, name, namespace)
			return nil, err
		}
	}
	return Put(namespace, kind, name, jsonIn)
}

// Create or Update resources
func Put(namespace, kind, name, jsonIn string) (resources []map[string]any, err error) {
	resourceTree, err := completeJsonInput(namespace, kind, name, jsonIn)
	if err != nil {
		return nil, err
	}
	log.Printf("PUT resourceTree: %s", resourceTree)
	// TODO: update tree with namespace kind and name ?
	storedResources, err := storeResource(namespace, resourceTree)
	return storedResources, err
}

// Update parts of resources
func Patch(namespace, kind, name, jsonText string) (resources []map[string]any, err error) {
	return nil, fmt.Errorf("Not implemented yet !")
}

// Delete resources
func Delete(namespace, kind, name, jsonText string) (resources []map[string]any, err error) {
	return nil, fmt.Errorf("Not implemented yet !")
}
