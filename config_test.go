package config

import (
	"testing"
)

func TestNewFromYAML(t *testing.T) {
	_, err := NewFromYAML("./t/config.yaml")
	if err != nil {
		t.Errorf("NewFromYAML, %s", err)
	}
}

func TestGet(t *testing.T) {
	// read config file
	cfg, _ := NewFromYAML("./t/config.yaml")

	// get key
	key := "key1"
	expectedValue := "value for key1"
	gotValue, _ := cfg.Get(key).(string)

	// check
	if gotValue != expectedValue {
		t.Errorf("Get(%q), got: %q, want: %q", key, gotValue, expectedValue)
	}
}
