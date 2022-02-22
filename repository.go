package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

const REFUEL_TABLE_NAME = "refuel"
const MAX_RESPONSE_SIZE = 8
const MAX_YEARS_FOR_STATS = 100

func getUserIdByCredentials(username string, password string) int {
	var user_id int
	err := conn.QueryRow(context.Background(), "SELECT users_id FROM users WHERE (username=$1 AND pass_key=$2)", username, password).Scan(&user_id)
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
	if err != nil {
		return err, "", ""
	}
	if username == "" {
		return errors.New("Username " + requestedUsername + " does not exist"), "", ""
	}
	return nil, username, password
}

func saveRefuelByUserId(refuel Refuel, userId int) (int, error) {
	lastInsertId := 0
	err := conn.QueryRow(context.Background(), "INSERT INTO "+REFUEL_TABLE_NAME+"(users_id, description, date_time, price_per_liter_euro, total_liter, price_per_liter, currency, mileage, license_plate) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", userId, refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, strings.ToUpper(refuel.LicensePlate)).Scan(&lastInsertId)

	if err != nil {
		return -1, err
	}
	if lastInsertId <= 0 {
		return -1, errors.New("Insert was not successful")
	}

	return lastInsertId, nil
}

func updateRefuelByUserId(refuel Refuel, userId int) error {
	commandTag, err := conn.Exec(context.Background(), "UPDATE "+REFUEL_TABLE_NAME+" SET description=$1, date_time=$2, price_per_liter_euro=$3, total_liter=$4, price_per_liter=$5, currency=$6, mileage=$7, license_plate=$8 WHERE (id=$9 AND users_id=$10)", refuel.Description, refuel.DateTime, refuel.PricePerLiterInEuro, refuel.TotalAmount, refuel.PricePerLiter, refuel.Currency, refuel.Mileage, refuel.LicensePlate, refuel.Id, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New(fmt.Sprintf("No row found to update, refuelId: %d ,userId: %d", refuel.Id, userId))
	}
	return nil
}

func deleteRefuelByUserId(refuelId int, userId int) error {
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM "+REFUEL_TABLE_NAME+" WHERE (id=$1 AND users_id=$2)", refuelId, userId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New(fmt.Sprintf("No row found to delete, refuelId: %d ,userId: %d", refuelId, userId))
	}
	return nil
}

func getAllRefuelsByUserId(userId int, startIndex int, licensePlate string, month int, year int) (RefuelResponse, error) {
	var err error = nil
	rows, err := conn.Query(context.Background(), "SELECT * FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'ALL') OR (license_plate = $2)) AND (($3 = 0) OR (date_part('month', date_time) = $3)) AND (($4 = 0) OR (date_part('year', date_time) = $4)) ORDER BY date_time DESC", userId, licensePlate, month, year)
	if err != nil {
		logger.Error("Getting all reufels failed: %s", err.Error())
		rows.Close()
		return RefuelResponse{}, err
	}

	defer rows.Close()

	var refuelListBuffer [MAX_RESPONSE_SIZE]Refuel

	var responseListIndex = 0
	var totalCounter = 0
	var validCounter = 0

	var tripMap = getTripsByLicensePlate(userId, licensePlate)

	for rows.Next() {
		var id int = -1
		var users_id int
		var description string
		var dateTime time.Time
		var pricePerLiterInEuro float64
		var totalAmount float64
		var pricePerLiter float64
		var currency string
		var mileage int
		var licensePlate string
		var lastChanged time.Time

		err := rows.Scan(&id, &users_id, &description, &dateTime, &pricePerLiterInEuro, &totalAmount, &pricePerLiter, &currency, &mileage, &licensePlate, &lastChanged)
		if err != nil {
			logger.Error("Scanning single row failed: %s", err.Error())
			rows.Close()
			return RefuelResponse{}, err
		}

		logger.Debug("data:")
		logger.Info("id: %d", id)
		logger.Info("users_id: %d", users_id)
		logger.Info("license plate: %s", licensePlate)
		logger.Info("mileage: %d", mileage)
		logger.Debug("\n")

		if totalCounter >= startIndex && totalCounter < startIndex+MAX_RESPONSE_SIZE {
			if id != -1 {
				validCounter += 1
			}
			refuelListBuffer[responseListIndex] = Refuel{
				Id:                  id,
				Description:         description,
				DateTime:            dateTime,
				PricePerLiterInEuro: pricePerLiterInEuro,
				TotalAmount:         totalAmount,
				PricePerLiter:       pricePerLiter,
				Currency:            currency,
				Mileage:             mileage,
				LicensePlate:        licensePlate,
				Trip:                tripMap[id],
				LastChanged:         lastChanged,
			}
			responseListIndex += 1
		}
		totalCounter += 1
	}

	var response RefuelResponse

	if validCounter < MAX_RESPONSE_SIZE {
		response = RefuelResponse{
			Refuels:    refuelListBuffer[:validCounter],
			TotalCount: totalCounter,
		}
	} else {
		response = RefuelResponse{
			Refuels:    refuelListBuffer[:MAX_RESPONSE_SIZE],
			TotalCount: totalCounter,
		}
	}

	return response, err
}

