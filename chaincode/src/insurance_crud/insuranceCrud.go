package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
)


// InsuranceChaincode
type InsuranceCrudChaincode struct {}

func main() {
	err := shim.Start(new(InsuranceCrudChaincode))
	if err != nil {
		fmt.Printf("Error starting Insurance CRUD chaincode: %s", err)
	}
}

func (t *InsuranceCrudChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke is the entry point to invoke a chaincode function
func (t *InsuranceCrudChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)
	stub.GetCallerCertificate()
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "save" {
		return t.save(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *InsuranceCrudChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "retrieve" {
		return t.retrieve(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//=================================================================================================================================
//	 Saves a value into state, with the specified key
//		args - key, value
//=================================================================================================================================
func (t *InsuranceCrudChaincode) save(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (Key, Value)")
	}

	err := stub.PutState(args[0], []byte(args[1]))

	return nil, err
}

//=================================================================================================================================
//	 Retrieves a value from state, with the specified key
//		args - key
//=================================================================================================================================
func (t *InsuranceCrudChaincode) retrieve(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 (Key)")
	}

	return stub.GetState(args[0])
}
