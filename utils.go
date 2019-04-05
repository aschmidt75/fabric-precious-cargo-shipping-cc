/**
 */
package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func newId(stub shim.ChaincodeStubInterface, indexName string) (string, error) {
	ckIndex, err := stub.CreateCompositeKey(ns, []string{".", indexName, ".", "index"})
	if err != nil {
		return "", err
	}

	var lastIndex uint64 = 0
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
