package src

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Parser interface {
	Parse(file string) (map[string]EnvConfig, error)
}

type YamlParser struct{}

func (y YamlParser) Parse(file string) (map[string]EnvConfig, error) {
	configMap := map[string]EnvConfig{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return configMap, err
	}

	if err := yaml.Unmarshal(b, &configMap); err != nil {
		return configMap, err
	}

	return configMap, nil
}

type JsonParser struct{}

func (j JsonParser) Parse(file string) (map[string]EnvConfig, error) {
	configMap := map[string]EnvConfig{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return configMap, err
	}

	if err := json.Unmarshal(b, &configMap); err != nil {
		return configMap, err
	}

	return configMap, nil
}
