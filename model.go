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
}

type Asset struct {
	Id
}

type Shipment struct {
	Asset

	Shipper ShipmentCo            `json:"by"`
	From    IndividualParticipant `json:"from"`
	To      IndividualParticipant `json:"to"`

	Status string `json:"status"` // one of: submitted, intransit, delivered

	SubmittedAt time.Time `json:"submittime"`
	DelivererAt time.Time `json:"delivertime"`
}

type TrackingDataPoint struct {
	Asset

	ShipmentID  Id      `json:"shipmentId"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
	Temperature float32 `json:"temp"`
	Humidity    float32 `json:"hum"`
}
