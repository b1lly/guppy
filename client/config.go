package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/imdario/mergo"
)

// GuppyConfig is used to tell guppy how to handle packages.
// e.g. What registry to pull from, and what directories to clone them to.
type GuppyConfig struct {
	Cwd             string
	Directory       string
	RegistryPrivate string
	RegistryPublic  string
	Color           bool
}

// NewGuppyConfig will return a config based on a few things;
// It first creates some default settings and merges/overwrites that with the root config (~/.guppy)
// It then looks for a project specific config (current working directory) and will merge/overwrite
// that previous result with the project specific config settings.
func NewGuppyConfig() *GuppyConfig {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	gc := &GuppyConfig{
		cwd,
		"/guppy_modules",
		"localhost:13379",
		"localhost:13379",
		true,
	}

	user, err := user.Current()
	if err == nil {
		err = gc.load(user.HomeDir)

		if err != nil {
			glog.Println(fmt.Sprintf("Failed to load guppy root config `%v`", err))
		}
	}

	if cwd != user.HomeDir {
		err = gc.load(cwd)
		if err != nil {
			glog.Println(fmt.Sprintf("Failed to load guppy project config `%v`", err))
		}
	}

	return gc
}

func (gc *GuppyConfig) load(cfgPath string) error {
	cfg, err := ioutil.ReadFile(path.Join(cfgPath, guppyFile))
	if err != nil || cfg == nil {
		return err
	}

	var fileConfig *GuppyConfig
	err = json.Unmarshal(cfg, &fileConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err = mergo.Merge(gc, *fileConfig); err != nil {
		return err
	}

	return nil
}

func (gc *GuppyConfig) Save() error {

	return nil
}
