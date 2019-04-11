/**
 */
package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/xeipuuv/gojsonschema"
)

var (
	getIndividualParticipantSchema = `
{
	"$id": "PreciousCargoShippping:getIndividualParticipantSchema",
	"type": "object",
	"properties": {
		"id": {
			"type": "string",
			"description": "ID of IndividualParticipant",
			"pattern": "^([0-9]{4,32})$"
		}
	},
	"required": [ "id" ]
}
`
	getIndividualParticipantSchemaLoader = gojsonschema.NewStringLoader(getIndividualParticipantSchema)
)

// Retrieves Participant data by Id, returns data structure
type getIndividualParticipantArg struct {
	Id string `json:"id"`
}

// Returns ID of shipment
type getIndividualParticipantResult struct {
	Participant IndividualParticipant `json:"participant"`
}

type getIndividualParticipantInvocation struct {
	arg getIndividualParticipantArg
	res getIndividualParticipantResult
}

func (inv *getIndividualParticipantInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter getIndividualParticipantInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("Expecting JSON input as first param")
	}

	result, err := gojsonschema.Validate(getIndividualParticipantSchemaLoader,
		gojsonschema.NewStringLoader(args[0]))
	if err != nil {
		logger.Println(err)
		return errors.New("Error parsing/validating JSON arg")
	}
	if !result.Valid() {
		logger.Printf("JSON input not valid:\n")
		for _, err := range result.Errors() {
			logger.Printf("- %s\n", err)
		}
		return errors.New("JSON not valid according to schema")
	}

	inv.arg = getIndividualParticipantArg{}
	err = json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("Error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}
	return nil
}

func (inv *getIndividualParticipantInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter getIndividualParticipant.process")
	logger.Printf("arg=%#v\n", inv.arg)

	ck, err := stub.CreateCompositeKey(ns, []string{".", "IndividualParticipant", "#", inv.arg.Id})
	if err != nil {
		logger.Println(err)
		return errors.New("internal error generating composite key (1).")
	}
	logger.Printf("key=%s\n", ck)

	data, err := stub.GetState(ck)
	if err != nil {
		logger.Println(err)
		return errors.New("internal error reading from world state.")
	}
	if data == nil {
		logger.Println("Nothing found for given key.")
		return errors.New("Not found")
	}

	// we could easily return data as-is (because its JSON), but to
	// make sure data is right we're trying to unmarshal it into right type
	err = json.Unmarshal(data, &inv.res.Participant)
	if err != nil {
		logger.Println(err)
		return errors.New("internal error reading from world state (2).")
	}
	logger.Printf("Found %#v\n", inv.res.Participant)

	return nil
}

func (inv *getIndividualParticipantInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
