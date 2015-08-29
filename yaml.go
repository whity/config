package config

import (
	"fmt"
	"github.com/whity/file-utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// NewFromYAML read configuration from a yaml file
func NewFromYAML(filename string) (*Config, error) {
	// check if file exists
	exists := fileUtils.FileExists(filename)
	if !exists {
		return New(make(dict)), fmt.Errorf("%s: no such file", filename)
	}

	// read yaml file
	dct := make(map[string]interface{})
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return New(dct), err
	}

	//
	err = yaml.Unmarshal(yamlFile, &dct)
	if err != nil {
		return New(dct), err
	}

	return New(dct), nil
}
