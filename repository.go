package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const REFUEL_TABLE_NAME = "refuel"

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
		log.Println("ERROR - Saving reufel failed:", err)
		return false
	}
	return true
}

func getAllRefuelsByUserId(userId int) (RefuelResponse, error) {
	var err error = nil
	rows, err := conn.Query(context.Background(), "SELECT * FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 ORDER BY date_time DESC", userId)
	if err != nil {
		log.Println("ERROR - Getting all reufels failed:", err)
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
			fmt.Fprintf(os.Stderr, "row next failed: %v\n", err)
			os.Exit(1)
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
