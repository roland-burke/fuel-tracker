package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var authToken = "willbeoverwritten"
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
	authToken = config.AuthToken
	var port = config.Port
	var urlPrefix = config.UrlPrefix
	printConfig(config)
	initDb()
	startServer(port, urlPrefix)
}

func initDb() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer conn.Close(context.Background())
}

func readConfig() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Error while reading config:", err)
	}
	return configuration
}
