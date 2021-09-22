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

/*
func dbTest() {
	urlExample := "postgres://postgres:fj498h89fm89dhfi3@db:5432/fuel_tracker"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var pricePerLiterInEuro float64
	var totalLiter float64
	err = conn.QueryRow(context.Background(), "select price_per_liter_euro, total_liter, name from refuel where id=$1", 0).Scan(&pricePerLiterInEuro, &totalLiter, &name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("name: %s, price per liter: %f, total liter: %f", name, pricePerLiterInEuro, totalLiter)

}
*/

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
