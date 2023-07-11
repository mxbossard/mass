package repo

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	//corev1 "k8s.io/api/core/v1"

	"mby.fr/utils/collections"

	scribble "github.com/nanobox-io/golang-scribble"
)

var (
	db    *scribble.Driver
	dbDir = "mydb"
)

func init() {
	initDb()
}

func initDb() {
	if db != nil {
		return
	}
	var err error
	db, err = scribble.New(dbDir, nil)
	if err != nil {
		log.Fatalf("Error initializing DB: %w", err)
	}
}

func Interract(namespace, kind, name, jsonIn, method string) (jsonOut string, err error) {
	if kind != "Pod" && kind != "Namespace" {
		// Not supported resource kind
		err = fmt.Errorf("Not supported kind: %s !", kind)
	}

	if method == "GET" {
		return Get(namespace, kind, name, jsonIn)
	} else if method == "POST" {
		return Post(namespace, kind, name, jsonIn)
	} else if method == "PUT" {
		return Put(namespace, kind, name, jsonIn)
	} else if method == "PATCH" {
		return Patch(namespace, kind, name, jsonIn)
	} else if method == "DELETE" {
		return Delete(namespace, kind, name, jsonIn)
	} else {
		err = fmt.Errorf("Not supported method: %s !", method)
	}
	return
}

func resCollection(namespace, kind string) string {
	return fmt.Sprintf("%s___%s", namespace, kind)
}

func readAllMaps(collection string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	records, err := db.ReadAll(collection)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		var i map[string]any
		if err := json.Unmarshal([]byte(r), &i); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func readAny[T any](collection string) ([]T, error) {
	var items []T
	records, err := db.ReadAll(collection)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		var i T
		if err := json.Unmarshal([]byte(r), &i); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func listNamespaces(namespaceIn string) (namespaces []string, err error) {
	if namespaceIn == "all" {
		// Browse all NS
		allNs, err := readAllMaps("Namespace")
		if err != nil {
			return nil, err
		}
		for _, ns := range allNs {
			metadata, ok := ns["metadata"].(map[string]interface{})
			if !ok {
				err = fmt.Errorf("Bad NS format: %s !", ns)
				return nil, err
			}
			name, ok := metadata["name"].(string)
			if !ok {
				err = fmt.Errorf("Bad NS format: %s !", ns)
				return nil, err
			}
			namespaces = append(namespaces, name)
		}
	} else {
		namespaces = append(namespaces, namespaceIn)
	}
	return
}

func listResourcesCollections(namespace string) (collections []string, err error) {
	collectionsCollection := resCollection(namespace, "collections")
	allCollections, err := readAny[string](collectionsCollection)
	if err != nil {
		return nil, err
	}
	for _, value := range allCollections {
		collections = append(collections, value)
	}
	return
}

func listMappedResources(namespace, kind, name string) ([]map[string]any, error) {
	if namespace == "" {
		// Liste all namespaces
		return listMappedResources("__meta", "Namespace", "")
	}
	if kind == "" {
		// List all resources in namespace
		var resources []map[string]any
		collections, err := listResourcesCollections(namespace)
		if err != nil {
			return nil, err
		}
		for _, collection := range collections {
			allRes, err := readAny[map[string]any](collection)
			if err != nil {
				return nil, err
			}
			resources = append(resources, allRes...)
		}
		return resources, nil
	}

	kindCollection := resCollection(namespace, kind)
	records, err := readAny[map[string]any](kindCollection)
	if err != nil {
		return nil, err
	}

	if name == "" {
		// List all resources in namespace of kind
		return records, nil
	}

	// Return one resource in namespace of kind with name
	var mappingError error
	record := collections.Filter(records, func(i map[string]any) bool {
		if metadata, ok := i["metadata"].(map[string]any); ok {
			return metadata["name"] == name
		}
		mappingError = fmt.Errorf("Bad metadata in resource: %s", i)
		return false
	})

	return record, mappingError
}

func Get(namespace, kind, name, jsonIn string) (string, error) {
	mappedResources, err := listMappedResources(namespace, kind, name)
	if err != nil {
		return "", err
	}
	var mappingError error
	jsonResources := collections.Map(mappedResources, func(i map[string]any) string {
		outBytes, err := json.Marshal(i)
		if err != nil {
			mappingError = err
			return ""
		}
		return string(outBytes)
	})
	if mappingError != nil {
		return "", mappingError
	}
	jsonOut := strings.Join(jsonResources, "\n---\n")
	return jsonOut, nil
}

func Post(namespace, kind, name, jsonIn string) (jsonOut string, err error) {
	return
}

func Put(namespace, kind, name, jsonIn string) (jsonOut string, err error) {
	return
}

func Patch(namespace, kind, name, jsonText string) (jsonOut string, err error) {
	return
}

func Delete(namespace, kind, name, jsonText string) (jsonOut string, err error) {
	return
}
