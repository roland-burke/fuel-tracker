package main

import (
	"testing"
	"time"
)

func init() {
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

func TestSaveRefuels(t *testing.T) {
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

	err := saveRefuelsByUserId([]Refuel{refuel}, userId)

	if err != nil {
		t.Errorf("Save user with id: %d failed", userId)
	}
}
