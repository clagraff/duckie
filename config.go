package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var CFG_FILE_NAME string = ".ducky"

type Config struct {
	Path   string `json:"path"`
	Dir    string `json:"dir"`
	AddTxn bool   `json:"add_transactions"`
}

func ReadCfg() (Config, error) {
	cfg := Config{}

	currDir, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	path := filepath.Join(currDir, CFG_FILE_NAME)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(contents, &cfg)
	return cfg, err
}

func WriteCfg(cfg Config) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(currDir, CFG_FILE_NAME)

	content, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, content, 0777)
	return err
}
