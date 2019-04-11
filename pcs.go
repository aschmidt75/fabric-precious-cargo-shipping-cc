/**
 */
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

	// iterate through map of known functions
	for key, invType := range cci.handlers {
		if function == key {
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
