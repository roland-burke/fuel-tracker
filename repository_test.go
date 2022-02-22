package main

import (
	"testing"
	"time"

	"github.com/roland-burke/rollogger"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Mute the logger
	logger = rollogger.Init(rollogger.ERROR_LEVEL, true, true)
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

func TestGetUserIdByCredentials(t *testing.T) {
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
	refuelId, err := saveRefuelByUserId(exampleRefuelObj1_repository, userId)

	// Then
	assert.Nil(err)
	assert.Equal(refuelId, refuelId)

	// Cleanup
	err = deleteRefuelByUserId(refuelId, userId)
	assert.Nil(err)
}

func TestUpdateRefuelByUserId(t *testing.T) {
	assert := assert.New(t)
	var userId = 1

	// Setup
	refuelId, err := saveRefuelByUserId(exampleRefuelObj1_repository, userId)
	assert.Nil(err)

	exampleRefuelObj2_repository.Id = refuelId

	// When
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

func TestDeleteRefuelByUserId(t *testing.T) {
	assert := assert.New(t)

	// Setup
	var userId = 1

	refuelId, err := saveRefuelByUserId(exampleRefuelObj1_repository, userId)
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

	statistics := getStatisticsByUserId(1, "ALL")

	assert.Equal(expectedStats.TotalCost, statistics.TotalCost)
	assert.Equal(expectedStats.TotalMileage, statistics.TotalMileage)
}

func TestPaging(t *testing.T) {
	assert := assert.New(t)
	var userId = 1

	testObj1 := Refuel{
		Id:                  0,
		Description:         "paging-test1",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200500,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj2 := Refuel{
		Id:                  0,
		Description:         "paging-test2",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200600,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj3 := Refuel{
		Id:                  0,
		Description:         "paging-test4",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200700,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj4 := Refuel{
		Id:                  0,
		Description:         "paging-test4",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200800,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj5 := Refuel{
		Id:                  0,
		Description:         "paging-test5",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             200900,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj6 := Refuel{
		Id:                  0,
		Description:         "paging-test6",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             201000,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj7 := Refuel{
		Id:                  0,
		Description:         "paging-test7",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             201100,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj8 := Refuel{
		Id:                  0,
		Description:         "paging-test8",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             201200,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	testObj9 := Refuel{
		Id:                  0,
		Description:         "paging-test9",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             201300,
		LicensePlate:        "KN-KN-2323",
		LastChanged:         time.Now(),
	}

	refuelId1, err := saveRefuelByUserId(testObj1, userId)
	assert.Nil(err)
	refuelId2, err := saveRefuelByUserId(testObj2, userId)
	assert.Nil(err)
	refuelId3, err := saveRefuelByUserId(testObj3, userId)
	assert.Nil(err)

	result, err := getAllRefuelsByUserId(userId, 0, "KN-KN-2323", 0, 0)
	assert.Nil(err)

	assert.Equal(3, result.TotalCount)
	assert.Equal(3, len(result.Refuels))

	result, err = getAllRefuelsByUserId(userId, 1, "KN-KN-2323", 0, 0)
	assert.Nil(err)

	assert.Equal(3, result.TotalCount)
	assert.Equal(2, len(result.Refuels))

	refuelId4, err := saveRefuelByUserId(testObj4, userId)
	assert.Nil(err)
	refuelId5, err := saveRefuelByUserId(testObj5, userId)
	assert.Nil(err)
	refuelId6, err := saveRefuelByUserId(testObj6, userId)
	assert.Nil(err)
	refuelId7, err := saveRefuelByUserId(testObj7, userId)
	assert.Nil(err)
	refuelId8, err := saveRefuelByUserId(testObj8, userId)
	assert.Nil(err)
	refuelId9, err := saveRefuelByUserId(testObj9, userId)
	assert.Nil(err)

	result1, err := getAllRefuelsByUserId(userId, 0, "KN-KN-2323", 0, 0)
	assert.Nil(err)

	assert.Equal(9, result1.TotalCount)
	assert.Equal(8, len(result1.Refuels))

	result2, err := getAllRefuelsByUserId(userId, 5, "KN-KN-2323", 0, 0)
	assert.Nil(err)

	assert.Equal(9, result2.TotalCount)
	assert.Equal(4, len(result2.Refuels))

	// cleanup
	err = deleteRefuelByUserId(refuelId1, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId2, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId3, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId4, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId5, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId6, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId7, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId8, userId)
	assert.Nil(err)
	err = deleteRefuelByUserId(refuelId9, userId)
	assert.Nil(err)

}

// === Test Database Constraints ===

func TestSaveRefuelWithBadMileage(t *testing.T) {
	assert := assert.New(t)

	// Setup
	var userId = 1

	var refuelWithInvalidMileage = Refuel{
		Id:                  4,
		Description:         "TestRefuel1",
		DateTime:            time.Now(),
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             300,
		LicensePlate:        "KN-KN-9999",
		LastChanged:         time.Now(),
	}

	// When
	refuelId, err := saveRefuelByUserId(refuelWithInvalidMileage, userId)

	// Then
	assert.NotNil(err)
	assert.Equal(-1, refuelId)
	assert.Equal("ERROR: Mileage has already been reached (SQLSTATE P0001)", err.Error())
}

func TestSaveRefuelWithUnrealisticDate(t *testing.T) {
	assert := assert.New(t)

	// Setup
	var userId = 1
	var unrealisticTime, _ = time.Parse("2006-02-01T15:04:05", "1950-09-05T16:34:25")

	var refuelWithInvalidMileage = Refuel{
		Id:                  4,
		Description:         "TestRefuel1",
		DateTime:            unrealisticTime,
		PricePerLiterInEuro: 1.2,
		TotalAmount:         35,
		PricePerLiter:       40,
		Currency:            "Chf",
		Mileage:             390789,
		LicensePlate:        "KN-KN-9999",
		LastChanged:         time.Now(),
	}

	// When
	refuelId, err := saveRefuelByUserId(refuelWithInvalidMileage, userId)

	// Then
	assert.NotNil(err)
	assert.Equal(-1, refuelId)
	assert.Equal("ERROR: new row for relation \"refuel\" violates check constraint \"realistic_date\" (SQLSTATE 23514)", err.Error())
}
