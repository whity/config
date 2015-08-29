package config

import (
	"fmt"
	"reflect"
	"strings"
)

type dict map[string]interface{}

func (dct dict) rawData() map[string]interface{} {
	return rawData(dct).(map[string]interface{})
}

func (dct dict) clone() dict {
	return clone(dct).(dict)
}

type list []interface{}

func (lst list) rawData() []interface{} {
	return rawData(lst).([]interface{})
}

func (lst list) clone() list {
	return clone(lst).(list)
}

// Config config type
type Config struct {
	data dict
}

// New build a config passing map[interface{}]interface{}
func New(data dict) *Config {
	// clone data
	clonedData := data.clone()

	return &Config{data: clonedData}
}

// Get value from the config
func (config *Config) Get(key string) interface{} {
	// split key by '.'
	keys := strings.Split(key, ".")

	var result interface{}
	lastIdx := len(keys) - 1
	data := config.data
	for idx, k := range keys {
		tmpData := data[k]

		// on last idx or nil, just return
		if tmpData == nil || idx >= lastIdx {
			result = tmpData
			break
		}

		// isn't the last idx ( key level ) and isn't a Dict, return nil
		typeStr := reflect.TypeOf(tmpData).String()
		if typeStr != "config.dict" {
			result = nil
			break
		}

		data = tmpData.(dict)
	}

	if result != nil {
		result = rawData(result)
	}

	return result
}

// Merge with other config object
func (config *Config) Merge(other *Config) *Config {
	newData := merge(config.data, other.data)

	return &Config{data: newData}
}

func rawData(data interface{}) interface{} {
	dataType := reflect.TypeOf(data).String()

	if dataType == "config.dict" {
		newData := map[string]interface{}{}
		for key, value := range data.(dict) {
			newData[key] = rawData(value)
		}

		return newData
	}

	if dataType == "config.list" {
		newData := []interface{}{}
		for _, value := range data.(list) {
			newData = append(newData, rawData(value))
		}

		return newData
	}

	return data
}

func clone(data interface{}) interface{} {
	dataType := reflect.TypeOf(data).Kind()

	if dataType == reflect.Map {
		newData := dict{}
		tmpMap := reflect.ValueOf(data)
		mapKeys := tmpMap.MapKeys()
		for _, key := range mapKeys {
			value := tmpMap.MapIndex(key)
			newData[fmt.Sprintf("%v", key.Interface())] = clone(value.Interface())
		}

		return newData
	}

	if dataType == reflect.Slice {
		newData := list{}
		slice := reflect.ValueOf(data)
		for idx := 0; idx < slice.Len(); idx++ {
			value := slice.Index(idx)
			newData = append(newData, clone(value.Interface()))
		}

		return newData
	}

	return data
}

func merge(left dict, right dict) dict {
	merged := left.clone()

	for key, value := range right {
		/*
			check if key exists on the left
				- if not: just set it
				- otherwise:
					- if both are Dict, call 'merge' on them
					- if not, just set it
		*/

		// if key doesn't exists on left, just add it
		currentValue, exists := left[key]
		if !exists ||
			reflect.TypeOf(currentValue).String() != "config.dict" ||
			reflect.TypeOf(value).String() != "config.dict" {
			merged[key] = value
			continue
		}

		merged[key] = merge(currentValue.(dict), value.(dict))
	}

	return merged
}
