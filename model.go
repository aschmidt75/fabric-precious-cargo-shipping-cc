/**
 */
package main

import (
	"time"
)

// World State/Ledger data structs

// ID is a generic identifier
type ID struct {
	ID string `json:"id"`
}

// Participant is a simple Participant identified by Id and a name
type Participant struct {
	ID
	Name string `json:"name"`
}

// IndividualParticipant has an address
type IndividualParticipant struct {
	Participant
	Address string `json:"address"` // simplified address field as one-liner
}

// ShipmentCo is a Shipment Company
type ShipmentCo struct {
	Participant
	Address string `json:"address"`
}

// Asset is identified by Id
type Asset struct {
	ID
}

// Shipment combines Shipper, From and To Participants and
// Status
type Shipment struct {
	Asset

	ShipperID string `json:"by"`
	FromID    string `json:"from"`
	ToID      string `json:"to"`

	Status string `json:"status"`

	SubmittedAt time.Time `json:"submittime"`
	DelivererAt time.Time `json:"delivertime,omitempty"`
}

// TrackingDataPoint combines a location and environmental
// parameters for a shipment, at a point in time.
type TrackingDataPoint struct {
	ShipmentID  ID        `json:"shipmentId"`
	At          time.Time `json:"at"`
	Latitude    float64   `json:"lat"`
	Longitude   float64   `json:"lng"`
	Temperature float32   `json:"temp"`
	Humidity    float32   `json:"hum"`
}
