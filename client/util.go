package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

func unmarshalJSONFile(filePath, fileName string, v interface{}) error {
	cfg, err := ioutil.ReadFile(path.Join(filePath, fileName))
	if err != nil || cfg == nil {
		return err
	}

	err = json.Unmarshal(cfg, &v)
	if err != nil {
		return err
	}

	return nil
}
