package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	expectedConf := Configuration{
		Description: "testing",
		ApiKey:      "apikey",
		Port:        9008,
		UrlPrefix:   "/fuel-tracker",
	}

	var conf = readConfig()

	if conf != expectedConf {
		t.Errorf("Expected Config was wrong, got: %s, want: %s.", convertJsonObjectToString(conf), convertJsonObjectToString(expectedConf))
	}
}