// Get cost and mileage per year
func getStatsForEveryYear(userId int, licensePlate string) ([MAX_YEARS_FOR_STATS]Stat, int) {
	var index = 0
	var statListBuffer [MAX_YEARS_FOR_STATS]Stat

	rows, err := conn.Query(context.Background(), "SELECT date_part('year', date_time) AS year, SUM (total_liter * price_per_liter_euro) AS cost, max(mileage) - min(mileage) AS mileage FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'ALL') OR (license_plate = $2)) GROUP BY year ORDER BY year DESC;", userId, licensePlate)
	if err != nil {
		logger.Error("Getting all sats for every year failed: %s", err.Error())
		rows.Close()
		return [MAX_YEARS_FOR_STATS]Stat{}, 0
	}

	defer rows.Close()

	for rows.Next() {
		var cost float64
		var mileage int
		var year int

		err := rows.Scan(&year, &cost, &mileage)
		if err != nil {
			logger.Error("Scanning year, cost, mileage failed: %s", err.Error())
			rows.Close()
			return [MAX_YEARS_FOR_STATS]Stat{}, 0
		}

		statListBuffer[index] = Stat{
			Year:    year,
			Cost:    cost,
			Mileage: mileage,
		}
		index += 1
	}
	return statListBuffer, index
}

func convertStringArrayToIntArray(input []string) []int {
	var output = []int{}
	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			logger.Error("Cannot convert: %s to int: %s", i, err.Error())
			return []int{}
		}
		output = append(output, j)
	}
	return output
}

func calculateAverageTrip(inputMap map[int]int) int {
	var avrg int
	var size = len(inputMap)

	if size <= 0 {
		logger.Error("Input map is empty")
		return 0
	}

	// calculate average mileage difference
	for _, trip := range inputMap {
		avrg += trip
	}

	return avrg / size
}

func getTripsByLicensePlate(userId int, licensePlate string) map[int]int {
	var tripMap = make(map[int]int)

	// Get all mileages per license plate
	mileageRows, err := conn.Query(context.Background(), "select license_plate, STRING_AGG(mileage::varchar(10), ',' order by date_time desc) as mileages, STRING_AGG(id::varchar(10), ',' order by date_time desc) as ids FROM (select distinct on (mileage) * from "+REFUEL_TABLE_NAME+" order by mileage, date_time desc) s WHERE users_id=$1 group by license_plate;", userId)
	if err != nil {
		logger.Error("Getting all mileages failed: %s", err.Error())
		mileageRows.Close()
		return tripMap
	}

	defer mileageRows.Close()

	for mileageRows.Next() {
		values, err := mileageRows.Values()
		if err != nil {
			logger.Error("Failed to get mileage for license plate: %s", err.Error())
		}

		var idList = convertStringArrayToIntArray(strings.Split(values[2].(string), ","))

		logger.Debug("values: ")
		logger.Debug("%s", values[0].(string))
		logger.Debug("%s", values[1].(string))
		logger.Debug("%s", values[2].(string))
		logger.Debug("\n")

		// At Index 1 should be the mileage list e.g: 7300,6700,6200,...
		var stringList = strings.Split(values[1].(string), ",")
		logger.Debug("String list: %s", stringList)

		// reverse
		for i, j := 0, len(stringList)-1; i < j; i, j = i+1, j-1 {
			stringList[i], stringList[j] = stringList[j], stringList[i]
		}

		// reverse
		for i, j := 0, len(idList)-1; i < j; i, j = i+1, j-1 {
			idList[i], idList[j] = idList[j], idList[i]
		}

		logger.Debug("String list length: %d", len(stringList))
		// Get all trips for one license plate
		if len(stringList) >= 2 {
			var intList = convertStringArrayToIntArray(stringList)

			for i := 0; i < len(intList)-1; i++ {
				var trip = intList[i+1] - intList[i]
				logger.Debug("Trip: %d", trip)
				tripMap[idList[i+1]] = trip
			}
		}
	}

	return tripMap
}

