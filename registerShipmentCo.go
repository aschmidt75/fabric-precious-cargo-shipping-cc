/**
 */
package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Creates a new Participant, by name and address. Returns the Id
type registerShipmentCoArg struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Returns ID of shipment
type registerShipmentCoResult struct {
	Id string `json:"id"`
}

type registerShipmentCoInvocation struct {
	arg registerShipmentCoArg
	res registerShipmentCoResult
}

func (inv *registerShipmentCoInvocation) checkParseArguments(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter registerShipmentCoInvocation.checkParseArguments")

	_, args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return errors.New("Expecting JSON input as first param")
	}

	inv.arg = registerShipmentCoArg{}
	err := json.Unmarshal([]byte(args[0]), &inv.arg)
	if err != nil {
		logger.Printf("Error unmarshaling JSON: %s", err)
		return errors.New("Invalid JSON")
	}
	return nil
}

func (inv *registerShipmentCoInvocation) process(stub shim.ChaincodeStubInterface) error {
	logger.Println("enter registerShipmentCo.process")
	logger.Printf("arg=%#v\n", inv.arg)

	idStr, err := newId(stub, "ShipmentCo")
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error generating index key.")
	}

	p := ShipmentCo{
		Participant: Participant{
			Id: Id{
				Id: idStr,
			},
			Name: inv.arg.Name,
		},
		Address: inv.arg.Address,
	}

	ck, err := getShipmentCoKey(stub, idStr)
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error generating composite key (2).")
	}
	logger.Printf("key=%s\n", ck)

	data, err := json.Marshal(p)
	if err != nil {
		logger.Println(err)
		return errors.New("Internal JSON marshal error (1).")
	}
	err = stub.PutState(ck, []byte(data))
	if err != nil {
		logger.Println(err)
		return errors.New("Internal error writing world state.")
	}
	logger.Printf("PutState to key=%s, data=%#v\n", ck, p)

	inv.res = registerShipmentCoResult{
		Id: idStr,
	}

	return nil
}

func (inv *registerShipmentCoInvocation) getResponse(stub shim.ChaincodeStubInterface) interface{} {
	return inv.res
}
