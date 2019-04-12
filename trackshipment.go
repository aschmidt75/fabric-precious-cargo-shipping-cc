// fabric-precious-cargo-shipping-cc is a sample chaincode for Hyperledger Fabric
// Copyright (C) 2019 @aschmidt75
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"time"
)

// Retrieves Participant data by Id, returns data structure
type trackShipmentArg struct {
	ID          string  `json:"id"`
	At          string  `json:"at"` // time in RFC3339, e.g. 2006-01-02T15:04:05Z
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
	Temperature float32 `json:"temp"`
	Humidity    float32 `json:"hum"` // in [%]
}

type trackShipmentInvocation struct {
	arg trackShipmentArg

	at          time.Time
	shipmentKey string
	shipment    Shipment

	res struct{}
}

func (inv *trackShipmentInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter trackShipmentInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("expecting JSON input as first param")
	}

	var err error

	inv.arg = trackShipmentArg{}
	err = json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}

	// parse time
	inv.at, err = time.Parse(time.RFC3339, inv.arg.At)
	if err != nil {
		return errors.New("invalid at argument: Not parseable, please provide in RFC3339, e.g. 2006-01-02T15:04:05Z")
	}
	// must be somewhat recent. (TODO)

	// check humidity
	if inv.arg.Humidity < 0 || inv.arg.Humidity > 100 {
		return errors.New("invalid hum argument: Must be [0..100] [%]")
	}

	// load shipment
	y, x, err := shipmentRegistry().get(stub, inv.arg.ID)
	if err != nil {
		println(err)
		return errors.New("unable to locate shipment for this ID")
	}
	inv.shipmentKey = y
	inv.shipment = x.(Shipment)

	return nil
}

func (inv *trackShipmentInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter submitShipmentInvocation.process")
	logger.Printf("arg=%#v\n", inv.arg)

	// use a generic registry but with a special type string specific to
	// the shipmennt
	r := &registry{
		typeStr: fmt.Sprintf("trackingDataPoint[%s]", inv.shipment.ID.ID),
		typeRT:  reflect.TypeOf(&TrackingDataPoint{}),
	}

	tdp := TrackingDataPoint{
		ShipmentID:  ID{inv.arg.ID},
		At:          inv.at,
		Latitude:    inv.arg.Latitude,
		Longitude:   inv.arg.Longitude,
		Temperature: inv.arg.Temperature,
		Humidity:    inv.arg.Humidity,
	}

	key, err := r.create(stub, tdp)
	if err != nil {
		logger.Println(err)
		return errors.New("unable to write trackment data")
	}
	logger.Printf("Tracked: %s\n", key)

	return err
}

func (inv *trackShipmentInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
