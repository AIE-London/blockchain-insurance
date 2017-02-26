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
const   CURRENT_POLICY_ID_KEY	= "currentPolicyId"
const	CURRENT_USER_ID_KEY	= "currentUserId"
const   CURRENT_CLAIM_ID_KEY	= "currentClaimId"

//==============================================================================================================================
//	 Prefixes for the different domain object type ids
//==============================================================================================================================
const   POLICY_ID_PREFIX	= "P"
const	USER_ID_PREFIX		= "U"
const   CLAIM_ID_PREFIX		= "C"

const CRUD_ID_KEY = "CRUD_ID_KEY"

const SAVE_FUNCTION = "save"
const RETRIEVE_FUNCTION = "retrieve"

type Dao struct {}

func InitDao(stub shim.ChaincodeStubInterface, crudChaincodeId string) {
	stub.PutState(CRUD_ID_KEY, []byte(crudChaincodeId))
	configureIdState(stub)
}

func SavePolicy(stub shim.ChaincodeStubInterface, policy Policy) (Policy, error) {
	if policy.Id == "" {
		policyId := getNextPolicyId(stub)

		policy.Id = policyId
	}

	err := saveObject(stub, policy.Id, policy)

	return policy, err
}

func RetrievePolicy(stub shim.ChaincodeStubInterface, id string) (Policy, error){
	var policy Policy

	err := retrieveObject(stub, id, policy)

	return policy, err
}

func RetrieveAllPolicies(stub shim.ChaincodeStubInterface) ([]Policy){
	var policies []Policy

	numberOfPolicies := getCurrentPolicyIdNumber(stub)

	for i := 1; i <= numberOfPolicies; i++ {

		policyId := POLICY_ID_PREFIX + strconv.Itoa(i)

		policy, err := RetrievePolicy(stub, policyId)

		if err != nil {	fmt.Printf("RetrieveAllPolicies: Error but continuing: %s", err); }

		policies = append(policies, policy)
	}

	return policies
}

func SaveClaim(stub shim.ChaincodeStubInterface, claim Claim) (Claim, error) {
	if claim.Id == "" {
		claimId := getNextClaimId(stub)

		claim.Id = claimId
	}

	err := saveObject(stub, claim.Id, claim)

	return claim, err
}

func RetrieveClaim(stub shim.ChaincodeStubInterface, id string) (Claim, error){
	var claim Claim

	bytes, err := retrieve(stub, id)

	if err != nil {	fmt.Printf("RetrieveClaim: Cannot retrieve claim: %s", err); return claim, err}

	err = unmarshal(bytes, &claim);

	if err != nil {	fmt.Printf("RetrieveClaim: Cannot marshall claim: %s", err); return claim, err}

	return claim, nil
}

func RetrieveAllClaims(stub shim.ChaincodeStubInterface) ([]Claim){
	var claims []Claim

	numberOfClaims := getCurrentClaimIdNumber(stub)

	for i := 1; i <= numberOfClaims; i++ {

		claimId := CLAIM_ID_PREFIX + strconv.Itoa(i)

		claim, err := RetrieveClaim(stub, claimId)

		if err != nil {	fmt.Printf("RetrieveAllClaims: Error but continuing: %s", err); }

		claims = append(claims, claim)
	}

	return claims
}

func SaveVehicle(stub shim.ChaincodeStubInterface, vehicle Vehicle) (Vehicle, error) {
	if vehicle.Id == "" {
		vehicle.Id = vehicle.Details.Registration
	}

	err := saveObject(stub, vehicle.Id, vehicle)

	return vehicle, err
}

func RetrieveVehicle(stub shim.ChaincodeStubInterface, id string) (Vehicle, error){
	var vehicle Vehicle

	err := retrieveObject(stub, id, vehicle)

	return vehicle, err
}

func retrieve(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	return query(stub, getCrudChaincodeId(stub), RETRIEVE_FUNCTION, []string{id})
}

func retrieveObject(stub shim.ChaincodeStubInterface, id string, toStoreObject interface{}) (error){
	bytes, err := retrieve(stub, id)

	if err != nil {	fmt.Printf("RetrievePolicy: Cannot retrieve object with id: " + id + " : %s", err); return err}

	err = unmarshal(bytes, &toStoreObject);

	if err != nil {	fmt.Printf("RetrievePolicy: Cannot unmarshall object with id: " + id + " : %s", err); return err}

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
	bytes, _ := stub.GetState(CRUD_ID_KEY)
	return string(bytes)
}

//==============================================================================================================================
//	 ID Functions - The current id of both policies and claims are stored in blockchain state.
//   This value is incremented when a new policy or claim is created.
//   A prefix is added to the id's to differentiate between policies and claims
//==============================================================================================================================
func getCurrentPolicyIdNumber(stub shim.ChaincodeStubInterface) (int) {
	return getCurrentIdNumber(stub, CURRENT_POLICY_ID_KEY);
}

func getNextPolicyId(stub shim.ChaincodeStubInterface) (string) {

	return getNextId(stub, CURRENT_POLICY_ID_KEY, POLICY_ID_PREFIX);
}

func getNextUserId(stub shim.ChaincodeStubInterface) (string) {

	return getNextId(stub, CURRENT_USER_ID_KEY, USER_ID_PREFIX);
}

func getCurrentClaimIdNumber(stub shim.ChaincodeStubInterface) (int) {
	return getCurrentIdNumber(stub, CURRENT_CLAIM_ID_KEY);
}

func getNextClaimId(stub shim.ChaincodeStubInterface) (string) {

	return getNextId(stub, CURRENT_CLAIM_ID_KEY, CLAIM_ID_PREFIX);
}

func getNextId(stub shim.ChaincodeStubInterface, idKey string, idPrefix string) (string) {

	currentId := getCurrentIdNumber(stub, idKey)

	nextIdNum := strconv.Itoa(currentId + 1)

	save(stub, idKey, []byte(nextIdNum))

	return idPrefix + nextIdNum
}

func getCurrentIdNumber(stub shim.ChaincodeStubInterface, idKey string) (int) {
	bytes, err := retrieve(stub, idKey)

	if err != nil { fmt.Printf("getCurrentIdNumber Error getting id %s", err); return -1}

	currentId, err := strconv.Atoi(string(bytes))

	return currentId;
}

func configureIdState(stub shim.ChaincodeStubInterface) {
	bytes, err := retrieve(stub, CURRENT_POLICY_ID_KEY)
	if (err != nil || len(bytes) == 0) {
		fmt.Println("No existing policy id key, adding all....")
		save(stub, CURRENT_POLICY_ID_KEY, []byte("2"))
		save(stub, CURRENT_CLAIM_ID_KEY, []byte("0"))
	}
}
