// fabric-precious-cargo-shipping-cc is a sample chaincode for Hyperledger Fabric
// Copyright (C) 2019 @aschmidt75
package main

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// registryInterface can be used to create and retrieve data items
// from the world state. It works on generic interfaces that can
// be marshalled from/to a storage format (e.g. JSON)
type registryInterface interface {
	// key creates a composite key from an id
	key(stub shim.ChaincodeStubInterface, id string) (string, error)

	// creates a new data item by marshaling into JSON. Returns
	// ID of newly created item.
	create(stub shim.ChaincodeStubInterface, item interface{}) (string, error)

	// get retrieves an item by its ID
	get(stub shim.ChaincodeStubInterface, id string) (string, interface{}, error)
}

// registry is a concrete registry with a type, given by its name (for creating keys)
// and its reflect.Type (for creating structs dynamically)
type registry struct {
	typeStr string
	typeRT  reflect.Type
}

func (r registry) key(stub shim.ChaincodeStubInterface, id string) (string, error) {
	ck, err := stub.CreateCompositeKey(ns, []string{".", r.typeStr, "#", id})
	if err != nil {
		logger.Println(err)
		return "", err
	}
	return ck, nil

}

func (r registry) create(stub shim.ChaincodeStubInterface, item interface{}) (string, error) {
	idStr, err := newID(stub, r.typeStr)
	if err != nil {
		logger.Println(err)
		return "", errors.New("internal error generating index key")
	}
	ck, err := getShipmentKey(stub, idStr)
	if err != nil {
		logger.Println(err)
		return "", errors.New("internal error generating composite key")
	}
	logger.Printf("key=%s\n", ck)

	data, err := json.Marshal(&item)
	if err != nil {
		logger.Println(err)
		return "", errors.New("internal JSON marshal error")
	}
	err = stub.PutState(ck, []byte(data))
	if err != nil {
		logger.Println(err)
		return "", errors.New("internal error writing world state")
	}
	logger.Printf("PutState to key=%s, data=%#v\n", ck, item)

	return ck, nil
}

func (r registry) get(stub shim.ChaincodeStubInterface, id string) (string, interface{}, error) {
	ck, err := r.key(stub, id)
	if err != nil {
		return "", nil, err
	}
	data, err := stub.GetState(ck)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (1)")
	}
	res := reflect.New(r.typeRT)
	err = json.Unmarshal(data, &res)
	if err != nil {
		logger.Println(err)
		return "", nil, errors.New("internal error reading from world state (2)")
	}
	logger.Printf("Found value=%#v for key=%s\n", res, id)

	return ck, res, nil

}
