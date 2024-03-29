package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/roland-burke/fuel-tracker/internal/config"
	"github.com/roland-burke/fuel-tracker/internal/model"
	"github.com/roland-burke/fuel-tracker/internal/repository"
)

func StartServer(port int, urlPrefix string) {
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
	r.Use(middleware)

	config.Logger.Info("Listening on port: %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func sendResponseWithMessageAndStatus(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		intermediate(w, r, next)
	})
}

func intermediate(w http.ResponseWriter, r *http.Request, next http.Handler) {
	config.Logger.Info("Request from: %s %s", r.RemoteAddr, r.URL)
	var apiKeyFromClient = r.Header.Get("auth")

	if apiKeyFromClient == config.ApiKey {
		credentialsValid := checkCredentials(r)

		if credentialsValid {
			next.ServeHTTP(w, r)
			return
		} else {
			sendResponseWithMessageAndStatus(w, http.StatusUnauthorized, "Credentials Check failed")
		}
	} else {
		// No api permission
		config.Logger.Warn("Invalid Apikey: %s", apiKeyFromClient)
		sendResponseWithMessageAndStatus(w, http.StatusUnauthorized, "API access denied!")
	}
}

func checkCredentials(r *http.Request) bool {
	creds := model.Credentials{
		Username: r.Header.Get("username"),
		Password: r.Header.Get("password"),
	}

	err, username, password := repository.GetCredentials(creds.Username)

	if err != nil {
		config.Logger.Error("Credentials check failed: %s", err.Error())
		return false
	}

	var usernameEqual = strings.Compare(strings.TrimRight(username, "\n"), strings.TrimRight(creds.Username, "\n")) == 0
	var passwordEqual = strings.Compare(strings.TrimRight(password, "\n"), strings.TrimRight(creds.Password, "\n")) == 0

	return (usernameEqual && passwordEqual)
}

func getDefaultRequestObj(w http.ResponseWriter, r *http.Request) (model.DefaultRequest, error) {
	decoder := json.NewDecoder(r.Body)

	var defaultRequest model.DefaultRequest
	err := decoder.Decode(&defaultRequest)
	if err != nil {
		config.Logger.Error("Decoding request failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusBadRequest, "Failed to decode request")
		return model.DefaultRequest{}, err
	}
	return defaultRequest, nil
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my Homepage!")
}

func getStatistics(w http.ResponseWriter, r *http.Request) {
	var params = r.URL.Query()["licensePlate"]
	var licensePlate = "ALL"
	if len(params) >= 1 {
		licensePlate = params[0]
	}
	response := repository.GetStatisticsByUserId(repository.GetUserIdByCredentials(r.Header.Get("username"), r.Header.Get("password")), licensePlate)

	responseJson, err := json.Marshal(response)
	if err != nil {
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, "Error while parsing statistics")
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusOK, string(responseJson))
}

func addRefuel(w http.ResponseWriter, r *http.Request) {
	request, err := getDefaultRequestObj(w, r)
	if err != nil {
		return
	}

	_, err = repository.SaveRefuelByUserId(request.Payload[0], repository.GetUserIdByCredentials(r.Header.Get("username"), r.Header.Get("password")))
	if err != nil {
		config.Logger.Error("Saving refuel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusCreated, "Successfully added")
}

func updateRefuel(w http.ResponseWriter, r *http.Request) {
	request, err := getDefaultRequestObj(w, r)
	if err != nil {
		return
	}
	err = repository.UpdateRefuelByUserId(request.Payload[0], repository.GetUserIdByCredentials(r.Header.Get("username"), r.Header.Get("password")))
	if err != nil {
		config.Logger.Error("Updating refuel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusOK, "Successfully updated")
}

func deleteRefuel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var deletion model.DeletionRequest
	err := decoder.Decode(&deletion)
	if err != nil {
		config.Logger.Error("Decoding deletion request failed:: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusBadRequest, "invalid delete request")
		return
	}

	err = repository.DeleteRefuelByUserId(deletion.Id, repository.GetUserIdByCredentials(r.Header.Get("username"), r.Header.Get("password")))
	if err != nil {
		config.Logger.Error("Deleting reufel failed: %s", err.Error())
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseWithMessageAndStatus(w, http.StatusOK, "Successfully deleted")
}

func getAllRefuels(w http.ResponseWriter, r *http.Request) {
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
			config.Logger.Error("Parsing query params failed: %s - %s", values, err)
		}
	}

	response, err := repository.GetAllRefuelsByUserId(repository.GetUserIdByCredentials(r.Header.Get("username"), r.Header.Get("password")), startIndex, licensePlate, month, year)
	if err != nil {
		sendResponseWithMessageAndStatus(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJson, _ := json.Marshal(response)
	sendResponseWithMessageAndStatus(w, http.StatusOK, string(responseJson))
}
