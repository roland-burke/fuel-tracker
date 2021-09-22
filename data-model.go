package main

import "time"

type Configuration struct {
	AuthToken string `json:"authToken"`
	Port      int    `json:"port"`
	UrlPrefix string `json:"urlPrefix"`
}

type Refuel struct {
	Id                  int       `json:"id"`
	Name                string    `json:"name"`
	DateTime            time.Time `json:"dateTime"`
	PricePerLiterInEuro float64   `json:"pricePerLiterInEuro"`
	TotalAmount         float64   `json:"totalAmount"`
	PricePerLiter       float64   `json:"pricePerLiter"`
	Currency            string    `json:"currency"`
	LastChanged         time.Time `json:"lastChanged"`
}

type Deletion struct {
	Id int `json:"id"`
}
