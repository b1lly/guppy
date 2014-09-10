package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/b1lly/guppy"
)

var packages = make(map[string][]*guppy.Package)

func main() {
	http.HandleFunc("/register", RegisterPkg)
	http.HandleFunc("/get", GetPkg)
	err := http.ListenAndServe(":13379", nil)
	if err != nil {
		log.Fatal("ListenAndServ: ", err)
	}

	log.Printf("Listening on 13379...")
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

func validatePkg(params url.Values) (*guppy.Package, error) {
	var (
		name    string
		version guppy.Version
		remote  string
		hash    string
	)

	if name = params.Get("pkgname"); name == "" {
		return nil, fmt.Errorf("no package name provided in request")
	}

	version = guppy.Version{}
	version.Set(params.Get("version"))

	if remote = params.Get("remote"); remote == "" {
		return nil, fmt.Errorf("no remote name specified")
	}

	if hash = params.Get("hash"); hash == "" {
		return nil, fmt.Errorf("no hash provided")
	}

	return &guppy.Package{name, version, remote, hash}, nil
}

// ResponseMsg is used to send a structured JSON response with an error
// back to the request
type ResponseMsg struct {
	Msg      string
	ErrorMsg string
}

// RegisterPkg is used to register a new package with the remote server
func RegisterPkg(res http.ResponseWriter, req *http.Request) {
	pkg, err := validatePkg(req.URL.Query())
	if err != nil {
		writeResponseJSON(res, req, ResponseMsg{"Could not register package", err.Error()}, 406)
		return
	}

	packages[pkg.Name] = append(packages[pkg.Name], pkg)
	writeResponseJSON(res, req, pkg, 200)
}

// GetPkg will look up a specified package and return it back to the request as JSON
func GetPkg(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var name string

	if name = params.Get("pkgname"); name == "" ||
		len(packages[name]) < 1 {
		writeResponseJSON(res, req, ResponseMsg{fmt.Sprintf("No packages found called `%v`", name), ""}, 404)
		return
	}

	v := guppy.Version{}
	v.Set(params.Get("version"))

	for _, p := range packages[name] {
		if p.Version.Major == v.Major &&
			p.Version.Minor == v.Minor {
			writeResponseJSON(res, req, p, 200)
			return
		}
	}
}
