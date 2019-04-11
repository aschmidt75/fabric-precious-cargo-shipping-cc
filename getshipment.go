/**
 */
package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/xeipuuv/gojsonschema"
	"reflect"
)

var (
	getShipmentSchema = `
 {
	 "$id": "PreciousCargoShippping:getShipmentSchema",
	 "type": "object",
	 "properties": {
		 "id": {
			 "type": "string",
			 "description": "ID of Shipment"
		 }
	 },
	 "required": [ "id" ]
 }
 `
)

// Retrieves Participant data by Id, returns data structure
type getShipmentArg struct {
	Id string `json:"id"`
}

// Returns ID of shipment
type getShipmentResult struct {
	Shipment Shipment `json:"participant"`
}

type getShipmentInvocation struct {
	arg getShipmentArg
	res getShipmentResult
}

func (inv *getShipmentInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter getShipmentInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("expecting JSON input as first param")
	}

	sl := gojsonschema.NewStringLoader(getShipmentSchema)
	dl := gojsonschema.NewStringLoader(args[0])

	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		return errors.New("error parsing/validating JSON arg")
	}
	if !result.Valid() {
		logger.Printf("JSON input not valid:\n")
		for _, err := range result.Errors() {
			logger.Printf("- %s\n", err)
		}
		return errors.New("json not valid according to schema")
	}

	inv.arg = getShipmentArg{}
	err = json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}
	return nil
}

func (inv *getShipmentInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter submitShipmentInvocation.process")
	logger.Printf("arg=%#v\n", inv.arg)

	r := &registry{
		typeStr: "Shipment",
		typeRT:  reflect.TypeOf(&Shipment{}),
	}

	_, x, err := r.get(stub, inv.arg.Id)
	inv.res = x.(getShipmentResult)

	return err
}

func (inv *getShipmentInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
