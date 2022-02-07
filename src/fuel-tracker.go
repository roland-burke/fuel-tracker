package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var confPath = "conf.json"
var apiKey = "willbeoverwritten"

const API_KEY_MIN_LENGTH = 12

var conn *pgx.Conn

func convertJsonObjectToString(object interface{}) string {
	var jsonObj, err = json.MarshalIndent(object, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(jsonObj)
}

func printConfig(conf Configuration) {
	// Only show first 5 Characters of api key
	conf.ApiKey = conf.ApiKey[:len(conf.ApiKey)-(len(conf.ApiKey)-5)] + "..."
	fmt.Println("Configuration loaded:")
	fmt.Println(convertJsonObjectToString(conf))
}

func main() {
	var config = readConfig()
	apiKey = config.ApiKey
	if apiKey == "willbeoverwritten" || apiKey == "CHANGEME" {
		fmt.Printf("Invalid Apikey: %s\nEither it wasn't changed or something went wrong!\n", apiKey)
		return
	} else if len(apiKey) < API_KEY_MIN_LENGTH {
		fmt.Printf("Apikey too short: %s\nMust be at least %d characters long!\n", apiKey, API_KEY_MIN_LENGTH)
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func readConfig() Configuration {
	file, err := os.Open(confPath)
	defer file.Close()

	if err != nil {
		log.Fatalf("ERROR - Cannot open config file from '%s': %s", confPath, err)
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("ERROR - cannot decode config: ", err)
	}
	return configuration
}
