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

type GuppyConfig struct {
	Cwd             string
	Directory       string
	RegistryPrivate string
	RegistryPublic  string
	Color           bool
}

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
		err = gc.Load(user.HomeDir)

		if err != nil {
			glog.Println(fmt.Sprintf("Failed to load guppy root config `%v`", err))
		}
	}

	if cwd != user.HomeDir {
		err = gc.Load(cwd)
		if err != nil {
			glog.Println(fmt.Sprintf("Failed to load guppy project config `%v`", err))
		}
	}

	return gc
}

func (gc *GuppyConfig) Load(cfgPath string) error {
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
