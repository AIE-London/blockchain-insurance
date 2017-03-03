package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"encoding/json"
	"strconv"
)


// SettlementsChaincode
type SettlementsChaincode struct {}

//Holds a list of keys added to state
type KeyHolder struct {
	Keys 	[]string 	`json:"keys"`
}

const KEYHOLDER_KEY = "KEYHOLDER_KEY"

func main() {
	err := shim.Start(new(SettlementsChaincode))
	if err != nil {
		fmt.Printf("Error starting Insurance Settlements chaincode: %s", err)
	}
}

func (t *SettlementsChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {return nil, errors.New("Incorrect number of init arguments. Expecting 1 (CrudChaincodeId)")}

	InitDao(stub, args[0])
	return nil, nil
}

// Invoke is the entry point to invoke a chaincode function
func (t *SettlementsChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	//Check that this chaincode was invoked internally by another chaincode contract
	/*if !t.isTxIdValid(stub, args) {
		fmt.Println("TxId does not match, aborting")
		return nil, errors.New("TxId does not match")
	}*/
	//Allow anyone to invoke for now

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addPayment" {
		return t.addPayment(stub, args)
	} else if function == "settlePayments" {
		return t.settlePayments(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SettlementsChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	//Check that this chaincode was invoked internally by another chaincode contract
	/*if !t.isTxIdValid(stub, args) {
		fmt.Println("TxId does not match, aborting")
		return nil, errors.New("TxId does not match")
	}*/
	//Allow anyone to query for now

	if function == "retrieveUnsettledPayments" {
		return t.retrieveUnsettledPaymentsJSON(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SettlementsChaincode) addPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("in addPayment()")

	if len(args) != 6 {
		fmt.Println("addPayment: Incorrect number of arguments. Expecting 6 (txId, paymentType, from, to, amount, linkedId)")
		return nil, errors.New("addPayment: Incorrect number of arguments. Expecting 6 (txId, paymentType, from, to, amount, linkedId)")
	}

	if !t.isTxIdValid(stub, args) {
		fmt.Println("TxId does not match, aborting")
		return nil, errors.New("TxId does not match")
	}

	amount, _ := strconv.Atoi(args[4])
	_, err := SavePayment(stub, NewPayment("", args[1], args[2], args[3], amount, args[5]))

	if err != nil {fmt.Printf("\nUnable to save payment: %s", err)}

	return nil, err
}

func (t *SettlementsChaincode) settlePayments(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
 return nil, nil
}

//TODO Inefficient.  Need to come up with a model that makes retrieving unsettled payments faster
func (t *SettlementsChaincode) retrieveUnsettledPaymentsJSON(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	payments := RetrieveAllPayments(stub)

	result := "["

	for _, payment := range payments {
		if payment.Status == PAYMENT_STATUS_PENDING {
			paymentJSON, err := json.Marshal(payments)

			if err != nil {
				fmt.Printf("retrieveUnsettledPaymentsJSON: Cannot marshall policy with id " + payment.Id + ": %s", err)
			} else {
				result += string(paymentJSON) + ","
			}
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//=================================================================================================================================
//	 This chaincode should not be invoked externally.  To confirm that the invoke/query was performed from
//	 another chaincode, the txId must be passed as the first argument, and if this does not match the txId
//   returned from the stub, we should not allow the request to proceed.
//   (Only another chaincode can know the txId before sending the request)
//=================================================================================================================================
func (t *SettlementsChaincode) isTxIdValid(stub shim.ChaincodeStubInterface, args []string) (bool){

	if len(args) == 0 {
		return false
	}

	return stub.GetTxID() == args[0]
}
