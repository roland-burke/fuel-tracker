package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func startServer(port int, urlPrefix string) {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc(fmt.Sprintf("%s/api/add", urlPrefix), addRefuel).Methods("POST")
	r.HandleFunc(fmt.Sprintf("%s/api/delete", urlPrefix), deleteRefuel).Methods("DELETE")
	r.HandleFunc(fmt.Sprintf("%s/api/update", urlPrefix), updateRefuel).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("%s/api/get", urlPrefix), getRefuel).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/api/get/all", urlPrefix), getAllRefuels).Methods("GET")
	r.Use(Middleware)

	println(fmt.Sprintf("Listening on port: %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request from:", r.RemoteAddr, r.URL)
		var apiKey = r.Header.Get("auth")

		if apiKey == authToken {
			next.ServeHTTP(w, r)
			return
		}
		// No permission
		log.Println("Invalid Auth Key: " + apiKey)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Acess Denied!")
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my Homepage!")
}

func addRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var refuel Refuel
	err := decoder.Decode(&refuel)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO refuel(name, date_time, price_per_liter_euro, total_liter, price_per_liter, currency) VALUES($1, $2, $3, $4, $5, $6)", refuel.Name, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func updateRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var refuel Refuel
	err := decoder.Decode(&refuel)
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(context.Background(), "UPDATE refuel SET name=$1, date_time=$2, price_per_liter_euro=$3, total_liter=$4, price_per_liter=$5, currency=$6 where id=$7", refuel.Name, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Id)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("updated"))
}

func deleteRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var deletion Deletion
	err := decoder.Decode(&deletion)
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(context.Background(), "DELETE FROM refuel WHERE id=$1", deletion.Id)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("deleted"))
}

func getRefuel(w http.ResponseWriter, r *http.Request) {
	refuel := Refuel{
		Name:                "mocked data",
		DateTime:            time.Now(),
		PricePerLiter:       1.38,
		TotalAmount:         45.32,
		Currency:            "euro",
		PricePerLiterInEuro: 1.38,
	}

	reponseJson, err := json.Marshal(refuel)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(reponseJson)

}
func getAllRefuels(w http.ResponseWriter, r *http.Request) {
	var err error
	rows, err := conn.Query(context.Background(), "SELECT * FROM refuel")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}

	var refuelListBuffer [100]Refuel

	var index = 0

	for rows.Next() {
		var id int
		var name string
		var dateTime time.Time
		var pricePerLiterInEuro float64
		var totalAmount float64
		var pricePerLiter float64
		var currency string
		var lastChanged time.Time

		err := rows.Scan(&id, &name, &dateTime, &pricePerLiterInEuro, &totalAmount, &pricePerLiter, &currency, &lastChanged)
		if err != nil {
			fmt.Fprintf(os.Stderr, "row next failed: %v\n", err)
			os.Exit(1)
		}

		refuelListBuffer[index] = Refuel{
			Id:                  id,
			Name:                name,
			DateTime:            dateTime,
			PricePerLiter:       pricePerLiter,
			TotalAmount:         totalAmount,
			Currency:            currency,
			PricePerLiterInEuro: pricePerLiterInEuro,
			LastChanged:         lastChanged,
		}
		index += 1
		fmt.Printf("id: %d, name: %s, totalliter: %f\n", id, name, totalAmount)
	}

	reponseJson, err := json.Marshal(refuelListBuffer[:index])

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(reponseJson)

}
