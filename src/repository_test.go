package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	initLogger()
	initDb()
}

func TestGetUserId(t *testing.T) {
	assert := assert.New(t)

	// When
	var expectedUserIdJohn = 1
	var expectedUserIdMary = 2
	var expecteduserIdInvalid = -1

	var userIdJohn = getUserIdByName("john")
	var userIdMary = getUserIdByName("mary")
	var userIdInvalid = getUserIdByName("not-existing")

	// Then
	assert.Equal(expectedUserIdJohn, userIdJohn)
	assert.Equal(expectedUserIdMary, userIdMary)
	assert.Equal(expecteduserIdInvalid, userIdInvalid)
}

func TestGetCredentials(t *testing.T) {
	assert := assert.New(t)

	// When
	err, username, password := getCredentials("john")

	// Then
	assert.Nil(err)
	assert.Equal("john", username)
	assert.Equal("john", password)

	// When
	err, username, password = getCredentials("not_exist")

	// Then
	assert.NotNil(err)
	assert.Equal("", username)
	assert.Equal("", password)
}

func TestSaveRefuel(t *testing.T) {
	assert := assert.New(t)

	// Setup
	refuel := Refuel{
		Id:                  0,
		Description:         "Test",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200450,
		LicensePlate:        "KN-KN-420",
		LastChanged:         time.Now(),
	}

	var userId = 1

	// When
	err, refuelId := saveRefuelByUserId(refuel, userId)

	// Then
	assert.Nil(err)
	assert.Equal(4, refuelId)

	// Cleanup
	deleteRefuelByUserId(refuelId, userId)
}

func TestUpdateRefuel(t *testing.T) {
	assert := assert.New(t)

	timeObj1, err := time.Parse("2006-02-01T15:04:05", "2021-09-04T13:10:25")

	if err != nil {
		logger.Error(err.Error())
	}

	expectedRefuel := Refuel{
		Id:                  1,
		Description:         "Test",
		DateTime:            timeObj1,
		PricePerLiterInEuro: 1.439,
		TotalAmount:         42.0,
		PricePerLiter:       1.488,
		Currency:            "Chf",
		Mileage:             40100,
		LicensePlate:        "KN-KN-9999",
		LastChanged:         time.Now(),
	}

	// When
	var userId = 1
	err = updateRefuelByUserId(expectedRefuel, userId)
	assert.Nil(err)

	refuelResponse, err := getAllRefuelsByUserId(userId, 0, "KN-KN-9999", 0, 0)
	assert.Nil(err)

	// Then
	var targetRefuel = refuelResponse.Refuels[len(refuelResponse.Refuels)-1]

	assert.Equal(expectedRefuel.Description, targetRefuel.Description)
	assert.Equal(expectedRefuel.DateTime, targetRefuel.DateTime)
	assert.Equal(expectedRefuel.PricePerLiterInEuro, targetRefuel.PricePerLiterInEuro)
	assert.Equal(expectedRefuel.TotalAmount, targetRefuel.TotalAmount)
	assert.Equal(expectedRefuel.PricePerLiter, targetRefuel.PricePerLiter)
	assert.Equal(expectedRefuel.Currency, targetRefuel.Currency)
	assert.Equal(expectedRefuel.Mileage, targetRefuel.Mileage)
	assert.Equal(expectedRefuel.LicensePlate, targetRefuel.LicensePlate)
}

func TestDeleteRefuel(t *testing.T) {
	assert := assert.New(t)

	// Setup
	refuel := Refuel{
		Id:                  9999,
		Description:         "Test",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.67,
		TotalAmount:         42,
		PricePerLiter:       1.78,
		Currency:            "Chf",
		Mileage:             200460,
		LicensePlate:        "KN-KN-420",
		LastChanged:         time.Now(),
	}

	var userId = 1

	err, refuelId := saveRefuelByUserId(refuel, userId)
	assert.Nil(err)

	// Test
	err = deleteRefuelByUserId(refuelId, userId)
	assert.Nil(err)

	err = deleteRefuelByUserId(refuelId, userId)
	assert.NotNil(err)
}

func TestGetAllRefuels(t *testing.T) {
	assert := assert.New(t)

	refuel := Refuel{
		Id:                  0,
		Description:         "Test",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200450,
		LicensePlate:        "KN-KN-9999",
		LastChanged:         time.Now(),
	}

	expectedRefuelResponse := RefuelResponse{
		Refuels:    []Refuel{refuel},
		TotalCount: 2,
	}

	var userId = 1

	refuelResponse, err := getAllRefuelsByUserId(userId, 0, "KN-KN-9999", 0, 0)
	assert.Nil(err)

	assert.Equal(expectedRefuelResponse.TotalCount, refuelResponse.TotalCount)
}

func TestGetStatistics(t *testing.T) {
	assert := assert.New(t)

	// Setup
	expectedStats := StatisticsResponse{
		Stats:        []Stat{},
		TotalMileage: 700,
		TotalCost:    123.75,
	}

	statistics, err := getStatisticsByUserId(1)
	assert.Nil(err)

	assert.Equal(expectedStats.TotalCost, statistics.TotalCost)
	assert.Equal(expectedStats.TotalMileage, statistics.TotalMileage)
}
