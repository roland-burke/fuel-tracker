package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roland-burke/rollogger"
)

var confPath = "conf.json"
var apiKey = "willbeoverwritten"

const API_KEY_MIN_LENGTH = 12

var conn *pgxpool.Pool
var logger *rollogger.Log

func printConfig(conf Configuration) {
	// Only show first 5 Characters of api key
	conf.ApiKey = conf.ApiKey[:len(conf.ApiKey)-(len(conf.ApiKey)-5)] + "..."
	logger.Info("Configuration loaded:")
	logger.InfoObj(conf)
}

func main() {
	logger = rollogger.Init(rollogger.INFO_LEVEL, true, true)
	var config = readConfig()
	updateConfigFromEnvironment(config)
	apiKey = config.ApiKey
	if apiKey == "willbeoverwritten" || apiKey == "CHANGEME" {
		logger.Warn("Invalid Apikey: %s\nEither it wasn't changed or something went wrong!\n", apiKey)
		return
	} else if len(apiKey) < API_KEY_MIN_LENGTH {
		logger.Warn("Apikey '%s' too short, must be at least %d characters long!\n", apiKey, API_KEY_MIN_LENGTH)
		return
	}
	var port = config.Port
	var urlPrefix = config.UrlPrefix
	printConfig(config)
	initDb()
	fmt.Printf("=======================================================\n\n")
	startServer(port, urlPrefix)
}

func initDb() {
	var err error

	conn, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	logger.Debug(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Unable to connect to database: %s", err.Error())
		os.Exit(1)
	}
}

func updateConfigFromEnvironment(config Configuration) {
	var desc = os.Getenv("FT_DESCRIPTION")
	if desc != "" {
		config.Description = desc
		logger.Info("Set config.description from ENV: '%s'", desc)
	}
	var apiKey = os.Getenv("FT_API-KEY")
	if apiKey != "" {
		config.ApiKey = apiKey
		logger.Info("Set config.apikey from ENV: '%s'", apiKey[0:5]+"...")
	}
	var port = os.Getenv("FT_PORT")
	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			logger.Warn("Invalid port from ENV: %s", portInt)
		} else {
			logger.Info("Set config.port from ENV: %d", portInt)
			config.Port = portInt
		}
	}
	var urlPrefix = os.Getenv("FT_URL-PREFIX")
	if urlPrefix != "" {
		config.UrlPrefix = urlPrefix
		logger.Info("Set config.urlPrefifx from ENV: '%s'", urlPrefix)
	}

}

func readConfig() Configuration {
	file, err := os.Open(confPath)
	defer file.Close()

	if err != nil {
		logger.Error("Cannot open config file from '%s': %s", confPath, err.Error())
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		logger.Error("Cannot decode config: %s", err.Error())
	}
	return configuration
}
