package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/b1lly/guppy"
)

var registry *Registry

func main() {
	initDB()

	registry = NewRegistry()
	err := registry.Load()
	if err != nil {
		log.Fatal("Could not load the registry properly from the database. ", err)
	}

	http.HandleFunc("/register", RegisterPkg)
	http.HandleFunc("/search", SearchPkg)

	log.Printf("Listening on 13379...")
	log.Fatal("ListenAndServ: ", http.ListenAndServe(":13379", nil))
}

func writeResponseJSON(res http.ResponseWriter, req *http.Request, msg interface{}, statusCode int) {
	json, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Could not unmarshal package properly")
		http.NotFound(res, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(json)
}

// RegisterPkg will register a new package with the remote server based on the
// required request params
func RegisterPkg(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	pkg, err := guppy.NewPackage(params.Get("pkgname"), params.Get("version"), params.Get("remote"), params.Get("hash"))
	if err != nil {
		writeResponseJSON(res, req, err, 406)
		return
	}

	err = registry.Add(pkg)
	if err != nil {
		writeResponseJSON(res, req, err, 406)
		return
	}

	writeResponseJSON(res, req, pkg, 200)
}

// SearchPkg will look up a specified package and return it back to the request as JSON
func SearchPkg(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var name string

	if name = params.Get("pkgname"); name == "" {
		writeResponseJSON(res, req, guppy.PackageError{"Please provide a valid package name"}, 406)
		return
	}

	version := guppy.NewVersion(params.Get("version"))
	pkg := registry.PackageByNameAndVersion(name, version)
	if pkg != nil {
		writeResponseJSON(res, req, pkg, 201)
		return
	}

	writeResponseJSON(res, req, guppy.PackageError{fmt.Sprintf("Package `%s` version `%s` does not exist.", name, version.String())}, 404)
}
