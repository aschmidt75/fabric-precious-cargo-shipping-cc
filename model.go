/**
 */
package main

import (
	"time"
)

// World State/Ledger data structs

type Id struct {
	Id string `json:"id"`
}

type Participant struct {
	Id
	Name string `json:"name"`
}

type IndividualParticipant struct {
	Participant
	Address string `json:"address"` // simplified address field as one-liner
}

type ShipmentCo struct {
	Participant
	Address string `json:"address"`
}

type Asset struct {
	Id
}

type Shipment struct {
	Asset

	ShipperId string `json:"by"`
	FromId    string `json:"from"`
	ToId      string `json:"to"`

	Status string `json:"status"`

	SubmittedAt time.Time `json:"submittime"`
	DelivererAt time.Time `json:"delivertime,omitempty"`
}

type TrackingDataPoint struct {
	ShipmentID  Id        `json:"shipmentId"`
	At          time.Time `json:"at"`
	Latitude    float64   `json:"lat"`
	Longitude   float64   `json:"lng"`
	Temperature float32   `json:"temp"`
	Humidity    float32   `json:"hum"`
}
