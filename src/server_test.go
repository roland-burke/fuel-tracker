package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	var jsonStr = []byte(`{
		"payload": [
			{
			"description": "testmethod",
			"dateTime": "2021-09-04T13:10:25Z",
			"pricePerLiterInEuro": 1.34,
			"totalAmount": 45.0,
			"pricePerLiter": 0.0,
			"currency": "chf",
			"mileage": 340.0,
			"licensePlate": "KN-KN-9999"
		}]
		
	}`)

	timeObj, err := time.Parse("2006-01-02T15:04:05", "2021-09-04T13:10:25")
	assert.Nil(err)

	var expectedRequest = DefaultRequest{
		Payload: []Refuel{
			{
				Description:         "testmethod",
				DateTime:            timeObj,
				PricePerLiterInEuro: 1.34,
				TotalAmount:         45.0,
				PricePerLiter:       0.0,
				Currency:            "chf",
				Mileage:             340.0,
				LicensePlate:        "KN-KN-9999",
			},
		},
	}

	req, err := http.NewRequest("GET", "/unimportant", bytes.NewBuffer(jsonStr))
	assert.Nil(err)

	req.Header.Set("username", "john")
	req.Header.Set("password", "john")

	defaultReq, err := getDefaultRequestObj(recorder, req)
	assert.Nil(err)

	assert.Equal(expectedRequest, defaultReq)
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
