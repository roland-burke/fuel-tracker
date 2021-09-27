package main

import (
	"time"
)

type Configuration struct {
	AuthToken string `json:"authToken"`
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

type RefuelResposne struct {
	Refuels []Refuel `json:"refuels"`
}

type Deletion struct {
	Id int `json:"id"`
}
