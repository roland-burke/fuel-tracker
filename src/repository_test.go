package main

import (
	"testing"
	"time"
)

func init() {
	initLogger()
	initDb()
}

func TestGetUserId(t *testing.T) {
	var expectedUserIdJohn = 1
	var expectedUserIdMary = 2
	var expecteduserIdInvalid = -1

	var userIdJohn = getUserIdByName("john")
	var userIdMary = getUserIdByName("mary")
	var userIdInvalid = getUserIdByName("not-existing")

	if userIdJohn != expectedUserIdJohn {
		t.Errorf("Got user id: %d, expected: %d", userIdJohn, expectedUserIdJohn)
	}
	if userIdMary != expectedUserIdMary {
		t.Errorf("Got user id: %d, expected: %d", userIdMary, expectedUserIdMary)
	}
	if userIdInvalid != -1 {
		t.Errorf("Got user id: %d, expected: %d", userIdInvalid, expecteduserIdInvalid)
	}
}

func TestGetCredentials(t *testing.T) {
	err, username, password := getCredentials("john")

	if err != nil {
		t.Errorf("Failed getting credentials: %s", err.Error())
	}

	if username != "john" || password != "john" {
		t.Errorf("Found wrong credentials: %s, %s, expected: john, john", username, password)
	}

	err, username, password = getCredentials("not_exist")

	if err == nil {
		t.Errorf("Did not throw error on non existing username")
	}

	if username != "" || password != "" {
		t.Errorf("Username or password was not empty on error: %s, %s", username, password)
	}
}

func TestSaveRefuel(t *testing.T) {
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

	// Test
	err, refuelId := saveRefuelByUserId(refuel, userId)

	if err != nil {
		t.Errorf("Save user with id: %d failed: %s", userId, err.Error())
	}

	if refuelId != 4 {
		t.Errorf("SaveRefuelByUserId returned wrong refuelId: %d, expected: %d", refuelId, 4)
	}

	// cleanup
	deleteRefuelByUserId(refuelId, userId)
}

func TestUpdateRefuel(t *testing.T) {
	timeObj, err := time.Parse("2006-02-01T15:04:05", "2021-09-04T13:10:25")

	if err != nil {
		logger.Error(err.Error())
	}

	newRefuel := Refuel{
		Id:                  1,
		Description:         "Test",
		DateTime:            timeObj,
		PricePerLiterInEuro: 1.439,
		TotalAmount:         42.0,
		PricePerLiter:       1.488,
		Currency:            "Chf",
		Mileage:             40100,
		LicensePlate:        "KN-KN-9999",
		LastChanged:         time.Now(),
	}

	var userId = 1
	err = updateRefuelByUserId(newRefuel, userId)

	if err != nil {
		t.Errorf("Updating refuel with userId: %d failed: %s", userId, err.Error())
	}

	refuelResponse, err := getAllRefuelsByUserId(userId, 0, "KN-KN-9999", 0, 0)

	var targetRefuel = refuelResponse.Refuels[len(refuelResponse.Refuels)-1]

	if targetRefuel.Description != "Test" {
		t.Errorf("Updating refuel with userId: %d failed, description: %s", userId, targetRefuel.Description)
	}
}

func TestDeleteRefuel(t *testing.T) {
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

	if err != nil {
		t.Errorf("Save refuel with userId: %d failed: %s", userId, err.Error())
	}

	// Test
	err = deleteRefuelByUserId(refuelId, userId)
	if err != nil {
		t.Errorf("Delete refuel with userId: %d failed: %s", userId, err.Error())
	}

	err = deleteRefuelByUserId(refuelId, userId)
	if err == nil {
		t.Errorf("Deleted refuel with userId twice: userId: %d, refuelId: %d", userId, refuelId)
	}
}

func TestGetAllRefuels(t *testing.T) {
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

	if refuelResponse.TotalCount != expectedRefuelResponse.TotalCount {
		t.Errorf("Get all Refuels returned wrong totalCount: %d, expected: %d", refuelResponse.TotalCount, expectedRefuelResponse.TotalCount)
	}

	if err != nil {
		t.Errorf("Get all Refuels failed: %s", err.Error())
	}
}

func TestGetStatistics(t *testing.T) {

	expectedStats := StatisticsResponse{
		Stats:        []Stat{},
		TotalMileage: 700,
		TotalCost:    123.75,
	}

	statistics, err := getStatisticsByUserId(1)

	if err != nil {
		t.Errorf("Get Statistics failed: %s", err.Error())
	}

	if statistics.TotalCost != expectedStats.TotalCost {
		t.Errorf("Got Total Cost: %f, expected: %f", statistics.TotalCost, expectedStats.TotalCost)
	}

	if statistics.TotalMileage != expectedStats.TotalMileage {
		t.Errorf("Got Total Mileage: %f, expected: %f", statistics.TotalMileage, expectedStats.TotalMileage)
	}
}
