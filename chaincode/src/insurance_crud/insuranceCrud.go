package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
)


// InsuranceChaincode
type InsuranceCrudChaincode struct {}

//Holds a list of keys added to state
type KeyHolder struct {
	Keys 	[]string 	`json:"keys"`
}

const KEYHOLDER_KEY = "KEYHOLDER_KEY"

func main() {
	err := shim.Start(new(InsuranceCrudChaincode))
	if err != nil {
		fmt.Printf("Error starting Insurance CRUD chaincode: %s", err)
	}
}

func (t *InsuranceCrudChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//Only applies in dev mode as the chaincode will be deployed with a fresh chaincodeId in prod mode
	t.wipeData(stub)

	//Add new key holder
	var keyHolder KeyHolder
	keyHolder.Keys = []string{}
	bytes, err := json.Marshal(keyHolder)

	if (err != nil) {fmt.Printf("\nError when creating empty key holder: %s", err); return nil, err}

	err = stub.PutState(KEYHOLDER_KEY, bytes)

	return nil, err
}

// Invoke is the entry point to invoke a chaincode function
func (t *InsuranceCrudChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	//Allow init to be run externally
	if function == "init" {
		return t.Init(stub, "init", args)
	}

	//Check that this chaincode was invoked internally by another chaincode contract
	if !t.isTxIdValid(stub, args) {
		fmt.Println("TxId does not match, aborting")
		return nil, errors.New("TxId does not match")
	}

	if function == "save" {
		return t.save(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *InsuranceCrudChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	//Check that this chaincode was invoked internally by another chaincode contract
	if !t.isTxIdValid(stub, args) {
		fmt.Println("TxId does not match, aborting")
		return nil, errors.New("TxId does not match")
	}

	if function == "retrieve" {
		return t.retrieve(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//=================================================================================================================================
//	 Saves a value into state, with the specified key
//		args - txId, key, value
//=================================================================================================================================
func (t *InsuranceCrudChaincode) save(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 (TxId, Key, Value)")
	}

	err := stub.PutState(args[1], []byte(args[2]))

	if (err == nil) {fmt.Println("SAVE (KEY " + args[1] + "): " + string(args[2]))}

	err = t.addKey(stub, args[1])
	return nil, err
}

//=================================================================================================================================
//	 Retrieves a value from state, with the specified key
//		args - txId, key
//=================================================================================================================================
func (t *InsuranceCrudChaincode) retrieve(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (TxId, Key)")
	}

	bytes, err := stub.GetState(args[1])

	if (err == nil) {fmt.Println("RETRIEVE (KEY " + args[1] + "): " + string(bytes))}

	return bytes, err;
}

//=================================================================================================================================
//	 This chaincode should not be invoked externally.  To confirm that the invoke/query was performed from
//	 another chaincode, the txId must be passed as the first argument, and if this does not match the txId
//   returned from the stub, we should not allow the request to proceed.
//   (Only another chaincode can know the txId before sending the request)
//=================================================================================================================================
func (t *InsuranceCrudChaincode) isTxIdValid(stub shim.ChaincodeStubInterface, args []string) (bool){

	if len(args) == 0 {
		return false
	}

	return stub.GetTxID() == args[0]
}

//Wipes all data
func (t *InsuranceCrudChaincode) wipeData(stub shim.ChaincodeStubInterface){

	fmt.Println("WIPING STATE")

	keyHolder, err := t.getKeyHolder(stub)

	if (err != nil) {fmt.Printf("\nError when wiping state: %s", err)}
	if err == nil {
		for _,key := range keyHolder.Keys {
			stub.DelState(key)
			fmt.Println("Deleted key: " + key)
		}
	}
}

func (t *InsuranceCrudChaincode) addKey(stub shim.ChaincodeStubInterface, key string) (error) {
	bytes, _ := stub.GetState(KEYHOLDER_KEY)

	keyHolder, err := t.getKeyHolder(stub)

	if err == nil {
		keyHolder.Keys = append(keyHolder.Keys, key)

		bytes, err = json.Marshal(keyHolder)
		stub.PutState(KEYHOLDER_KEY, bytes)
	}

	return err
}

func (t *InsuranceCrudChaincode) getKeyHolder(stub shim.ChaincodeStubInterface) (KeyHolder, error) {
	bytes, _ := stub.GetState(KEYHOLDER_KEY)

	var keyHolder KeyHolder

	err := json.Unmarshal(bytes, &keyHolder);

	return keyHolder, err
}
