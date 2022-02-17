package main

import (
	"encoding/json"
	"fmt"
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

	logger.Info("Listening on port: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func sendResponseWithMessageAndStatus(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request from: %s %s", r.RemoteAddr, r.URL)
		var apiKeyFromClient = r.Header.Get("auth")

		if apiKeyFromClient == apiKey {
			credentialsValid := checkCredentials(r)

			if credentialsValid {
				next.ServeHTTP(w, r)
				return
			} else {
				sendResponseWithMessageAndStatus(w, http.StatusUnauthorized, "Credentials Check failed")
			}
		} else {
			// No api permission
			logger.Warn("Invalid Apikey: %s", apiKeyFromClient)
			sendResponseWithMessageAndStatus(w, http.StatusUnauthorized, "API access denied!")
		}
	})
}

func checkCredentials(r *http.Request) bool {
	creds := Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	err, username, password := getCredentials(creds.Username)

	if err != nil {
		logger.Error("Credentials check failed: %s", err.Error())
		return false
	}

	var usernameEqual = strings.Compare(strings.TrimRight(username, "\n"), strings.TrimRight(creds.Username, "\n")) == 0
	var passwordEqual = strings.Compare(strings.TrimRight(password, "\n"), strings.TrimRight(creds.Password, "\n")) == 0

	return (usernameEqual && passwordEqual)
}

func getDefaultRequestObj(w http.ResponseWriter, r *http.Request) (DefaultRequest, error) {
	decoder := json.NewDecoder(r.Body)

	var defaultRequest DefaultRequest
	err := decoder.Decode(&defaultRequest)
	if err != nil {
		logger.Error("Decoding request failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusBadRequest, err.Error())
		return DefaultRequest{}, err
	}
	return defaultRequest, nil
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my Homepage!")
}

func getStatistics(w http.ResponseWriter, r *http.Request) {
	var username = r.Header.Get("username")

	response, err := getStatisticsByUserId(getUserIdByName(username))
	if err != nil {
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, "error while getting statistics")
		return
	}
	responseJson, _ := json.Marshal(response)
	sendResponseWithMessageAndStatus(w, http.StatusOK, string(responseJson))
}

func addRefuel(w http.ResponseWriter, r *http.Request) {
	var username = r.Header.Get("username")
	request, err := getDefaultRequestObj(w, r)
	if err != nil {
		return
	}

	err, _ = saveRefuelByUserId(request.Payload[0], getUserIdByName(username))
	if err != nil {
		logger.Error("Saving refuel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusCreated, "Successfully added")
}

func updateRefuel(w http.ResponseWriter, r *http.Request) {
	var username = r.Header.Get("username")
	request, err := getDefaultRequestObj(w, r)
	if err != nil {
		return
	}
	err = updateRefuelByUserId(request.Payload[0], getUserIdByName(username))
	if err != nil {
		logger.Error("Updating reufel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusOK, "Successfully updated")
}

func deleteRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var username = r.Header.Get("username")

	var deletion DeletionRequest
	err := decoder.Decode(&deletion)
	if err != nil {
		logger.Error("Decoding deletion request failed:: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusBadRequest, "invalid delete request")
		return
	}

	err = deleteRefuelByUserId(deletion.Id, getUserIdByName(username))
	if err != nil {
		logger.Error("Deleting reufel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusOK, "Successfully deleted")
}

func getAllRefuels(w http.ResponseWriter, r *http.Request) {
	var username = r.Header.Get("username")

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
			logger.Error("Parsing query params failed: %s - %s", values, err)
		}
	}

	response, err := getAllRefuelsByUserId(getUserIdByName(username), startIndex, licensePlate, month, year)
	if err != nil {
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJson, _ := json.Marshal(response)
	sendResponseWithMessageAndStatus(w, http.StatusOK, string(responseJson))
}
