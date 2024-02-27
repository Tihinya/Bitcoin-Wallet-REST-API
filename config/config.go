package config

import (
	"encoding/json"
	"fmt"
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
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return nil
}
