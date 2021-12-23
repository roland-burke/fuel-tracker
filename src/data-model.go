package main

import (
	"time"
)

type Configuration struct {
	ApiKey    string `json:"apiKey"`
	Port      int    `json:"port"`
	UrlPrefix string `json:"urlPrefix"`
}

type Refuel struct {
	Id                  int       `json:"id"`
	Description         string    `json:"description"`
	DateTime            time.Time `json:"dateTime"`
	PricePerLiterInEuro float64   `json:"pricePerLiterInEuro"`
	TotalAmount         float64   `json:"totalAmount"`
	PricePerLiter       float64   `json:"pricePerLiter"`
	Currency            string    `json:"currency"`
	Mileage             float64   `json:"mileage"`
	LicensePlate        string    `json:"licensePlate"`
	LastChanged         time.Time `json:"lastChanged"`
}

type Stat struct {
	Year    int     `json:"year"`
	Cost    float64 `json:"cost"`
	Mileage float64 `json:"mileage"`
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
	Stats        []Stat  `json:"stats"`
	TotalCost    float64 `json:"totalCost"`
	TotalMileage float64 `json:"totalMileage"`
}

type DefaultRequest struct {
	Payload []Refuel `json:"payload"`
}

type DeletionRequest struct {
	Id int `json:"id"`
}
