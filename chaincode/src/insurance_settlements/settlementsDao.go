package main

import (
	"github.com/hyperledger/fabric/core/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"strconv"
	"fmt"
)

//==============================================================================================================================
//	 Keys for obtaining the current id for the different domain object types.
//   Ids are incremental so knowing the latest id is useful when querying for
//   all domain objects of a certain type.
//==============================================================================================================================
const   CURRENT_PAYMENT_ID_KEY	= "currentPaymentId"
const	CURRENT_INVOICE_ID_KEY	= "currentInvoiceId"

//Stored the chaincodeId of the CRUD chaincode
const CRUD_CHAINCODE_ID_KEY = "CRUD_CHAINCODE_ID"

//==============================================================================================================================
//	 Prefixes for the different domain object type ids
//==============================================================================================================================
const   PAYMENT_ID_PREFIX	= "PAY"
const	INVOICE_ID_PREFIX	= "I"

const SAVE_FUNCTION = "save"
const RETRIEVE_FUNCTION = "retrieve"

type Dao struct {}

func InitDao(stub shim.ChaincodeStubInterface, crudChaincodeId string) {
	stub.PutState(CRUD_CHAINCODE_ID_KEY, []byte(crudChaincodeId))
	configureIdState(stub)
}

func SavePayment(stub shim.ChaincodeStubInterface, payment Payment) (Payment, error) {
	if payment.Id == "" {
		paymentId := getNextPaymentId(stub)

		payment.Id = paymentId
	}

	err := saveObject(stub, payment.Id, payment)

	return payment, err
}

func RetrievePayment(stub shim.ChaincodeStubInterface, id string) (Payment, error){
	var payment Payment

	err := retrieveObject(stub, id, &payment)

	return payment, err
}

func RetrieveAllPayments(stub shim.ChaincodeStubInterface) ([]Payment){
	var payments []Payment

	numberOfPayments := getCurrentPaymentIdNumber(stub)

	for i := 1; i <= numberOfPayments; i++ {

		paymentId := PAYMENT_ID_PREFIX + strconv.Itoa(i)

		payment, err := RetrievePayment(stub, paymentId)

		if err != nil {	fmt.Printf("RetrieveAllPayments: Error but continuing: %s", err); }

		payments = append(payments, payment)
	}

	return payments
}

func retrieve(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	return query(stub, getCrudChaincodeId(stub), RETRIEVE_FUNCTION, []string{id})
}

func retrieveObject(stub shim.ChaincodeStubInterface, id string, toStoreObject interface{}) (error){
	bytes, err := retrieve(stub, id)

	if err != nil {	fmt.Printf("RetrieveObject: Cannot retrieve object with id: " + id + " : %s", err); return err}

	err = unmarshal(bytes, toStoreObject)

	if err != nil {	fmt.Printf("RetrieveObject: Cannot unmarshall object with id: " + id + " : %s", err); return err}

	return nil
}

func saveObject(stub shim.ChaincodeStubInterface, id string, object interface{}) (error){
	bytes, err := marshall(object)

	if err != nil {fmt.Printf("\nUnable to marshall object with id: " + id + " : %s", err); return err}

	err = save(stub, id, bytes)

	if err != nil {fmt.Printf("Unable to save policy: %s",  err); return err}

	return nil
}

func save(stub shim.ChaincodeStubInterface, id string, toSave []byte) (error){
	return invoke(stub, getCrudChaincodeId(stub), SAVE_FUNCTION, []string{id, string(toSave)})
}

func invoke(stub shim.ChaincodeStubInterface, chaincodeId, functionName string, args []string) (error){
	_, error := stub.InvokeChaincode(chaincodeId, createArgs(functionName, stub.GetTxID(), args))

	return error
}

func query(stub shim.ChaincodeStubInterface, chaincodeId, functionName string, args []string) ([]byte, error){
	return stub.QueryChaincode(chaincodeId, createArgs(functionName, stub.GetTxID(), args))
}

func createArgs(functionName string, txId string, args[]string) ([][]byte){
	funcAndArgs := append([]string{functionName}, txId)
	funcAndArgs = append(funcAndArgs, args...)
	return util.ArrayToChaincodeArgs(funcAndArgs)
}

func marshall(toMarshall interface{}) ([]byte, error) {
	return json.Marshal(toMarshall)
}

func unmarshal(data []byte, object interface{}) (error){
	return json.Unmarshal(data, object)
}

func getCrudChaincodeId(stub shim.ChaincodeStubInterface) (string) {
	bytes, _ := stub.GetState(CRUD_CHAINCODE_ID_KEY)
	return string(bytes)
}

//==============================================================================================================================
//	 ID Functions - The current id of both policies and claims are stored in blockchain state.
//   This value is incremented when a new policy or claim is created.
//   A prefix is added to the id's to differentiate between policies and claims
//==============================================================================================================================
func getCurrentPaymentIdNumber(stub shim.ChaincodeStubInterface) (int) {
	return getCurrentIdNumber(stub, CURRENT_PAYMENT_ID_KEY);
}

func getNextPaymentId(stub shim.ChaincodeStubInterface) (string) {

	return getNextId(stub, CURRENT_PAYMENT_ID_KEY, PAYMENT_ID_PREFIX);
}

func getNextInvoiceId(stub shim.ChaincodeStubInterface) (string) {

	return getNextId(stub, CURRENT_INVOICE_ID_KEY, INVOICE_ID_PREFIX);
}

func getNextId(stub shim.ChaincodeStubInterface, idKey string, idPrefix string) (string) {

	currentId := getCurrentIdNumber(stub, idKey)

	nextIdNum := strconv.Itoa(currentId + 1)

	save(stub, idKey, []byte(nextIdNum))
	//Comment below regarding why we need to store this locally
	stub.PutState(idKey, []byte(nextIdNum))

	return idPrefix + nextIdNum
}

func getCurrentIdNumber(stub shim.ChaincodeStubInterface, idKey string) (int) {
	bytesFromCrud, err := retrieve(stub, idKey)

	if err != nil { fmt.Printf("getCurrentIdNumber Error getting id %s", err); return -1}

	currentIdFromCrud, _ := strconv.Atoi(string(bytesFromCrud))

	//There seems to be a bug in hyperledger where the state within a transaction is not correct within a different chaincode
	//when invoking that chaincode multiple times within the same transaction, so we need to keep track of the id state locally
	bytesFromLocal, _ := stub.GetState(idKey)
	currentIdLocal, _ := strconv.Atoi(string(bytesFromLocal))

	//Check which is higher and use that
	if (currentIdLocal > currentIdFromCrud) { return currentIdLocal }

	return currentIdFromCrud
}

func configureIdState(stub shim.ChaincodeStubInterface) {
	bytes, err := retrieve(stub, CURRENT_PAYMENT_ID_KEY)
	if (err != nil || len(bytes) == 0) {
		fmt.Println("No existing policy id key, adding all....")
		save(stub, CURRENT_PAYMENT_ID_KEY, []byte("0"))
		save(stub, CURRENT_INVOICE_ID_KEY, []byte("0"))
	}

	//There seems to be a bug in hyperledger where the state within a transaction is not correct within a different chaincode
	//when invoking that chaincode multiple times within the same transaction, so we need to keep track of the id state locally
	stub.PutState(CURRENT_PAYMENT_ID_KEY, []byte("0"))
	stub.PutState(CURRENT_INVOICE_ID_KEY, []byte("0"))
}
