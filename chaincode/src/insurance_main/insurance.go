package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
)

// InsuranceChaincode example simple Chaincode implementation
type InsuranceChaincode struct {}

//==============================================================================================================================
//	PoliceReport - Defines the structure for a PoliceReport object.
//==============================================================================================================================
type PoliceReport struct {
	Description         string          `json:"description"`
	Location            Coordinates     `json:"coordinates"`
	DriverAtFault       bool            `json:"driver_at_fault"`
}

//==============================================================================================================================
//	Coordinates - Defines the structure for a Coordinates object.
//==============================================================================================================================
type Coordinates struct {
	x float32 `json:"x"`
	y float32 `json:"y"`
}

//==============================================================================================================================
//	 Keys for obtaining the current id for the different domain object types.
//   Ids are incremental so knowing the latest id is useful when querying for
//   all domain objects of a certain type.
//==============================================================================================================================
const   CURRENT_POLICY_ID_KEY	= "currentPolicyId"
const	CURRENT_VEHICLE_ID_KEY	= "currentVehicleId"
const	CURRENT_USER_ID_KEY	= "currentUserId"
const   CURRENT_CLAIM_ID_KEY	= "currentClaimId"

//Used to store all current approved garages
const	APPROVED_GARAGES_KEY	= "approvedGarages"

//==============================================================================================================================
//	 Prefixes for the different domain object type ids
//==============================================================================================================================
const   POLICY_ID_PREFIX	= "P"
const	VEHICLE_ID_PREFIX	= "V"
const	USER_ID_PREFIX		= "U"
const   CLAIM_ID_PREFIX		= "C"

func main() {
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting Insurance chaincode: %s", err)
	}
}

func (t *InsuranceChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	stub.PutState(CURRENT_POLICY_ID_KEY, []byte("0"))
	stub.PutState(CURRENT_VEHICLE_ID_KEY, []byte("0"))
	stub.PutState(CURRENT_USER_ID_KEY, []byte("0"))
	stub.PutState(CURRENT_CLAIM_ID_KEY, []byte("0"))

	t.initSetup(stub)

	return nil, nil
}

// Invoke is the entry point to invoke a chaincode function
func (t *InsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	caller, caller_affiliation, _ := t.get_caller_data(stub)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addPolicy" {
		return t.addPolicy(stub, caller, caller_affiliation, args)
	} else if function == "createClaim" {
		return t.createClaim(stub, caller, caller_affiliation, args)
	} else if function == "addPoliceReport" {
		//TODO
	} else if function == "addGarageReport" {
		return t.addGarageReport(stub, caller, caller_affiliation, args)
	} else if function == "confirmWork" {
		//TODO
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *InsuranceChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	caller, caller_affiliation, _ := t.get_caller_data(stub)

	// Handle different functions
	if function == "retrieveAllPolicies" {
		return t.retrieveAllPoliciesJSON(stub, caller, caller_affiliation)
	} else if function == "retrieveAllClaims" {
		return t.retrieveAllClaimsJSON(stub, caller, caller_affiliation)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//=================================================================================================================================
//	 Add Policy  - Creates a Policy object and then saves it to the ledger.
//          args - startDate, endDate, excess, vehicle
//=================================================================================================================================
func (t *InsuranceChaincode) addPolicy(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addPolicy()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4 (ActivationDate,ExpiryDate,Excess,VehicleReg)")
	}

	excess, _ := strconv.Atoi(args[2])
	policy := NewPolicy(t.getNextPolicyId(stub), caller, args[0], args[1], excess, args[3])

	bytes, err := json.Marshal(policy)

	if err != nil { fmt.Printf("addPolicy Error converting policy record: %s", err); return nil, errors.New("Error converting policy record") }

	err = stub.PutState(policy.Id, bytes)

	if err != nil { fmt.Printf("addPolicy: Error storing policy record: %s", err); return nil, errors.New("Error storing policy record") }

	return nil, nil
}

//=================================================================================================================================
//	 Add Vehicle  - Creates a Vehicle object and then saves it to the ledger.
//          args - make, model, registration, year, mileage
//=================================================================================================================================
func (t *InsuranceChaincode) addVehicle(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addVehicle()")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 (Make, Model, Registration, Year, Mileage)")
	}

	mileage, _ := strconv.Atoi(args[4])
	vehicle := NewVehicle(t.getNextVehicleId(stub), args[0], args[1], args[2], args[3], mileage)

	bytes, err := json.Marshal(vehicle)

	if err != nil { fmt.Printf("addVehicle Error converting vehicle record: %s", err); return nil, errors.New("Error converting vehicle record") }

	err = stub.PutState(vehicle.Id, bytes)

	if err != nil { fmt.Printf("addVehicle: Error storing vehicle record: %s", err); return nil, errors.New("Error storing vehicle record") }

	return nil, nil
}

