package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/b1lly/guppy"
	_ "github.com/yext/glog"
)

const (
	projectFile = "gup.json"
)

type Project struct {
	guppy.Package
	Version string
	Deps    []string
	Private bool
}

func NewProject() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var project *Project
	err = unmarshalJSONFile(cwd, projectFile, &project)
	return project, err
}

func (p *Project) Install(args []string) error {
	if len(args) == 1 {
		// Convert dep into request
		// Fetch single package
		// Check for "-save" flag to add to package dep
		return nil
	}

	// Install the main package
	for _, dep := range p.Deps {
		var segs = make([]string, 3)
		segs = strings.Split(dep, "#")

		if segs[1] == "private" {
			segs[1] = ""
			segs[2] = "private"
		}

		isPrivate, err := strconv.ParseBool(segs[2])
		if err != nil {
			isPrivate = false
		}

		resp, err := http.Get(searchRegistryUrl(segs[0], segs[1], isPrivate))
		if err != nil {
			glog.Info("Could not fetch package information:", err)
			continue
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.Info("Failed to get package information: ", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != 200 {
			glog.Info("Failed to fetch package", string(bytes))
			continue
		}

		pkg, err := guppy.NewPackageFromJSON(bytes)
		if err != nil {
			glog.Info(err)
			continue
		}

		glog.Info(pkg)

		// Unmarshal reseponse into a guppy meta package
		// Git clone repo from meta guppy meta package (target=GuppyConfig.Cwd + GuppyConfig.Directory)
		// Do next
	}

	return nil
}

func (p *Project) Register(args []string) error {
	// - convert settings into request
	// - confirm remote pkg repo exists
	// - send register request to gup registry server
	return nil
}

func (p *Project) Save() {
	// Write state to disc
}

func searchRegistryUrl(pkgName string, version string, private bool) string {
	var registry string
	switch private {
	case true:
		registry = guppyCfg.RegistryPrivate
	case false:
		registry = guppyCfg.RegistryPublic
	}

	params := url.Values{}
	params.Add("pkgname", pkgName)
	params.Add("version", version)

	// e.g. common#1.0.0 => localhost/storm-common
	return fmt.Sprintf("http://%s/search?%s", registry, params.Encode())
}
