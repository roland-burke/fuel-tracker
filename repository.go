package main

import (
	"context"
	"log"
	"strings"
	"time"
)

const REFUEL_TABLE_NAME = "refuel"

func getUserIdByName(username string) int {
	var user_id int
	err := conn.QueryRow(context.Background(), "SELECT users_id FROM users WHERE username=$1", username).Scan(&user_id)
	if err != nil {
		log.Println("ERROR - Cannot get user id", err)
		return -1
	}

	return user_id
}

func deleteRefuelByUserId(refuelId int, userId int) bool {
	_, err := conn.Exec(context.Background(), "DELETE FROM "+REFUEL_TABLE_NAME+" WHERE (id=$1 AND users_id=$2)", refuelId, userId)
	if err != nil {
		log.Println("ERROR - Deleting reufel failed:", err)
		return false
	}
	return true
}

func updateRefuelByUserId(refuel *Refuel, userId int) bool {
	_, err := conn.Exec(context.Background(), "UPDATE "+REFUEL_TABLE_NAME+" SET description=$1, date_time=$2, price_per_liter_euro=$3, total_liter=$4, price_per_liter=$5, currency=$6, mileage=$7, license_plate=$8 where (id=$9 AND users_id=$10)", refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, refuel.LicensePlate, refuel.Id, userId)
	if err != nil {
		log.Println("ERROR - Updating reufel failed:", err)
		return false
	}
	return true
}

func saveRefuelByUserId(refuel *Refuel, userId int) bool {
	_, err := conn.Exec(context.Background(), "INSERT INTO "+REFUEL_TABLE_NAME+"(users_id, description, date_time, price_per_liter_euro, total_liter, price_per_liter, currency, mileage, license_plate) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", userId, refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, strings.ToUpper(refuel.LicensePlate))
	if err != nil {
		log.Println("ERROR - Saving refuel failed:", err)
		return false
	}
	return true
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
		log.Println("ERROR - Getting all reufels failed:", err)
		return StatisticsResponse{}, err
	}

	for rows.Next() {
		var cost float64
		var mileage float64
		var year int

		err := rows.Scan(&year, &cost, &mileage)
		if err != nil {
			log.Println("ERROR - scan single row failed:", err)
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
		log.Println("ERROR - Getting statistics failed:", err)
		return StatisticsResponse{}, err
	}

	response := StatisticsResponse{
		Stats:        statListBuffer[:index],
		TotalCost:    totalCost,
		TotalMileage: totalMileage,
	}

	return response, err
}

func getAllRefuelsByUserId(userId int) (RefuelResponse, error) {
	var err error = nil
	rows, err := conn.Query(context.Background(), "SELECT * FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 ORDER BY date_time DESC", userId)
	if err != nil {
		log.Println("ERROR - Getting all reufels failed:", err)
		return RefuelResponse{}, err
	}

	var refuelListBuffer [100]Refuel

	var index = 0

	for rows.Next() {
		var id int
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
			log.Println("ERROR - scan single row failed:", err)
			return RefuelResponse{}, err
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

	response := RefuelResponse{
		Refuels: refuelListBuffer[:index],
	}

	return response, err
}
