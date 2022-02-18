package main

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"
)

const REFUEL_TABLE_NAME = "refuel"
const MAX_RESPONSE_SIZE = 8

func getUserIdByName(username string) int {
	var user_id int
	err := conn.QueryRow(context.Background(), "SELECT users_id FROM users WHERE username=$1", username).Scan(&user_id)
	if err != nil {
		logger.Error("Cannot get user id: %s", err.Error())
		return -1
	}

	return user_id
}

func getCredentials(requestedUsername string) (error, string, string) {
	var username string
	var password string
	var err = conn.QueryRow(context.Background(), "SELECT username, pass_key FROM users WHERE username=$1", requestedUsername).Scan(&username, &password)
	if username == "" {
		return errors.New("Username " + requestedUsername + " does not exist"), "", ""
	}
	if err != nil {
		return err, "", ""
	}
	return nil, username, password
}

func deleteRefuelByUserId(refuelId int, userId int) error {
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM "+REFUEL_TABLE_NAME+" WHERE (id=$1 AND users_id=$2)", refuelId, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}

func updateRefuelByUserId(refuel Refuel, userId int) error {
	commandTag, err := conn.Exec(context.Background(), "UPDATE "+REFUEL_TABLE_NAME+" SET description=$1, date_time=$2, price_per_liter_euro=$3, total_liter=$4, price_per_liter=$5, currency=$6, mileage=$7, license_plate=$8 WHERE (id=$9 AND users_id=$10)", refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, refuel.LicensePlate, refuel.Id, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

func saveRefuelByUserId(refuel Refuel, userId int) (error, int) {
	lastInsertId := 0
	err := conn.QueryRow(context.Background(), "INSERT INTO "+REFUEL_TABLE_NAME+"(users_id, description, date_time, price_per_liter_euro, total_liter, price_per_liter, currency, mileage, license_plate) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", userId, refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, strings.ToUpper(refuel.LicensePlate)).Scan(&lastInsertId)

	if err != nil {
		return err, -1
	}
	if lastInsertId <= 0 {
		return errors.New("Insert was not successful"), -1
	}
	return nil, lastInsertId
}

func getStatisticsByUserId(userId int) (StatisticsResponse, error) {
	var err error
	var totalCost float64
	var totalMileage float64

	// Get total cost and mileage
	err = conn.QueryRow(context.Background(), "SELECT SUM (total_liter * price_per_liter_euro) AS cost FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1", userId).Scan(&totalCost)
	err = conn.QueryRow(context.Background(), "SELECT max(mileage) - min(mileage) FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1", userId).Scan(&totalMileage)

	var statListBuffer [100]Stat

	var index = 0

	// Get cost and mileage per year
	rows, err := conn.Query(context.Background(), "SELECT date_part('year', date_time) AS year, SUM (total_liter * price_per_liter_euro) AS cost, max(mileage) - min(mileage) AS mileage FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 GROUP BY year ORDER BY year DESC;", userId)
	if err != nil {
		logger.Error("Getting all reufels failed: %s", err.Error())
		return StatisticsResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var cost float64
		var mileage float64
		var year int

		err := rows.Scan(&year, &cost, &mileage)
		if err != nil {
			logger.Error("Scanning single row failed: %s", err.Error())
			return StatisticsResponse{}, err
		}

		statListBuffer[index] = Stat{
			Year:    year,
			Cost:    cost,
			Mileage: mileage,
		}
		index += 1
	}

	if err != nil {
		logger.Error("Getting statistics failed: %s", err.Error())
		return StatisticsResponse{}, err
	}

	response := StatisticsResponse{
		Stats:        statListBuffer[:index],
		TotalCost:    math.Round(totalCost*100) / 100,
		TotalMileage: math.Round(totalMileage*100) / 100, // Round to the 2. decimal place
	}

	return response, err
}

func getAllRefuelsByUserId(userId int, startIndex int, licensePlate string, month int, year int) (RefuelResponse, error) {
	var err error = nil
	rows, err := conn.Query(context.Background(), "SELECT * FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'DEFAULT') OR (license_plate = $2)) AND (($3 = 0) OR (date_part('month', date_time) = $3)) AND (($4 = 0) OR (date_part('year', date_time) = $4)) ORDER BY date_time DESC", userId, licensePlate, month, year)
	if err != nil {
		logger.Error("Getting all reufels failed: %s", err.Error())
		return RefuelResponse{}, err
	}

	defer rows.Close()

	var refuelListBuffer [MAX_RESPONSE_SIZE]Refuel

	var index = 0
	var counter = 0
	var validCounter = 0

	for rows.Next() {
		var id int = -1
		var users_id int
		var description string
		var dateTime time.Time
		var pricePerLiterInEuro float64
		var totalAmount float64
		var pricePerLiter float64
		var currency string
		var mileage float64
		var licensePlate string
		var lastChanged time.Time

		err := rows.Scan(&id, &users_id, &description, &dateTime, &pricePerLiterInEuro, &totalAmount, &pricePerLiter, &currency, &mileage, &licensePlate, &lastChanged)
		if err != nil {
			logger.Error("Scanning single row failed: %s", err.Error())
			return RefuelResponse{}, err
		}

		if counter >= startIndex && counter < startIndex+MAX_RESPONSE_SIZE {
			if id != -1 {
				validCounter += 1
			}
			refuelListBuffer[index] = Refuel{
				Id:                  id,
				Description:         description,
				DateTime:            dateTime,
				PricePerLiterInEuro: pricePerLiterInEuro,
				TotalAmount:         totalAmount,
				PricePerLiter:       pricePerLiter,
				Currency:            currency,
				Mileage:             mileage,
				LicensePlate:        licensePlate,
				LastChanged:         lastChanged,
			}
			index += 1
		}
		counter += 1
	}

	var response RefuelResponse

	if validCounter < MAX_RESPONSE_SIZE {
		response = RefuelResponse{
			Refuels:    refuelListBuffer[:validCounter],
			TotalCount: counter,
		}
	} else {
		response = RefuelResponse{
			Refuels:    refuelListBuffer[:MAX_RESPONSE_SIZE],
			TotalCount: counter,
		}
	}

	return response, err
}