func getAverageDistancePerRefuel(userId int, licensePlate string) int {
	var mileageRows pgx.Rows
	var totalEntries = 0

	if strings.Compare(strings.ToUpper(licensePlate), "ALL") == 0 {
		var allTrips = getTripsByLicensePlate(userId, licensePlate)
		return calculateAverageTrip(allTrips)

	} else {
		err := conn.QueryRow(context.Background(), "SELECT COUNT(*) AS totalEntries FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND license_plate = $2", userId, licensePlate).Scan(&totalEntries)

		if err != nil {
			logger.Error("Getting totalEntries failed: %s", err.Error())
			return 0
		}

		mileageRows, err = conn.Query(context.Background(), "SELECT mileage FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND license_plate = $2 ORDER BY date_time DESC", userId, licensePlate)
		if err != nil {
			logger.Error("Getting all distances failed: %s", err.Error())
			mileageRows.Close()
			return 0
		}
	}

	defer mileageRows.Close()

	// get array of all mileages
	var allMileages []int
	for mileageRows.Next() {
		var mileage int = 0
		err := mileageRows.Scan(&mileage)
		if err != nil {
			logger.Error("Getting mileage value failed: %s", err.Error())
			return 0
		}
		allMileages = append(allMileages, mileage)
	}

	// reverse
	for i, j := 0, len(allMileages)-1; i < j; i, j = i+1, j-1 {
		allMileages[i], allMileages[j] = allMileages[j], allMileages[i]
	}

	if totalEntries <= 1 {
		mileageRows.Close()
		return 0
	}

	var avrg int

	// calculate average mileage difference
	for i := 0; i < totalEntries-1; i++ {
		var trip = allMileages[i+1] - allMileages[i]
		avrg += trip
	}
	// because we don't have the first value
	return avrg / (totalEntries - 1)
}

func getStatisticsByUserId(userId int, licensePlate string) StatisticsResponse {
	var totalCost float64 = -1.0
	var totalMileage int = -1
	var avrgCost float64 = -1.0

	// Get total cost, mileage and average cost
	err := conn.QueryRow(context.Background(), "SELECT SUM (total_liter * price_per_liter_euro) AS cost FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'ALL') OR (license_plate = $2))", userId, licensePlate).Scan(&totalCost)
	if err != nil {
		logger.Error("Failed to get total cost: %s", err.Error())
	}

	err = conn.QueryRow(context.Background(), "SELECT MAX(mileage) - MIN(mileage) FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'ALL') OR (license_plate = $2))", userId, licensePlate).Scan(&totalMileage)
	if err != nil {
		logger.Error("Failed to get total mileage: %s", err.Error())
	}

	err = conn.QueryRow(context.Background(), "SELECT AVG(price_per_liter_euro) FROM "+REFUEL_TABLE_NAME+" WHERE users_id=$1 AND (($2 = 'ALL') OR (license_plate = $2))", userId, licensePlate).Scan(&avrgCost)
	if err != nil {
		logger.Error("Failed to get average cost: %s", err.Error())
	}

	//
	err = nil

	var avrgDistancePerRefuel = getAverageDistancePerRefuel(userId, licensePlate)
	statListBuffer, amount := getStatsForEveryYear(userId, licensePlate)

	response := StatisticsResponse{
		Stats:                   statListBuffer[:amount],
		TotalCost:               math.Round(totalCost*100) / 100,
		TotalMileage:            totalMileage, // Round to the 2. decimal place
		AverageMileagePerRefuel: avrgDistancePerRefuel,
		AverageCost:             math.Round(avrgCost*1000) / 1000, // Round to the 3. decimal place
	}

	return response
}
