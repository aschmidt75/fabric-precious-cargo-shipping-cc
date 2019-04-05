/**
 */
package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	logger = log.New(os.Stdout, "PreciousCargoChaincode: ", log.Ldate|log.Ltime|log.Lmicroseconds)
	ns     = "sample.PreciousCargoChaincode"
)

type InvocationHandler interface {
	/**
	 * extracts arguments from the stub and checks them
	 */
	checkParseArguments(stub shim.ChaincodeStubInterface) error

	/**
	 * runs the chaincode
	 */
	process(stub shim.ChaincodeStubInterface) error

	/**
	 * returns result
	 */
	getResponse(stub shim.ChaincodeStubInterface) interface{}
}

type PreciousCargoChaincode struct {
	// map function names to function implementation types
	handlers map[string]reflect.Type
}

func (cci *PreciousCargoChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Println("enter Init")

	// we could process init arguments here using the stub

	return shim.Success(nil)
}

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
			logger.Printf("invHandler(%x)=%#v\n", inv, inv)
			if err := inv.checkParseArguments(stub); err != nil {
				return shim.Error(err.Error())
			}
			if err := inv.process(stub); err != nil {
				return shim.Error(err.Error())
			}
			r, err := json.Marshal(inv.getResponse(stub))
			if err != nil {
				logger.Println(err)
				return shim.Error("Internal JSON marshal error.")
			}
			return shim.Success([]byte(r))
		}
	}

	return shim.Error("Invalid function name.")
}

func main() {
	logger.Println("Instantiating chaincode.")

	cc := new(PreciousCargoChaincode)
	// add all functions as InvocationHandlers
	cc.handlers = make(map[string]reflect.Type)
	cc.handlers["submitShipment"] = reflect.TypeOf((*submitShipmentInvocation)(nil)).Elem()
	cc.handlers["registerIndividualParticipant"] = reflect.TypeOf((*registerIndividualParticipantInvocation)(nil)).Elem()
	cc.handlers["getIndividualParticipant"] = reflect.TypeOf((*getIndividualParticipantInvocation)(nil)).Elem()

	err := shim.Start(cc)
	if err != nil {
		logger.Fatalf("Error starting chaincode: %s", err)
	}
}
