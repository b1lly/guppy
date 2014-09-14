package main

import (
	"fmt"
	"os"

	"github.com/b1lly/guppy"
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
	// Iterate through deps
	// Convert dep into request
	// Fetch package data from registry
	// Unmarshal reseponse into a guppy meta package
	// Git clone repo from meta guppy meta package (target=GuppyConfig.Cwd + GuppyConfig.Directory)
	// Do next

	if len(args) == 1 {
		// Fetch single package
		// Check for "-save" flag to add to package dep
		return nil
	}

	for _, deps := range p.Deps {
		fmt.Println(deps)
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

func remoteURL(pkgName string, version string) {
	// Create URL for Deps (e.g. common#1.0.0 => git@x.xx.x.x/storm-common)
}
