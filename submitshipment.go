/**
 */
package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// creates a new shipment structure from given Ids of shipper and Participants.
type submitShipmentArg struct {
	Shipper string `json:"by"`
	From    string `json:"from"`
	To      string `json:"to"`
}

// Returns ID of shipment
type submitShipmentResult struct {
	Id string `json:"id"`
}

type submitShipmentInvocation struct {
	arg submitShipmentArg
	res submitShipmentResult
}

func (inv *submitShipmentInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter submitShipmentInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("Expecting JSON input as first param")
	}

	inv.arg = submitShipmentArg{}
	err := json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("Error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}
	return nil
}

func (inv *submitShipmentInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter submitShipmentInvocation.process")

	inv.res = submitShipmentResult{}

	return nil
}

func (inv *submitShipmentInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
