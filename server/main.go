package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Package struct {
	Name       string
	Version    Version
	Remote     string
	CommitHash string
}

type Version struct {
	Major int
	Minor int
	Patch int
}

// Set will take a string representation of versions and convert it to
// a Version struct. It fills in empty fields with the value of 0
func (v *Version) Set(version string) *Version {
	segments := strings.Split(version, ".")
	if len(segments) == 0 {
		return &Version{0, 0, 0}
	}

	// Convert our string to integers for storage
	var vers []int
	for _, seg := range segments {
		i, err := strconv.Atoi(seg)
		if err != nil {
			vers = append(vers, 0)
			continue
		}
		vers = append(vers, i)
	}

	// Fill in remaining version fields with 0 (if necessary)
	for i := len(vers); i < 4; i++ {
		vers = append(vers, 0)
	}

	return &Version{vers[0], vers[1], vers[2]}
}

var packages = make(map[string][]*Package)

func main() {
	http.HandleFunc("/register", RegisterPkg)
	http.HandleFunc("/get", GetPkg)
	err := http.ListenAndServe(":13379", nil)
	if err != nil {
		log.Fatal("ListenAndServ: ", err)
	}

	log.Printf("Listening on 1337...")
}

func writeResponseJSON(res http.ResponseWriter, json []byte, statusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(json)
}

func validatePkg(params url.Values) (*Package, error) {
	var (
		name    string
		version Version
		remote  string
		hash    string
	)

	if name = params.Get("pkgname"); name == "" {
		return nil, fmt.Errorf("no package name provided in request")
	}

	version = Version{}
	version.Set(params.Get("version"))

	if remote = params.Get("remote"); remote == "" {
		return nil, fmt.Errorf("no remote name specified")
	}

	if hash = params.Get("hash"); hash == "" {
		return nil, fmt.Errorf("no hash provided")
	}

	return &Package{name, version, remote, hash}, nil
}

type ResponseMsg struct {
	Msg      string
	ErrorMsg string
}

func RegisterPkg(res http.ResponseWriter, req *http.Request) {
	pkg, err := validatePkg(req.URL.Query())
	if err != nil {
		json, err := json.Marshal(ResponseMsg{"Could not register package", err.Error()})
		if err != nil {
			log.Printf("Could not unmarshal package properly")
			http.NotFound(res, req)
			return
		}

		writeResponseJSON(res, json, 406)
		return
	}

	json, err := json.Marshal(ResponseMsg{"Great success, package added to repo!", ""})
	if err != nil {
		log.Printf("Could not unmarshal package properly")
		http.NotFound(res, req)
		return
	}

	packages[pkg.Name] = append(packages[pkg.Name], pkg)
	writeResponseJSON(res, json, 200)
}

func GetPkg(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	var name string

	if name = params.Get("pkgname"); name == "" {
		http.NotFound(res, req)
		return
	}

	pkgs := packages[name]

	if len(pkgs) >= 1 {
		v := Version{}
		v.Set(params.Get("version"))

		for _, p := range pkgs {
			if p.Version.Major == v.Major &&
				p.Version.Minor == v.Minor {

				json, err := json.Marshal(p)
				if err != nil {
					log.Printf("Could not unmarshal package properly")
					http.NotFound(res, req)
					return
				}

				writeResponseJSON(res, json, 200)
			}
		}
	} else {
		http.NotFound(res, req)
		return
	}
}
