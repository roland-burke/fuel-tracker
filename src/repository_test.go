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

var timeObj1_repository, _ = time.Parse("2006-02-01T15:04:05", "2021-09-04T13:10:25")
var timeObj2_repository, _ = time.Parse("2006-02-01T15:04:05", "2021-09-05T16:34:25")

var exampleRefuelObj1_repository = Refuel{
	Id:                  4,
	Description:         "TestRefuel1",
	DateTime:            timeObj1_repository,
	PricePerLiterInEuro: 1.2,
	TotalAmount:         35,
	PricePerLiter:       40,
	Currency:            "Chf",
	Mileage:             42100,
	LicensePlate:        "KN-KN-9999",
	LastChanged:         time.Now(),
}

var exampleRefuelObj2_repository = Refuel{
	Id:                  5,
	Description:         "TestRefuel2",
	DateTime:            timeObj2_repository,
	PricePerLiterInEuro: 1.234,
	TotalAmount:         55,
	PricePerLiter:       40,
	Currency:            "Chf",
	Mileage:             43100,
	LicensePlate:        "KN-KN-9999",
	LastChanged:         time.Now(),
}

func TestGetUserId(t *testing.T) {
	assert := assert.New(t)

	// When
	var expectedUserIdJohn = 1
	var expectedUserIdMary = 2
	var expecteduserIdInvalid = -1

	var userIdJohn = getUserIdByCredentials("john", "john")
	var userIdMary = getUserIdByCredentials("mary", "mary")
	var userIdInvalid = getUserIdByCredentials("not-existing", "asdf")

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
	var userId = 1

	// When
	err, refuelId := saveRefuelByUserId(exampleRefuelObj1_repository, userId)

	// Then
	assert.Nil(err)
	assert.Equal(refuelId, refuelId)

	// Cleanup
	err = deleteRefuelByUserId(refuelId, userId)
	assert.Nil(err)
}

func TestUpdateRefuel(t *testing.T) {
	assert := assert.New(t)

	// Setup
	err, refuelId := saveRefuelByUserId(exampleRefuelObj1_repository, 1)
	assert.Nil(err)

	// When
	var userId = 1
	err = updateRefuelByUserId(exampleRefuelObj2_repository, userId)
	assert.Nil(err)

	refuelResponse, err := getAllRefuelsByUserId(userId, 0, "KN-KN-9999", 0, 0)
	assert.Nil(err)

	// Then
	var targetRefuel = refuelResponse.Refuels[len(refuelResponse.Refuels)-1]

	assert.Equal(exampleRefuelObj2_repository.Description, targetRefuel.Description)
	assert.Equal(exampleRefuelObj2_repository.DateTime, targetRefuel.DateTime)
	assert.Equal(exampleRefuelObj2_repository.PricePerLiterInEuro, targetRefuel.PricePerLiterInEuro)
	assert.Equal(exampleRefuelObj2_repository.TotalAmount, targetRefuel.TotalAmount)
	assert.Equal(exampleRefuelObj2_repository.PricePerLiter, targetRefuel.PricePerLiter)
	assert.Equal(exampleRefuelObj2_repository.Currency, targetRefuel.Currency)
	assert.Equal(exampleRefuelObj2_repository.Mileage, targetRefuel.Mileage)
	assert.Equal(exampleRefuelObj2_repository.LicensePlate, targetRefuel.LicensePlate)

	// Cleanup
	err = deleteRefuelByUserId(refuelId, userId)
	assert.Nil(err)
}

func TestDeleteRefuel(t *testing.T) {
	assert := assert.New(t)

	// Setup
	var userId = 1

	err, refuelId := saveRefuelByUserId(exampleRefuelObj1_repository, userId)
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

func TestGetStatisticsByUserId(t *testing.T) {
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