//=================================================================================================================================
//	 Add User  - Creates a User object and then saves it to the ledger.
//          args - forename, surname, email, relatedPolicy
//=================================================================================================================================
func (t *InsuranceChaincode) addUser(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addUser()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4 (Forename, Surname, Email, RelatedPolicy)")
	}

	user := NewUser(t.getNextUserId(stub), args[0], args[1], args[2], args[3])

	bytes, err := json.Marshal(user)

	if err != nil { fmt.Printf("addUser Error converting user record: %s", err); return nil, errors.New("Error converting user record") }

	err = stub.PutState(user.Id, bytes)

	if err != nil { fmt.Printf("addUser: Error storing user record: %s", err); return nil, errors.New("Error storing user record") }

	return nil, nil
}

//=================================================================================================================================
//	 appendApprovedGarages		-	Appends approved garages to the list of approved garages in the ledger
//=================================================================================================================================
func (t *InsuranceChaincode) appendApprovedGarages(stub shim.ChaincodeStubInterface, approvedGarages []string) ([]byte, error){

	garages, err := stub.GetState(APPROVED_GARAGES_KEY)

	if err != nil {
		return nil, errors.New("Failed to get allApprovedGarages")
	}

	var newGarages = ApprovedGarages{}

	var approvedGaragesList ApprovedGarages
	json.Unmarshal(garages, &approvedGaragesList)

	newGarages.Garages = approvedGarages

	approvedGaragesList.Garages = append(approvedGaragesList.Garages, newGarages.Garages...)

	bytes, err := json.Marshal(approvedGaragesList)

	if err != nil {
		fmt.Printf("addApprovedGarages: Error converting garages: %s", err); return nil, errors.New("Error converting garages")
	}

	err = stub.PutState(APPROVED_GARAGES_KEY, bytes)

	if err != nil {
		fmt.Printf("addApprovedGarages: Error adding approved garages: %s", err); return nil, errors.New("Error adding approved garages")
	}

	return nil, nil
}

//=================================================================================================================================
//	 createClaim - Creates a Claim object and then saves it to the ledger.
//          args - RelatedPolicy,Description,Date,IncidentType
//=================================================================================================================================
func (t *InsuranceChaincode) createClaim(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running createClaim()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4 (RelatedPolicy,Description, Date, IncidentType)")
	}

	claim := NewClaim(t.getNextClaimId(stub), args[0], args[1], args[2], args[3])

	isValid, err := t.checkClaimIsValid(stub, caller, claim);

	if !isValid {
		fmt.Printf("createClaim: Claim invalid: %s", err);
		return nil, errors.New("Claim is invalid");
	}

	bytes, err := json.Marshal(claim)

	if err != nil { fmt.Printf("createClaim Error converting claim record: %s", err); return nil, errors.New("Error converting claim record") }

	err = stub.PutState(claim.Id, bytes)

	if err != nil { fmt.Printf("createClaim: Error storing claim record: %s", err); return nil, errors.New("Error storing claim record") }

	return nil, nil
}

//=================================================================================================================================
//	 checkClaimIsValid - Performs checks on the claim to ensure it is valid
//      Returns true if valid, or false and an error if invalid
//=================================================================================================================================
func (t *InsuranceChaincode) checkClaimIsValid(stub shim.ChaincodeStubInterface, caller string, claim Claim) (bool, error) {

	//Check policy exists
	policy, err := t.retrievePolicy(stub, claim.Relations.RelatedPolicy);

	if err != nil {
		fmt.Printf("checkClaimIsValid: Error getting policy with id %s", claim.Relations.RelatedPolicy);
		return false, errors.New("Policy doesnt exist");
	}

	//Check policy owner matches current user
	if policy.Relations.Owner != caller {
		fmt.Printf("checkClaimIsValid: Policy owner is incorrect %s", claim.Relations.RelatedPolicy);
		return false, errors.New("Policy owner incorrect");
	}

	return true, nil;
}

