// fabric-precious-cargo-shipping-cc is a sample chaincode for Hyperledger Fabric
// Copyright (C) 2019 @aschmidt75
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.package main
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	ns     = "sample.PreciousCargoChaincode"
	logger = log.New(os.Stdout, fmt.Sprintf("%s: ", ns), log.Ldate|log.Ltime|log.Lmicroseconds)
)

// PreciousCargoChaincode is the Chaincode wrapper for PreciousCargoShipment
type PreciousCargoChaincode struct {
	// map function names to function implementation types
	handlers map[string]reflect.Type
}

// Init initializes chaincode
func (cci *PreciousCargoChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Println("enter Init")

	// we could process init arguments here using the stub

	return shim.Success(nil)
}

// Invoke a chaincode function according to function namen and handlers.
func (cci *PreciousCargoChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Println("enter Invoke")

	function, args := stub.GetFunctionAndParameters()
	logger.Printf("requested function=%s, with args=%#v", function, args)

	if invType, found := cci.handlers[function]; found {
		// from invType as reflect.Type, create a new object and
		// cast its interface to InvocationHandler.
		inv := reflect.New(invType).Interface().(InvocationHandler)
		// let it check its input
		if err := inv.checkParseArguments(stub); err != nil {
			return shim.Error(err.Error())
		}
		// run the transaction
		if err := inv.process(stub); err != nil {
			return shim.Error(err.Error())
		}
		// send out the response
		r, err := json.Marshal(inv.getResponse(stub))
		if err != nil {
			logger.Println(err)
			return shim.Error("Internal JSON marshal error (response).")
		}
		return shim.Success([]byte(r))
	}

	return shim.Error("Invalid function name.")
}

func main() {
	logger.Println("Instantiating chaincode.")

	cc := &PreciousCargoChaincode{
		// all functions as InvocationHandlers
		handlers: map[string]reflect.Type{
			"submitShipment":                reflect.TypeOf((*submitShipmentInvocation)(nil)).Elem(),
			"getShipment":                   reflect.TypeOf((*getShipmentInvocation)(nil)).Elem(),
			"registerIndividualParticipant": reflect.TypeOf((*registerIndividualParticipantInvocation)(nil)).Elem(),
			"getIndividualParticipant":      reflect.TypeOf((*getIndividualParticipantInvocation)(nil)).Elem(),
			"registerShipmentCo":            reflect.TypeOf((*registerShipmentCoInvocation)(nil)).Elem(),
		},
	}
	err := shim.Start(cc)
	if err != nil {
		logger.Fatalf("Error starting chaincode: %s", err)
	}
}
