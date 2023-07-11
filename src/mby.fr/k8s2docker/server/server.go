package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"mby.fr/k8s2docker/repo"
)

var (
	serveCoreResourcesRootPath = "/api/v1/namespaces"
	// pattern /api/v1/namespaces/NAMESPACE_NAME/RESOURCE_KINDs/RESOURCE_NAME
	// /api/v1/namespaces(?:/([a-z0-9][a-z0-9-]*[a-z0-9])(?:/([a-z]+))(?:/([a-z0-9][a-z0-9-.]*[a-z0-9]))?)?
	/*
		/api/v1/namespaces
		/api/v1/namespaces/
		/api/v1/namespaces/default
		/api/v1/namespaces/default/
		/api/v1/namespaces/default/pods
		/api/v1/namespaces/default/pods.foo
		/api/v1/namespaces/default/pods/
		/api/v1/namespaces/default/pods/name.foo
		/api/v1/namespaces/default/pods/name
		/api/v1/namespaces/default/pods/name/
	*/
	serveCoreResourcesPattern = regexp.MustCompile("^" + serveCoreResourcesRootPath + "(?:/(?P<namespace>[^/]+)(?:/(?P<kind>[^/]+)?(?:/(?P<name>[^/]+)?)?)?)?/?$")

	namespaceNamePattern = regexp.MustCompile("^[a-z0-9][a-z0-9-.]*[a-z0-9]$")
	resourceKindPattern  = regexp.MustCompile("^([a-z]+)s$")
	resourceNamePattern  = regexp.MustCompile("^[a-z0-9][a-z0-9-.]*[a-z0-9]$")
	ContentTypeHeader    = "Content-Type"
	JsonContentType      = "application/json"
)

func Start() (err error) {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil, // use DefaultServeMux by default
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer server.Close()

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc(serveCoreResourcesRootPath, coreResourcesHandler)
	http.HandleFunc(serveCoreResourcesRootPath+"/", coreResourcesHandler)

	err = server.ListenAndServe()
	log.Printf("Server error: %s", err)
	server = nil

	return
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	serveErrors(w, r, 404, fmt.Sprintf("Not supported call to Path: %s", r.URL.Path))
}

func coreResourcesHandler(w http.ResponseWriter, r *http.Request) {
	// Exemple d'URL à matcher /apis/apps/v1/namespaces/default/pods/<name>
	loc := serveCoreResourcesPattern.FindStringSubmatch(r.URL.Path)
	if loc != nil {
		namespace := loc[1]
		kind := loc[2]
		name := loc[3]
		var err error
		var errors []string

		if namespace != "" {
			err = assertNamespace(namespace)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
		if kind != "" {
			kind, err = assertKind(kind)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
		if name != "" {
			err = assertResourceName(name)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}

		if len(errors) > 0 {
			serveErrors(w, r, 400, errors...)
			return
		}

		if kind == "" {
			serveNamespacesHandler(w, r, namespace)
			return
		} else if name == "" {
			serveKindsHandler(w, r, namespace, kind)
			return
		} else {
			serveCoreResourcesHandler(w, r, namespace, kind, name)
			return
		}
	}
	//log.Printf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	serveErrors(w, r, 400, fmt.Sprintf("Not supported core resources URL: %s", r.URL.Path))
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, object any) {
	w.Header().Set(ContentTypeHeader, JsonContentType)
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Printf("Error Marshalling object to json: %s !", err)
	}
	w.WriteHeader(statusCode)
	w.Write(jsonBytes)

}

func readJsonRequest(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func serveErrors(w http.ResponseWriter, r *http.Request, statusCode int, messages ...string) {
	err := ServerError{
		Message: strings.Join(messages, "\n"),
		Path:    r.URL.Path,
	}
	writeJsonResponse(w, statusCode, err)
}

func serveMessages(w http.ResponseWriter, r *http.Request, messages ...string) {
	for _, msg := range messages {
		if msg != "" {
			w.Write([]byte(msg))
			w.Write([]byte("\n"))
		}
	}
}

func serveNamespacesHandler(w http.ResponseWriter, r *http.Request, namespace string) {
	json, err := readJsonRequest(r)
	if err != nil {
		serveErrors(w, r, 400, err.Error())
		return
	}
	jsonOut, err := repo.Interract(namespace, "", "", json, r.Method)
	if err != nil {
		serveErrors(w, r, 500, err.Error())
		return
	}

	writeJsonResponse(w, 200, jsonOut)
}

func serveKindsHandler(w http.ResponseWriter, r *http.Request, namespace, kind string) {
	json, err := readJsonRequest(r)
	if err != nil {
		serveErrors(w, r, 400, err.Error())
		return
	}
	jsonOut, err := repo.Interract(namespace, kind, "", json, r.Method)
	if err != nil {
		serveErrors(w, r, 500, err.Error())
		return
	}

	writeJsonResponse(w, 200, jsonOut)
}

func serveCoreResourcesHandler(w http.ResponseWriter, r *http.Request, namespace, kind, name string) {
	json, err := readJsonRequest(r)
	if err != nil {
		serveErrors(w, r, 400, err.Error())
		return
	}
	jsonOut, err := repo.Interract(namespace, kind, name, json, r.Method)
	if err != nil {
		serveErrors(w, r, 500, err.Error())
		return
	}

	writeJsonResponse(w, 200, jsonOut)
}

func assertNamespace(namespace string) (err error) {
	if !namespaceNamePattern.MatchString(namespace) {
		err = fmt.Errorf("Bad namespace format: [%s] !", namespace)
	}
	return
}

func assertKind(kind string) (formattedKind string, err error) {
	if !resourceKindPattern.MatchString(kind) {
		err = fmt.Errorf("Bad resource kind format: [%s] !", kind)
	} else {
		// rewrite kind: pods => Pod
		formattedKind = strings.ToUpper(kind[0:1]) + kind[1:len(kind)-1]
	}
	return
}

func assertResourceName(name string) (err error) {
	if !resourceNamePattern.MatchString(name) {
		err = fmt.Errorf("Bad resource name format: [%s] !", name)
	}
	return
}