//==============================================================================================================================
//	 retrieveAllPoliciesJSON - Iterates through all policy ids, retreiving each and returning a JSON representation
//==============================================================================================================================
func (t *InsuranceChaincode) retrieveAllPoliciesJSON(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	numberOfPolicies := t.getCurrentPolicyIdNumber(stub)

	result := "["

	for i := 1; i <= numberOfPolicies; i++ {

		policyId := POLICY_ID_PREFIX + strconv.Itoa(i)
		policyJSON, err := t.retrievePolicyJSON(stub, policyId)

		//TODO Check caller has rights to see this policy
		if err != nil {
			fmt.Printf("retrievePolicies: Cannot retrieve policy with id " + policyId+": %s", err)
		} else {
			result += string(policyJSON) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//==============================================================================================================================
//	 retrievePolicy - Gets the state of the data at policyId in the ledger then converts it from the stored
//					JSON into the Policy struct for use in the contract. Returns the Policy struct.
//					Returns empty policy if it errors.
//==============================================================================================================================
func (t *InsuranceChaincode) retrievePolicy(stub shim.ChaincodeStubInterface, policyId string) (Policy, error) {

	var policy Policy

	bytes, err := t.retrievePolicyJSON(stub, policyId)

	if err != nil {	fmt.Printf("retrievePolicy: Cannot read policy: %s", err); return policy, errors.New("retrievePolicy: Cannot read policy")}

	err = json.Unmarshal(bytes, &policy);

	if err != nil {	fmt.Printf("retrievePolicy: Corrupt policy record "+string(bytes)+": %s", err); return policy, errors.New("retrievePolicy: Corrupt policy record"+string(bytes))}

	return policy, nil
}

//==============================================================================================================================
//	 retrievePolicyJSON - Gets the state of the data at policyId in the ledger and returns the JSON representation
//==============================================================================================================================
func (t *InsuranceChaincode) retrievePolicyJSON(stub shim.ChaincodeStubInterface, policyId string) ([]byte, error) {

	bytes, err := stub.GetState(policyId);

	if err != nil {	fmt.Printf("retrievePolicy: Failed to invoke: %s", err); return nil, errors.New("retrievePolicy: Error retrieving policy with policyId = " + policyId) }

	return bytes, nil
}

//==============================================================================================================================
//	 retrieveAllClaimsJSON - Iterates through all claim ids, retreiving each and returning a JSON representation
//==============================================================================================================================
func (t *InsuranceChaincode) retrieveAllClaimsJSON(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	numberOfClaims := t.getCurrentClaimIdNumber(stub)

	result := "["

	for i := 1; i <= numberOfClaims; i++ {

		claimId := CLAIM_ID_PREFIX + strconv.Itoa(i)
		claimJSON, err := t.retrieveClaimJSON(stub, claimId)

		//TODO Check caller has rights to see this claim
		if err != nil {
			fmt.Printf("retrieveAllClaims: Cannot retrieve claim with id " + claimId+": %s", err)
		} else {
			result += string(claimJSON) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//==============================================================================================================================
//	 retrieveClaim - Gets the state of the data at claimId in the ledger then converts it from the stored
//					JSON into the Claim struct for use in the contract. Returns the Claim struct.
//					Returns empty claim if it errors.
//==============================================================================================================================
func (t *InsuranceChaincode) retrieveClaim(stub shim.ChaincodeStubInterface, claimId string) (Claim, error) {

	var claim Claim

	bytes, err := t.retrieveClaimJSON(stub, claimId)

	if err != nil {	fmt.Printf("retrieveClaim: Cannot read claim: %s", err); return claim, errors.New("retrieveClaim: Cannot read claim")}

	err = json.Unmarshal(bytes, &claim);

	if err != nil {	fmt.Printf("retrieveClaim: Corrupt claim record "+string(bytes)+": %s", err); return claim, errors.New("retrieveClaim: Corrupt claim record"+string(bytes))}

	return claim, nil
}

//==============================================================================================================================
//	 retrieveClaimJSON - Gets the state of the data at claimId in the ledger and returns the JSON representation
//==============================================================================================================================
func (t *InsuranceChaincode) retrieveClaimJSON(stub shim.ChaincodeStubInterface, claimId string) ([]byte, error) {

	bytes, err := stub.GetState(claimId);

	if err != nil {	fmt.Printf("retrieveClaimJSON: Failed to invoke: %s", err); return nil, errors.New("retrieveClaimJSON: Error retrieving claim with claimId = " + claimId) }

	return bytes, nil
}

//==============================================================================================================================
//	 ID Functions - The current id of both policies and claims are stored in blockchain state.
//   This value is incremented when a new policy or claim is created.
//   A prefix is added to the id's to differentiate between policies and claims
//==============================================================================================================================
func (t *InsuranceChaincode) getCurrentPolicyIdNumber(stub shim.ChaincodeStubInterface) (int) {
	return t.getCurrentIdNumber(stub, CURRENT_POLICY_ID_KEY);
}

func (t *InsuranceChaincode) getNextPolicyId(stub shim.ChaincodeStubInterface) (string) {

	return t.getNextId(stub, CURRENT_POLICY_ID_KEY, POLICY_ID_PREFIX);
}

func (t *InsuranceChaincode) getNextVehicleId(stub shim.ChaincodeStubInterface) (string) {

	return t.getNextId(stub, CURRENT_VEHICLE_ID_KEY, VEHICLE_ID_PREFIX);
}

func (t *InsuranceChaincode) getNextUserId(stub shim.ChaincodeStubInterface) (string) {

	return t.getNextId(stub, CURRENT_USER_ID_KEY, USER_ID_PREFIX);
}

func (t *InsuranceChaincode) getCurrentClaimIdNumber(stub shim.ChaincodeStubInterface) (int) {
	return t.getCurrentIdNumber(stub, CURRENT_CLAIM_ID_KEY);
}

func (t *InsuranceChaincode) getNextClaimId(stub shim.ChaincodeStubInterface) (string) {

	return t.getNextId(stub, CURRENT_CLAIM_ID_KEY, CLAIM_ID_PREFIX);
}

func (t *InsuranceChaincode) getNextId(stub shim.ChaincodeStubInterface, idKey string, idPrefix string) (string) {

	currentId := t.getCurrentIdNumber(stub, idKey)

	nextIdNum := strconv.Itoa(currentId + 1)

	stub.PutState(idKey, []byte(nextIdNum))

	return idPrefix + nextIdNum
}

func (t *InsuranceChaincode) getCurrentIdNumber(stub shim.ChaincodeStubInterface, idKey string) (int) {
	bytes, err := stub.GetState(idKey);

	if err != nil { fmt.Printf("getCurrentIdNumber Error getting id %s", err); return -1}

	currentId, err := strconv.Atoi(string(bytes))

	return currentId;
}

//==============================================================================================================================
//	 Init Data Setup
//==============================================================================================================================
// Sets up all the data on deployment of chaincode
//==============================================================================================================================
func (t *InsuranceChaincode)initSetup(stub shim.ChaincodeStubInterface) {
	//Policy init data
	policyId1 := t.getNextPolicyId(stub)
	t.addPolicy(stub, PolicyCaller1, PolicyCallerAffiliation1, PolicyArgs1)

	policyId2 := t.getNextPolicyId(stub)
	t.addPolicy(stub, PolicyCaller2, PolicyCallerAffiliation2, PolicyArgs2)

	//Vehicle init data
	t.addVehicle(stub, VehicleCaller1, VehicleCallerAffiliation1, VehicleArgs1)

	t.addVehicle(stub, VehicleCaller2, VehicleCallerAffiliation2, VehicleArgs2)

	//User init data
	completeUserArgs1 := append(UserArgs1, policyId1)
	t.addUser(stub, UserCaller1, UserCallerAffiliation1, completeUserArgs1)

	completeUserArgs2 := append(UserArgs2, policyId2)
	t.addUser(stub, UserCaller2, UserCallerAffiliation2, completeUserArgs2)

	//ApprovedGarages init data
	t.appendApprovedGarages(stub, ApprovedGaragesData)
}

//==============================================================================================================================
//	 Security Functions
//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *InsuranceChaincode)  get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error){

	user, err := t.get_username(stub)

	// if err != nil { return "", "", err }

	// ecert, err := t.get_ecert(stub, user);

	// if err != nil { return "", "", err }

	affiliation, err := t.check_affiliation(stub);

	if err != nil { return "", "", err }

	return user, affiliation, nil
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *InsuranceChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	//username, err := stub.ReadCertAttribute("username");
	//if err != nil {
	//	fmt.Printf("Couldn't get attribute 'username'. Error: %s", err)
	//	return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	//}
    //
	//return string(username), nil
	return "dummy-user", nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *InsuranceChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
	//affiliation, err := stub.ReadCertAttribute("role");
	//if err != nil { return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error()) }
	//return string(affiliation), nil
	return "dummy-affiliation", nil
}

//==============================================================================================================================
//	 addGarageReport - This method adds the garage report's details into the claim
//   args{claimId, garage, estimated_cost,  writeOff, Note}
//==============================================================================================================================
func (t *InsuranceChaincode) addGarageReport(stub shim.ChaincodeStubInterface,  caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addGarageReport()")

	if len(args) != 5 {
		fmt.Println("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 6 (claimId, garage, estimated_cost, writeOff, Note)")
		return nil, errors.New("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 6 (claimId, garage, estimated_cost, writeOff, Note)")
	}
		
	// Does a claim exist for this vehicle?
	var theClaim Claim
	var claimId string = args[0]
	
	theClaim, err := t.retrieveClaim(stub , claimId)

	if err != nil {	fmt.Printf("\nADD_GARAGE_REPORT: Failed to retrieve claim Id: %s", err); return nil, errors.New("ADD_GARAGE_REPORT: Error retrieving claim with claimId = " + claimId) }
		
	// Is the Claim in a valid state to receive a garage report?
    if STATE_AWAITING_GARAGE_REPORT != theClaim.Details.Status {
		fmt.Printf("ADD_GARAGE_REPORT: unexpected Input for claim with Id: %s\n", claimId)
		return nil, errors.New("ADD_GARAGE_REPORT: unexpected Input for claim with claimId = " + claimId)
	}
	
	var report ClaimDetailsClaimGarageReport
	report, err = NewGarageReport(args[1], args[2], args[3], args[4])

	theClaim.Details.Repair = report

	theClaim.Details.Status = STATE_PENDING_AFTER_REPORT_DECISION

	t.saveClaim(stub, theClaim)	

	t.afterReportProcess(stub, caller, caller_affiliation, args)
	
	return nil, nil
}

//=========================================================================================
// This Function process the claim after a report has arrived
//=========================================================================================
func (t *InsuranceChaincode) afterReportProcess(stub shim.ChaincodeStubInterface,  caller string, caller_affiliation string, args []string) ([]byte, error) {
	fmt.Println("running afterReportProcess()")
    var theClaim Claim
	var claimId string = args[0]	
	
	theClaim, err := t.retrieveClaim(stub , claimId)
	
	if err != nil {	fmt.Printf("\nAFTER_REPORT_PROCESS: Failed to retrieve claim Id: %s", err); return nil, errors.New("AFTER_REPORT_PROCESS: Error retrieving claim with claimId = " + claimId) }
	
	if theClaim.Details.Status == STATE_PENDING_AFTER_REPORT_DECISION{

		if theClaim.Details.Repair.WriteOff || theClaim.Details.Repair.Estimate > (t.getCarValue() * 50 / 100) {
			//process total_loss
			t.processTotalLoss(stub, theClaim)
		}else{
			theClaim.Details.Status = STATE_ORDER_GARAGE_WORK
			t.saveClaim(stub, theClaim)
		}
	}
		
	return nil, nil
}

//=========================================================================================
//
//=========================================================================================
func (t *InsuranceChaincode) processTotalLoss(stub shim.ChaincodeStubInterface,  theClaim Claim) ([]byte, error) {

	fmt.Println("running processTotalLoss()")
    if &theClaim != nil{
		var settlement ClaimDetailsSettlement
		settlement.Decision = TOTAL_LOSS
		settlement.Dispute = false
		theClaim.Details.Settlement = settlement
		theClaim.Details.Status = STATE_TOTAL_LOSS_ESTABLISHED
		t.saveClaim(stub, theClaim)
	}else {fmt.Printf("PROCESS_TOTAL_LOSS: Error can not process Total loss on unexisting claim\n"); return nil, errors.New("PROCESS_TOTAL_LOSS: Error can not process Total loss on unexisting claim") }
	return nil, nil
}

//===============================================================================
// This method saves the claim after updates
//===============================================================================
func (t *InsuranceChaincode) saveClaim(stub shim.ChaincodeStubInterface, theClaim Claim) (bool, error) {

	bytes, err := json.Marshal(theClaim)
	if err != nil { fmt.Printf("\nSAVE_CLAIM Error converting Claim details: %s", err); return false, errors.New("Error converting Claim details") }
	
	fmt.Printf("About to update claim %s \n", theClaim.Id)
	err = stub.PutState(theClaim.Id, bytes)
	if err != nil { fmt.Printf("\nSAVE_CLAIM: Error adding garage Report: %s", err); return false, errors.New("Error storing claim details") }
	
	return true, nil
}

//===============================================================================
//  Temporary function to return a hardcoded value of the car after the fix
//  Basicaly, this will call he insurer core API to get the Car value
//===============================================================================
func (t *InsuranceChaincode) getCarValue() (int) {

	return 3000;
}

