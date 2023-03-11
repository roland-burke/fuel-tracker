package config

import (
	"os"
	"testing"

	"github.com/roland-burke/fuel-tracker/internal/model"
	"github.com/roland-burke/rollogger"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger = rollogger.Init(rollogger.INFO_LEVEL, true, true)
}
func TestConfigOverridingConfig(t *testing.T) {
	assert := assert.New(t)

	config := model.Configuration{
		Port:        1234,
		UrlPrefix:   "prefix",
		ApiKey:      "apikey",
		Description: "desc",
	}

	assert.Equal(config.Port, 1234)
	assert.Equal(config.UrlPrefix, "prefix")
	assert.Equal(config.ApiKey, "apikey")
	assert.Equal(config.Description, "desc")

	os.Setenv("FT_DESCRIPTION", "new_description")
	os.Setenv("FT_API-KEY", "new_apikey")
	os.Setenv("FT_PORT", "9000")
	os.Setenv("FT_URL-PREFIX", "new_prefix")

	updateConfigFromEnvironment(&config)

	assert.Equal(config.Port, 9000)
	assert.Equal(config.ApiKey, "new_apikey")
	assert.Equal(config.Description, "new_description")
	assert.Equal(config.UrlPrefix, "new_prefix")
	assert.Contains(Logger.GetLastLog(), "Set config.urlPrefifx from ENV: 'new_prefix'")
}
