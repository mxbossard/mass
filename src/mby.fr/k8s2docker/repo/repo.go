package repo

import (
	scribble "github.com/nanobox-io/golang-scribble"
)

var (
	db *scribble.Driver
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

func resCollection(namespace, kind) string {
	return fmt.Sprintf("%s___%s", namespace, kind)
}

func readAllMaps(collection string) []map[string]any, error {
	items := []map[string]any
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

func listNamespaces(namespace string) (namespaces []string, err error) {
	if namespace == "all" {
		// Browse all NS
		allNs, err := readAllMaps("Namespace")
		if err != nil {
			return nil, err
		}
		for _, ns := range allNs {
			namespaces = append(namespaces, ns["ObjectMeta"]["Name"])
		}
	} else {
		namespaces = append(namespaces, namespace)
	}
	return
}

func Get(namespace, kind, name, jsonIn string) (jsonOut string, err error) {
	if namespace == "" {
		// List all NS
		records := db.ReadAll("Namespace")
		var allNs string[]
		for _, record := range records {
			allNs = append(allNs, record)
		}
		jsonOut = "[" + strings.Join(records, ",") + "]"
		return
	} 

	namespaces, err = listNamespaces(namespace)
	if err != nil {
		return
	}


	if name == "" {
		// Get all resources
		resources, err := db.ReadAll(kind)
		if err != nil {
			return
		}
	}

}

func Post(namespace, kind, name, jsonIn string, labelFilter) (jsonOut string, err error) {

}

func Put(namespace, kind, name, jsonIn string, labelFilter) (jsonOut string, err error) {

}

func Patch(namespace, kind, name, jsonText string) (jsonOut string, err error) {
	
}

func Delete(namespace, kind, name, jsonText string) (jsonOut string, err error) {
	
}
