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

const CRUD_ID_KEY = "CRUD_ID_KEY"

const SAVE_FUNCTION = "save"
const RETRIEVE_FUNCTION = "retrieve"

type Dao struct {}

func (t *Dao) init(stub shim.ChaincodeStubInterface, crudChaincodeId string) {
	stub.PutState(CRUD_ID_KEY, []byte(crudChaincodeId))
}

func SavePolicy(stub shim.ChaincodeStubInterface, policy Policy) (Policy, error) {
	if policy.Id == "" {
		policyId := getNextPolicyId(stub)

		policy.Id = policyId
	}

	bytes, err := marshall(policy)

	if err != nil {fmt.Printf("\nUnable to marshall policy: %s", err); return policy, err}

	err = save(stub, policy.Id, bytes)

	if err != nil {fmt.Printf("Unable to save policy: %s",  err); return policy, err}

	return policy, nil
}

func RetrievePolicy(stub shim.ChaincodeStubInterface, id string) (Policy, error){
	var policy Policy

	bytes, err := retrieve(stub, id)

	if err != nil {	fmt.Printf("RetrievePolicy: Cannot retrieve policy: %s", err); return policy, err}

	err = unmarshal(bytes, &policy);

	if err != nil {	fmt.Printf("RetrievePolicy: Cannot marshall policy: %s", err); return policy, err}

	return policy, nil
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

func retrieve(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	return query(stub, getCrudChaincodeId(stub), RETRIEVE_FUNCTION, []string{id})
}

func save(stub shim.ChaincodeStubInterface, id string, toSave []byte) (error){
	return invoke(stub, getCrudChaincodeId(stub), SAVE_FUNCTION, []string{id, string(toSave)})
}

func invoke(stub shim.ChaincodeStubInterface, chaincodeId, functionName string, args []string) (error){
	_, error := stub.InvokeChaincode(chaincodeId, createArgs(functionName, args))

	return error
}

func query(stub shim.ChaincodeStubInterface, chaincodeId, functionName string, args []string) ([]byte, error){
	return stub.QueryChaincode(chaincodeId, createArgs(functionName, args))
}

func createArgs(functionName string, args[]string) ([][]byte){
	funcAndArgs := append([]string{functionName}, args...)
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

	stub.PutState(idKey, []byte(nextIdNum))

	return idPrefix + nextIdNum
}

func getCurrentIdNumber(stub shim.ChaincodeStubInterface, idKey string) (int) {
	bytes, err := stub.GetState(idKey);

	if err != nil { fmt.Printf("getCurrentIdNumber Error getting id %s", err); return -1}

	currentId, err := strconv.Atoi(string(bytes))

	return currentId;
}
