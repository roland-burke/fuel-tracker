package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test data

var timeObj_server, err = time.Parse("2006-01-02T15:04:05", "2021-09-04T13:10:25")
var exampleRefuelObj_server = DefaultRequest{
	Payload: []Refuel{
		{
			Description:         "testmethod",
			DateTime:            timeObj_server,
			PricePerLiterInEuro: 1.34,
			TotalAmount:         45.0,
			PricePerLiter:       0.0,
			Currency:            "chf",
			Mileage:             98030.0,
			LicensePlate:        "KN-KN-9999",
		},
	},
}

func init() {
	initLogger()
	initDb()
}

func TestSendResponseWithMessageAndStatus(t *testing.T) {
	assert := assert.New(t)
	// setup
	recorder := httptest.NewRecorder()

	// test
	sendResponseWithMessageAndStatus(recorder, 200, "Test")

	assert.Equal(200, recorder.Code)
	assert.Equal("Test", recorder.Body.String())
}

func TestCheckCredentials(t *testing.T) {
	assert := assert.New(t)

	// setup
	req, err := http.NewRequest("GET", "/unimportant", nil)
	assert.Nil(err)

	// test right credentials
	req.Header.Set("username", "john")
	req.Header.Set("password", "john")

	result := checkCredentials(req)
	assert.True(result)

	// test wrong credentials
	req.Header.Set("username", "not existing")
	req.Header.Set("password", "also")

	result = checkCredentials(req)
	assert.False(result)
}

func TestGetDefaultReuqestObj(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	json, err := json.Marshal(exampleRefuelObj_server)

	req, err := http.NewRequest("GET", "/unimportant", bytes.NewBuffer(json))
	assert.Nil(err)

	req.Header.Set("username", "john")
	req.Header.Set("password", "john")

	defaultReq, err := getDefaultRequestObj(recorder, req)
	assert.Nil(err)

	assert.Equal(exampleRefuelObj_server, defaultReq)
}

func TestIntermediate(t *testing.T) {
	assert := assert.New(t)

	// setup
	apiKey = "asdfasdf"
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/whatever", nil)
	assert.Nil(err)

	// test wrong apikey, right credentials
	req.Header.Set("auth", "wrong")
	req.Header.Set("username", "john")
	req.Header.Set("password", "john")

	intermediate(recorder, req, nil)

	assert.Equal(http.StatusUnauthorized, recorder.Result().StatusCode)
	assert.Equal("API access denied!", recorder.Body.String())

	// test right apikey, wrong credentials
	recorder = httptest.NewRecorder()
	req.Header.Set("auth", "asdfasdf")
	req.Header.Set("username", "foo")
	req.Header.Set("password", "baaar")

	intermediate(recorder, req, nil)
	assert.Equal(http.StatusUnauthorized, recorder.Result().StatusCode)
	assert.Equal("Credentials Check failed", recorder.Body.String())
}

func TestGetStatistics(t *testing.T) {
	assert := assert.New(t)

	// setup
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/whatever", nil)
	assert.Nil(err)

	expectedStats := StatisticsResponse{
		Stats:        []Stat{},
		TotalMileage: 700,
		TotalCost:    123.75,
	}

	// test
	req.Header.Set("username", "john")

	// When
	getStatistics(recorder, req)

	res := recorder.Result()
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var statsResponse StatisticsResponse
	err = decoder.Decode(&statsResponse)
	assert.Nil(err)

	// Then
	assert.Equal(http.StatusOK, recorder.Result().StatusCode)
	assert.Equal(expectedStats.TotalMileage, statsResponse.TotalMileage)
	assert.Equal(expectedStats.TotalCost, statsResponse.TotalCost)
}

func TestAddRefuel(t *testing.T) {
	assert := assert.New(t)

	// setup
	recorder := httptest.NewRecorder()

	json, err := json.Marshal(exampleRefuelObj_server)
	req, err := http.NewRequest("POST", "/whatever", bytes.NewBuffer(json))
	assert.Nil(err)

	// test
	req.Header.Set("username", "john")

	// When
	addRefuel(recorder, req)

	// Then
	assert.Equal(http.StatusCreated, recorder.Result().StatusCode)
	assert.Equal("Successfully added", recorder.Body.String())

	// Cleanup

	deleteRefuelByUserId(exampleRefuelObj1_repository.Id, 1)
}
