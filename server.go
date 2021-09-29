package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

func getUserIdByName(username string) int {
	var user_id int
	var err = conn.QueryRow(context.Background(), "SELECT users_id FROM users WHERE username=$1", username).Scan(&user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return -1
	}

	return user_id
}

func checkCredentialsValid(creds *Credentials) bool {
	var username string
	var password string
	var err = conn.QueryRow(context.Background(), "SELECT username, pass_key FROM users WHERE username=$1", creds.Username).Scan(&username, &password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return false
	}

	return username == creds.Username && password == creds.Password
}

func addRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request DefaultRequest
	err := decoder.Decode(&request)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		refuel := request.Payload

		if saveRefuelByUserId(&refuel, getUserIdByName(creds.Username)) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("created"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}
}

func updateRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request DefaultRequest
	err := decoder.Decode(&request)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		if updateRefuelByUserId(&request.Payload, getUserIdByName(creds.Username)) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("updated"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}
}

func deleteRefuel(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var deletion DeletionRequest
	err := decoder.Decode(&deletion)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		if deleteRefuelByUserId(deletion.Id, getUserIdByName(creds.Username)) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("deleted"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}
}

func getAllRefuels(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		response, err := getAllRefuelsByUserId(getUserIdByName(creds.Username))

		if err != nil {
			log.Fatal(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			reponseJson, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(reponseJson)
		}

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
	}
}
