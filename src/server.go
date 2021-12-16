package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	r.HandleFunc(fmt.Sprintf("%s/api/statistics", urlPrefix), getStatistics).Methods("GET")
	r.Use(Middleware)

	log.Println(fmt.Sprintf("INFO - Listening on port: %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("INFO - Request from:", r.RemoteAddr, r.URL)
		var apiKey = r.Header.Get("auth")

		if apiKey == authToken {
			next.ServeHTTP(w, r)
			return
		}
		// No permission
		log.Println("ERROR - Invalid Auth Key: " + "'" + apiKey + "'")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Access Denied!")
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my Homepage!")
}

func checkCredentialsValid(creds *Credentials) bool {
	var username string
	var password string
	var err = conn.QueryRow(context.Background(), "SELECT username, pass_key FROM users WHERE username=$1", creds.Username).Scan(&username, &password)
	if err != nil {
		log.Println("ERROR - Credentials check failed:", err)
		return false
	}

	var usernameEqual = strings.Compare(strings.TrimRight(username, "\n"), strings.TrimRight(creds.Username, "\n")) == 0
	var passwordEqual = strings.Compare(strings.TrimRight(password, "\n"), strings.TrimRight(creds.Password, "\n")) == 0

	return (usernameEqual && passwordEqual)
}

func getDataAndCredentials(w http.ResponseWriter, r *http.Request) (DefaultRequest, Credentials, error) {
	decoder := json.NewDecoder(r.Body)

	var defaultRequest DefaultRequest
	err := decoder.Decode(&defaultRequest)
	if err != nil {
		log.Println("ERROR - Decoding request failed:", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return DefaultRequest{}, Credentials{}, err
	}

	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}
	return defaultRequest, creds, nil
}

func sendReponseWithMessageAndStatus(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func getStatistics(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		response, err := getStatisticsByUserId(getUserIdByName(creds.Username))
		if err != nil {
			sendReponseWithMessageAndStatus(w, http.StatusInternalServerError, "error while getting statistics")
		} else {
			reponseJson, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(reponseJson)
		}
	} else {
		sendReponseWithMessageAndStatus(w, http.StatusUnauthorized, "invalid credentials")
	}
}

func addRefuel(w http.ResponseWriter, r *http.Request) {
	request, creds, err := getDataAndCredentials(w, r)
	if err != nil {
		return
	}

	if checkCredentialsValid(&creds) {
		if saveRefuelsByUserId(request.Payload, getUserIdByName(creds.Username)) {
			sendReponseWithMessageAndStatus(w, http.StatusCreated, "created")
		} else {
			sendReponseWithMessageAndStatus(w, http.StatusInternalServerError, "error")
		}

	} else {
		sendReponseWithMessageAndStatus(w, http.StatusUnauthorized, "invalid credentials")
	}
}

func updateRefuel(w http.ResponseWriter, r *http.Request) {
	request, creds, err := getDataAndCredentials(w, r)
	if err != nil {
		return
	}

	if checkCredentialsValid(&creds) {
		if updateRefuelByUserId(request.Payload, getUserIdByName(creds.Username)) {
			sendReponseWithMessageAndStatus(w, http.StatusOK, "updated")
		} else {
			sendReponseWithMessageAndStatus(w, http.StatusInternalServerError, "error")
		}
	} else {
		sendReponseWithMessageAndStatus(w, http.StatusUnauthorized, "invalid credentials")
	}
}

func deleteRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var deletion DeletionRequest
	err := decoder.Decode(&deletion)
	if err != nil {
		log.Println("ERROR - Decoding deletion request failed:", err.Error())
		sendReponseWithMessageAndStatus(w, http.StatusBadRequest, "invalid delete request")
		return
	}

	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	if checkCredentialsValid(&creds) {
		if deleteRefuelByUserId(deletion.Id, getUserIdByName(creds.Username)) {
			sendReponseWithMessageAndStatus(w, http.StatusOK, "deleted")
		} else {
			sendReponseWithMessageAndStatus(w, http.StatusInternalServerError, "error while deleting")
		}
	} else {
		sendReponseWithMessageAndStatus(w, http.StatusUnauthorized, "invalid credentials")
	}
}

func getAllRefuels(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	var startIndex int = 0
	var licensePlate string = "ALL"

	// 0 means all
	var month int = 0
	var year int = 0

	var err error
	values := r.URL.Query()

	if len(values) >= 4 {
		startIndex, err = strconv.Atoi(values["startIndex"][0])
		licensePlate = values["licensePlate"][0]
		month, err = strconv.Atoi(values["month"][0])
		year, err = strconv.Atoi(values["year"][0])
		if err != nil {
			log.Printf("ERROR - while parsing query params: %s", values)
		}
	}

	if checkCredentialsValid(&creds) {
		response, err := getAllRefuelsByUserId(getUserIdByName(creds.Username), startIndex, licensePlate, month, year)
		if err != nil {
			sendReponseWithMessageAndStatus(w, http.StatusInternalServerError, "error while getting all refuels")
		} else {
			reponseJson, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(reponseJson)
		}
	} else {
		sendReponseWithMessageAndStatus(w, http.StatusUnauthorized, "credentials check failed")
	}
}
