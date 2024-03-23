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

	"mby.fr/k8s2docker/descriptor"
	"mby.fr/utils/collections"
	"mby.fr/utils/structz"

	"mby.fr/scribble"
)

var (
	db                 *scribble.Driver
	meta_ns_collection = "__meta"
)

func InitDb(dbDirPath string) {
	if db != nil {
		return
	}
	var err error
	db, err = scribble.New(dbDirPath, nil)
	if err != nil {
		log.Fatalf("Error initializing DB at path: [%s] ! Cause by: %s", dbDirPath, err)
	}
}

func validateJsonInput(resourceTree map[string]any, kind, name string) (string, error) {
	validatedTree, err := descriptor.ValidateMappedK8sResource(resourceTree, kind, name)
	if err != nil {
		return "", err
	}
	//log.Printf("Validated tree: %v", validatedTree)

	jsonBytes, err := json.Marshal(validatedTree)
	if err != nil {
		return "", err
	}
	jsonIn := string(jsonBytes)
	return jsonIn, nil
}

func Interract(namespace, kind, name, jsonIn, method string) (jsonOut string, err error) {
	if kind != "Pod" && kind != "Namespace" {
		// Not supported resource kind
		err = fmt.Errorf("Not supported kind: %s !", kind)
	}

	resourceTree, err := completeJsonInput(namespace, kind, name, jsonIn)
	if err != nil {
		return "", err
	}

	log.Printf("Validated input json: %v", jsonIn)
	namespace, kind, name = consolidateMetadata(namespace, kind, name, jsonIn)

	var resources []map[string]any
	if method == "GET" {
		resources, err = Get(namespace, kind, name, jsonIn)
	} else if method == "POST" {
		jsonIn, err = validateJsonInput(resourceTree, kind, name)
		if err != nil {
			return "", err
		}
		resources, err = Post(namespace, kind, name, jsonIn)
	} else if method == "PUT" {
		jsonIn, err = validateJsonInput(resourceTree, kind, name)
		if err != nil {
			return "", err
		}
		resources, err = Put(namespace, kind, name, jsonIn)
	} else if method == "PATCH" {
		jsonIn, err = validateJsonInput(resourceTree, kind, name)
		if err != nil {
			return "", err
		}
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
	if kind == "Namespace" {
		namespace = meta_ns_collection
	}
	return fmt.Sprintf("%s___%s", namespace, kind)
}

func loadResource(namespace, kind, name string) (map[string]any, error) {
	collection := forgeResCollection(namespace, kind)
	return scribble.Read[map[string]any](db, collection, name)
}

func storeResource(namespace string, resourceTree map[string]any) ([]map[string]any, error) {
	var storedResources []map[string]any
	explorer := structz.JsonMapExplorer(resourceTree)
	// Read kind from tree1
	kind, err := structz.Resolve[string](explorer, "/kind")
	if err != nil {
		return nil, err
	}
	// Read name from tree
	name, err := structz.Resolve[string](explorer, "/metadata/name")
	if err != nil {
		return nil, err
	}
	// swallow error namespace is optionnal in tree
	jsonNamespace, _ := structz.Resolve[string](explorer, "/metadata/namespace")
	if jsonNamespace != "" {
		namespace = jsonNamespace
	}

	collection := forgeResCollection(namespace, kind)
	log.Printf("Storing %s %s: %s\n", collection, name, resourceTree)
	err = scribble.Write(db, collection, name, resourceTree)
	if err != nil {
		return storedResources, err
	}

	// Ensure resources collection is referenced
	recordResourcesCollection(namespace, kind)

	if kind != "Namespace" {
		namespaceNames, err := ListNamespaceNames()
		if err != nil {
			return nil, err
		}

		if !collections.Contains(&namespaceNames, namespace) {
			// Namespace does not exists yet create it for simplicity
			log.Printf("Namespace %s does not exists yet. (among: %v)\n", namespace, namespaceNames)
			nsJson := forgeNamespace(namespace)
			res, err := storeResource(meta_ns_collection, nsJson)
			if err != nil {
				return nil, err
			}
			storedResources = append(storedResources, res...)
		}
	}
	storedResources = append(storedResources, resourceTree)

	return storedResources, err
}

func deleteMapResource(namespace string, res map[string]any) (err error) {
	explorer := structz.JsonMapExplorer(res)
	resKind, err := structz.Resolve[string](explorer, "/kind")
	if err != nil {
		return err
	}
	resName, err := structz.Resolve[string](explorer, "/metadata/name")
	if err != nil {
		return err
	}
	log.Printf("Deleting mapres: %s/%s/%s ...", namespace, resKind, resName)
	collection := forgeResCollection(namespace, resKind)
	err = db.Delete(collection, resName)
	if err != nil {
		return err
	}
	return
}

func deleteResource(namespace, kind, name string) (ressources []map[string]any, err error) {
	if namespace == "" {
		err = fmt.Errorf("A namespace must be supplied for Delete operation !")
		return
	}
	log.Printf("Deleting res: %s/%s/%s ...", namespace, kind, name)
	err = assertNamespaceExist(namespace)
	if err != nil {
		return
	}

	selectedResources, err := listResourcesAsMap(namespace, kind, name)
	if err != nil {
		return
	}

	log.Printf("selectedResources: %v", selectedResources)
	for _, res := range selectedResources {
		err = deleteMapResource(namespace, res)
		if err != nil {
			return
		}
		ressources = append(ressources, res)
	}
	if kind == "" {
		// Deleting namespace
		namespaceRes, err := loadResource("", "Namespace", namespace)
		log.Printf("namespace: %v", namespaceRes)
		if err != nil {
			return nil, err
		}
		err = deleteMapResource("", namespaceRes)
		if err != nil {
			return nil, err
		}
		ressources = append(ressources, namespaceRes)
	} else {
		if len(selectedResources) == 0 {
			err = fmt.Errorf("Cannot delete not existing resource of name: [%s] with kind: [%s] in namespace: [%s] !", name, kind, namespace)
		}
	}
	return
}

// Record a resources collection for further listing
func recordResourcesCollection(namespace, kind string) (err error) {
	if kind == "Namespace" {
		// Do not record meta collections
		return
	}
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

func ListNamespaceNames() (namespaces []string, err error) {
	// Browse all NS
	allNs, err := listResourcesAsMap("", "", "")
	if err != nil {
		return nil, err
	}
	for _, ns := range allNs {
		explorer := structz.JsonMapExplorer(ns)
		name, err := structz.Resolve[string](explorer, "/metadata/name")
		if err != nil {
			//err = fmt.Errorf("Bad NS format: %s !", ns)
			return nil, err
		}
		namespaces = append(namespaces, name)
	}
	return
}

func assertNamespaceExist(name string) (err error) {
	namespaceNames, err := ListNamespaceNames()
	if err != nil {
		return err
	}
	if !collections.Contains(&namespaceNames, name) {
		err = fmt.Errorf("Namespace %s does not exists !", name)
	}
	return
}

// List namespace names corresponding to input.
func developNamespaceNames(namespaceIn string) (namespaces []string, err error) {
	if namespaceIn == "" {
		return
	} else if namespaceIn == "all" {
		// Browse all NS
		allNs, err := ListNamespaceNames()
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, allNs...)
	} else {
		namespaces = append(namespaces, namespaceIn)
	}
	return
}

/*
List resources baked by map[string]any
if namespace == "" => List all NS
if kind == "" => List all resources in NS
if name == "" => List all resources of Kind in NS
*/
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
		explorer := structz.JsonStringExplorer(jsonIn)
		// Read namespace from jsonIn (swallow error because optional)
		ns, _ = structz.Resolve[string](explorer, "/metadata/namespace")
		// Read kind from jsonIn (swallow error because optional)
		k, _ = structz.Resolve[string](explorer, "/kind")
		// Read name from jsonIn (swallow error because optional)
		n, _ = structz.Resolve[string](explorer, "/metadata/name")
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
	resourceTree, err = structz.PatcherMap(resourceTree).
		Default("/kind", kind).
		Default("/metadata/name", name).
		Default("/metadata/namespace", namespace).
		Test("/kind", "Namespace").SwallowError().Then(structz.OpRemove("/metadata/namespace")).
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
		allNs, err := ListNamespaceNames()
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
func Patch(namespace, kind, name, jsonIn string) (resources []map[string]any, err error) {
	return nil, fmt.Errorf("repo.Patch() Not implemented yet !")
}

// Delete resources
func Delete(namespace, kind, name, jsonIn string) (resources []map[string]any, err error) {
	namespace, kind, name = consolidateMetadata(namespace, kind, name, jsonIn)
	deletedResources, err := deleteResource(namespace, kind, name)
	return deletedResources, err
}
