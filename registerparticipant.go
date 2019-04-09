/**
 */
package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Invocation struct to register an IndividualParticipant
type registerIndividualParticipantInvocation struct {
	// input arguments (from client)
	arg registerIndividualParticipantArg

	// temporary stuff: id for new Participant
	idStr string

	// result (to client)
	res registerIndividualParticipantResult
}

// Creates a new Participant, by name and address. Returns the Id
type registerIndividualParticipantArg struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Returns ID of shipment
type registerIndividualParticipantResult struct {
	Id string `json:"id"`
}

// Unmarshal input argument, optionally check them
func (inv *registerIndividualParticipantInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter registerIndividualParticipantInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("Expecting JSON input as first param")
	}

	inv.arg = registerIndividualParticipantArg{}
	err := json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("Error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}

	// programmatic input check
	if len(inv.arg.Name) < 3 {
		return errors.New("Invalid input, name is too short")
	}
	if len(inv.arg.Name) > 100 {
		return errors.New("Invalid input, name is too long")
	}

	return nil
}

// Processes the invocation. Updates the invocation struct. Returns an error or nil if successful.
func (inv *registerIndividualParticipantInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter registerIndividualParticipant.process")
	logger.Printf("arg=%#v\n", inv.arg)

	// Create an ID for the new participant
	s, err := newId(stub, "IndividualParticipant")
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error generating index key.")
	}
	inv.idStr = s

	// create data item for world state update
	p := IndividualParticipant{
		Participant: Participant{
			Id: Id{
				Id: inv.idStr,
			},
			Name: inv.arg.Name,
		},
		Address: inv.arg.Address,
	}

	// combine namespace, type and ID into a key
	ck, err := stub.CreateCompositeKey(ns, []string{".", "IndividualParticipant", "#", inv.idStr})
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error generating composite key (2).")
	}
	logger.Printf("key=%s\n", ck)

	// marshal data to json and ...
	data, err := json.Marshal(p)
	if err != nil {
		logger.Println(err)
		return errors.New("Internal JSON marshal error (1).")
	}
	// ... save to world state
	err = stub.PutState(ck, []byte(data))
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error writing world state.")
	}
	logger.Printf("PutState to key=%s, data=%#v\n", ck, p)

	// return struct to client contains ID
	inv.res = registerIndividualParticipantResult{
		Id: inv.idStr,
	}

	return nil
}

func (inv *registerIndividualParticipantInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
