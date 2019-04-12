// fabric-precious-cargo-shipping-cc is a sample chaincode for Hyperledger Fabric
// Copyright (C) 2019 @aschmidt75
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// InvocationHandler is a generic interface for wrapping a
// single chaincode transaction.
type InvocationHandler interface {
	/**
	 * extracts call arguments from the stub and checks them
	 */
	checkParseArguments(stub shim.ChaincodeStubInterface) error

	/**
	 * runs the chaincode. Returns nil if successful, an error
	 * otherwise.
	 */
	process(stub shim.ChaincodeStubInterface) error

	/**
	 * returns result
	 */
	getResponse(stub shim.ChaincodeStubInterface) interface{}
}

func newID(stub shim.ChaincodeStubInterface, indexName string) (string, error) {
	ckIndex, err := stub.CreateCompositeKey(ns, []string{".", indexName, ".", "index"})
	if err != nil {
		return "", err
	}

	var lastIndex uint64
	lastIndexBytes, err := stub.GetState(ckIndex)
	if err != nil {
		return "", err
	}
	if lastIndexBytes != nil {
		lastIndex, err = strconv.ParseUint(string(lastIndexBytes), 10, 64)
		if err != nil {
			return "", err
		}
	}
	lastIndex = lastIndex + 1

	err = stub.PutState(ckIndex, []byte(strconv.FormatUint(lastIndex, 10)))
	if err != nil {
		return "", err
	}
	logger.Printf("newId: PutState to key=%s, index=%v\n", ckIndex, lastIndex)

	return fmt.Sprintf("%010v", lastIndex), nil

}

func getGenericKey(stub shim.ChaincodeStubInterface, typeStr string, id string) (string, error) {
	ck, err := stub.CreateCompositeKey(ns, []string{".", typeStr, "#", id})
	if err != nil {
		logger.Println(err)
		return "", err
	}
	return ck, nil
}

func getIndividualParticipantKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	return getGenericKey(stub, "IndividualParticipant", id)
}

func getShipmentCoKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	return getGenericKey(stub, "ShipmentCo", id)
}

func getShipmentKey(stub shim.ChaincodeStubInterface, id string) (string, error) {
	return getGenericKey(stub, "Shipment", id)
}

func getIndividualParticipant(stub shim.ChaincodeStubInterface, id string) (string, *IndividualParticipant, error) {
	ck, err := getIndividualParticipantKey(stub, id)
	if err != nil {
		return "", nil, err
	}
	data, err := stub.GetState(ck)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (1)")
	}
	res := &IndividualParticipant{}
	err = json.Unmarshal(data, res)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (2)")
	}
	logger.Printf("Found %#v\n", res)

	return ck, res, nil
}

func getShipmentCo(stub shim.ChaincodeStubInterface, id string) (string, *ShipmentCo, error) {
	ck, err := getShipmentCoKey(stub, id)
	if err != nil {
		return "", nil, err
	}
	data, err := stub.GetState(ck)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (1)")
	}
	res := &ShipmentCo{}
	err = json.Unmarshal(data, res)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (2)")
	}
	logger.Printf("Found %#v\n", res)

	return ck, res, nil
}
