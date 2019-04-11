/**
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// InvocationHandler is a generic interface for wrapper a
// single chaincode transaction.
type InvocationHandler interface {
	/**
	 * extracts arguments from the stub and checks them
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

func newId(stub shim.ChaincodeStubInterface, indexName string) (string, error) {
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

/**
 * A registry can be used to create and retrieve data items
 * from the world state.
 */
type registryInterface interface {
	key(stub shim.ChaincodeStubInterface, id string) (string, error)
	create(stub shim.ChaincodeStubInterface, item interface{}) (string, error)
	get(stub shim.ChaincodeStubInterface, id string) (string, interface{}, error)
}

/**
 * A concrete registry has a type, given by its name (for creating keys)
 * and its reflect.Type (for creating structs dynamically)
 */
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
	idStr, err := newId(stub, r.typeStr)
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

func shipmentRegistry() registry {
	return registry{
		typeStr: "Shipment",
		typeRT:  reflect.TypeOf(&Shipment{}),
	}
}

func shipmentCoRegistry() registry {
	return registry{
		typeStr: "ShipmentCo",
		typeRT:  reflect.TypeOf(&ShipmentCo{}),
	}
}

func individualParticipantRegistry() registry {
	return registry{
		typeStr: "IndividualParticipant",
		typeRT:  reflect.TypeOf(&IndividualParticipant{}),
	}
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
