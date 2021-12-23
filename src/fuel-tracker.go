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
var conn *pgx.Conn

func convertJsonObjectToString(object interface{}) string {
	var jsonObj, err = json.MarshalIndent(object, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(jsonObj)
}

func printConfig(conf Configuration) {
	fmt.Println("Configuration loaded:")
	fmt.Println(convertJsonObjectToString(conf))
	fmt.Println("============================================")
}

func main() {
	var config = readConfig()
	apiKey = config.ApiKey
	if apiKey == "willbeoverwritten" || apiKey == "CHANGEME" {
		fmt.Printf("Invalid Apikey: %s\nEither it wasn't changed or something went wrong!\n", apiKey)
		return
	}
	var port = config.Port
	var urlPrefix = config.UrlPrefix
	printConfig(config)
	initDb()
	startServer(port, urlPrefix)
}

func initDb() {
	var err error
	fmt.Printf("db url: %s\n", os.Getenv("DATABASE_URL"))
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
		log.Fatal("ERROR - can't decode config: ", err)
	}
	return configuration
}
