package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/roland-burke/rollogger"
)

var confPath = "conf.json"
var apiKey = "willbeoverwritten"

const API_KEY_MIN_LENGTH = 12

var conn *pgx.Conn
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
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	logger.Debug(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Unable to connect to database: %s", err.Error())
		os.Exit(1)
	}
}

func readConfig() Configuration {
	file, err := os.Open(confPath)
	defer file.Close()

	if err != nil {
		logger.Error("Cannot open config file from '%s': %s", confPath, err)
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		logger.Error("Cannot decode config: ", err)
	}
	return configuration
}
