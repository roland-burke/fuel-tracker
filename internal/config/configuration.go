package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/roland-burke/fuel-tracker/internal/model"
	"github.com/roland-burke/rollogger"
)

const API_KEY_MIN_LENGTH = 12

var Logger *rollogger.Log
var ApiKey = "willbeoverwritten"
var confPath = "conf.json"

func InitConfig() (int, string) {
	var config = readConfig()
	updateConfigFromEnvironment(&config)
	ApiKey = config.ApiKey
	if ApiKey == "willbeoverwritten" || ApiKey == "CHANGEME" {
		Logger.Warn("Invalid Apikey: %s\nEither it wasn't changed or something went wrong!\n", ApiKey)
		os.Exit(3)
	} else if len(ApiKey) < API_KEY_MIN_LENGTH {
		Logger.Warn("Apikey '%s' too short, must be at least %d characters long!\n", ApiKey, API_KEY_MIN_LENGTH)
		os.Exit(4)
	}
	printConfig(config)
	return config.Port, config.UrlPrefix
}

func updateConfigFromEnvironment(config *model.Configuration) {
	var desc = os.Getenv("FT_DESCRIPTION")
	if desc != "" {
		config.Description = desc
		Logger.Info("Set config.description from ENV: '%s'", desc)
	}
	var apiKey = os.Getenv("FT_API-KEY")
	if apiKey != "" {
		config.ApiKey = apiKey
		Logger.Info("Set config.apikey from ENV: '%s'", apiKey[0:5]+"...")
	}
	var port = os.Getenv("FT_PORT")
	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			Logger.Warn("Invalid port from ENV: %s", port)
		} else {
			Logger.Info("Set config.port from ENV: %d", portInt)
			config.Port = portInt
		}
	}
	var urlPrefix = os.Getenv("FT_URL-PREFIX")
	if urlPrefix != "" {
		config.UrlPrefix = urlPrefix
		Logger.Info("Set config.urlPrefifx from ENV: '%s'", urlPrefix)
	}
}

func readConfig() model.Configuration {
	file, err := os.Open(confPath)
	defer file.Close()

	if err != nil {
		Logger.Error("Cannot open config file from '%s': %s", confPath, err.Error())
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	configuration := model.Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		Logger.Error("Cannot decode config: %s", err.Error())
		os.Exit(2)
	}
	return configuration
}

func printConfig(conf model.Configuration) {
	// Only show first 5 Characters of api key
	conf.ApiKey = conf.ApiKey[:len(conf.ApiKey)-(len(conf.ApiKey)-5)] + "..."
	Logger.Info("Configuration loaded:")
	Logger.InfoObj(conf)
}
