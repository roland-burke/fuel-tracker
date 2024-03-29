package model

import (
	"time"
)

type Configuration struct {
	Description string `json:"description"`
	ApiKey      string `json:"apiKey"`
	Port        int    `json:"port"`
	UrlPrefix   string `json:"urlPrefix"`
}

type Refuel struct {
	Id                  int       `json:"id"`
	Description         string    `json:"description"`
	DateTime            time.Time `json:"dateTime"`
	PricePerLiterInEuro float64   `json:"pricePerLiterInEuro"`
	TotalAmount         float64   `json:"totalAmount"`
	PricePerLiter       float64   `json:"pricePerLiter"`
	Currency            string    `json:"currency"`
	Mileage             int       `json:"mileage"`
	LicensePlate        string    `json:"licensePlate"`
	Trip                int       `json:"trip"`
	LastChanged         time.Time `json:"lastChanged"`
}

type Stat struct {
	Year    int     `json:"year"`
	Cost    float64 `json:"cost"`
	Mileage int     `json:"mileage"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefuelResponse struct {
	Refuels    []Refuel `json:"refuels"`
	TotalCount int      `json:"totalCount"`
}

type StatisticsResponse struct {
	Stats                   []Stat  `json:"stats"`
	TotalCost               float64 `json:"totalCost"`
	TotalMileage            int     `json:"totalMileage"`
	AverageCost             float64 `json:"avrgCost"`
	AverageMileagePerRefuel int     `json:"avrgDistancePerRefuel"`
}

type DefaultRequest struct {
	Payload []Refuel `json:"payload"`
}

type DeletionRequest struct {
	Id int `json:"id"`
}
