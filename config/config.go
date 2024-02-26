package config

import (
	"encoding/json"
	"io"
	"os"
)

type configFile struct {
	Port string `json:"port"`
}

var ConfigFile configFile

func ParseConfig(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &ConfigFile)
	if err != nil {
		return err
	}

	return nil
}
