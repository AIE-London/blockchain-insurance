package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
)

// InsuranceChaincode
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
//	 Roles within the system
//==============================================================================================================================
const   ROLE_POLICY_HOLDER	= "policyholder"
const	ROLE_GARAGE			= "garage"
const   ROLE_INSURER		= "insurer"
const   ROLE_SUPER_USER		= "superuser"
const   ROLE_ORACLE			= "oracle"

func main() {
	err := shim.Start(new(InsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting Insurance chaincode: %s", err)
	}
}

func (t *InsuranceChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	InitDao(stub, args[0])
	InitOracleService(stub, args)
	InitReferenceData(stub)

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
	} else if function == "declareLiability" {
		return t.declareLiability(stub, caller, caller_affiliation, args)
	} else if function == "addPoliceReport" {
		//TODO
	} else if function == "addGarageReport" {
		return t.addGarageReport(stub, caller, caller_affiliation, args)
	} else if function == "confirmWork" {
		//TODO
	} else if function == "agreePayoutAmount" {
		return t.agreePayoutAmount(stub, caller, caller_affiliation, args)
	} else if function == "vehicleValueOracleCallback" {
		fmt.Printf("vehicleValueOracleCallback: ReqeustId / Value" + args[0] + " / " + args[1]);
		return t.vehicleValueOracleCallback(stub, caller, caller_affiliation, args)
	} else if function == "confirmPaidOut" {
		return t.confirmPaidOut(stub, caller, caller_affiliation, args)
	} else if function == "closeClaim" {
		return t.closeClaim(stub, caller, caller_affiliation, args)
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
//          args - owner, startDate, endDate, excess, vehicle
//=================================================================================================================================
func (t *InsuranceChaincode) addPolicy(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addPolicy()")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 (Owner,ActivationDate,ExpiryDate,Excess,VehicleReg)")
	}

	if caller_affiliation != ROLE_INSURER {
		return nil, errors.New("Only an insurer can add a new policy")
	}

	excess, _ := strconv.Atoi(args[3])
	policy := NewPolicy("", args[0], caller, args[1], args[2], excess, args[4])

	_, err := SavePolicy(stub, policy)

	return nil, err
}

//=================================================================================================================================
//	 Add Vehicle  - Creates a Vehicle object and then saves it to the ledger.
//          args - make, model, registration, year, mileage, styleId
//=================================================================================================================================
func (t *InsuranceChaincode) addVehicle(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addVehicle()")

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6 (Make, Model, Registration, Year, Mileage, StyleId)")
	}

	mileage, _ := strconv.Atoi(args[4])
	vehicle := NewVehicle(args[2], args[0], args[1], args[3], mileage, args[5])

	_, err := SaveVehicle(stub, vehicle)

	return nil, err
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

	user := NewUser("", args[0], args[1], args[2], args[3])

	_, err := SaveUser(stub, user)

	return nil, err
}

//=================================================================================================================================
//	 createClaim - Creates a Claim object and then saves it to the ledger.
//          args - RelatedPolicy,Description,Date,IncidentType,OtherPartyReg(Optional),liable(Optional)
//=================================================================================================================================
func (t *InsuranceChaincode) createClaim(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running createClaim()")

	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting at least (RelatedPolicy,Description, Date, IncidentType)")
	}

	claim := NewClaim("", args[0], args[1], args[2], args[3])

	isValid, err := t.checkClaimIsValid(stub, caller, claim);

	if !isValid {
		fmt.Printf("createClaim: Claim invalid: %s", err);
		return nil, errors.New("Claim is invalid");
	}

	if (claim.Details.Incident.Type == SINGLE_PARTY) {
		claim.Details.Status = STATE_AWAITING_GARAGE_REPORT
		_, err = SaveClaim(stub, claim)

		return nil, err
	} else if (claim.Details.Incident.Type == MULTIPLE_PARTIES) {
		isLiable, err := strconv.ParseBool(args[5])
		if err != nil { return nil, err }
		return t.processMultiplePartyClaimCreation(stub, claim, args[4], isLiable)
	} else {
		return nil, errors.New("Unsupported claim type: " + claim.Details.Incident.Type)
	}
}

//=================================================================================================================================
//	 processMultiplePartyClaimCreation - Performs additional processing required when creating a multiple party claim
//=================================================================================================================================
func (t *InsuranceChaincode) processMultiplePartyClaimCreation(stub shim.ChaincodeStubInterface, claim Claim, otherPartyReg string, liable bool) ([]byte, error) {
	//Set reg and liable
	claim.Relations.OtherPartyReg = otherPartyReg
	claim.Details.IsLiable = liable

	var status = STATE_AWAITING_LIABILITY_ACCEPTANCE

	//If raising party accepts they're liable initially, then skip the acceptance step
	if liable {
		status = STATE_AWAITING_GARAGE_REPORT
	}

	claim.Details.Status = status;

	//Save claim so that id is generated
	claim, err := SaveClaim(stub, claim)
	if err != nil { return nil, err}

	policy, _ := RetrievePolicy(stub, claim.Relations.RelatedPolicy)
	//Get the other party policy
	otherPartyPolicy, err := t.findPolicyWithVehicleReg(stub, otherPartyReg)

	if err != nil {
		//TODO: what should we do here???
		return nil, err
	}

	//Create claim for other party
	otherClaim := NewClaim("", otherPartyPolicy.Id, claim.Details.Description, claim.Details.Incident.Date, MULTIPLE_PARTIES)
	otherClaim.Details.Status = status
	otherClaim.Relations.OtherPartyReg = policy.Relations.Vehicle

	//Start as the opposite of the initial claim, subject to dispute
	otherClaim.Details.IsLiable = !liable

	//Link with initial claim
	otherClaim.Relations.LinkedClaims = append(otherClaim.Relations.LinkedClaims, claim.Id)
	savedClaim, err := SaveClaim(stub, otherClaim)
	if err != nil { return nil, err}

	//Finally link saved 'other party' claim with the original claim and save again
	claim.Relations.LinkedClaims = append(claim.Relations.LinkedClaims, savedClaim.Id)

	_, err = SaveClaim(stub, claim)

	return nil, err
}

//=================================================================================================================================
//	 findPolicyWithVehicleReg - Finds a policy based on the vehicle registration
//=================================================================================================================================
func (t *InsuranceChaincode) findPolicyWithVehicleReg(stub shim.ChaincodeStubInterface, vehicleReg string) (Policy, error) {
	policies := RetrieveAllPolicies(stub)

	for _, policy := range policies {
		if policy.Relations.Vehicle == vehicleReg { return policy, nil }
	}

	return Policy{}, errors.New("Policy does not exist for registration: " + vehicleReg);
}

//=================================================================================================================================
//	 checkClaimIsValid - Performs checks on the claim to ensure it is valid
//      Returns true if valid, or false and an error if invalid
//=================================================================================================================================
func (t *InsuranceChaincode) checkClaimIsValid(stub shim.ChaincodeStubInterface, caller string, claim Claim) (bool, error) {

	//Check policy exists
	policy, err := RetrievePolicy(stub, claim.Relations.RelatedPolicy);

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
	policies := RetrieveAllPolicies(stub)

	result := "["

	for _, policy := range policies {
		if t.isPolicyRelevantToCaller(stub, policy, caller, caller_affiliation) {
			policyJSON, err := json.Marshal(policy)

			if err != nil {
				fmt.Printf("retrievePolicies: Cannot marshall policy with id " + policy.Id + ": %s", err)
			} else {
				result += string(policyJSON) + ","
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

//==============================================================================================================================
//	 isPolicyRelevantToCaller - Checks if a policy is relevant to the caller
//==============================================================================================================================
func (t *InsuranceChaincode) isPolicyRelevantToCaller(stub shim.ChaincodeStubInterface, policy Policy, caller string, caller_affiliation string) (bool){
	//Super user can see everything
	if caller_affiliation == ROLE_SUPER_USER { return true }

	//Should associate policy with insurer but for now, allow it
	if caller_affiliation == ROLE_INSURER { return true }

	//Is policy owned by caller?
	if caller == policy.Relations.Owner {fmt.Println("Policy owner and caller match, policy is relevant"); return true}

	fmt.Println("Policy owner and caller do not match, policy is not relevant - " + caller + " : " + policy.Relations.Owner)
	return false
}

//==============================================================================================================================
//	 retrieveAllClaimsJSON - Iterates through all claim ids, retreiving each and returning a JSON representation
//==============================================================================================================================
func (t *InsuranceChaincode) retrieveAllClaimsJSON(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	claims := RetrieveAllClaims(stub)

	result := "["

	for _, claim := range claims {

		if t.isClaimRelevantToCaller(stub, claim, caller, caller_affiliation) {

			claimJSON, err := json.Marshal(claim)

			if err != nil {
				fmt.Printf("retrieveAllClaims: Cannot marshall claim with id " + claim.Id + ": %s", err)
			} else {
				result += string(claimJSON) + ","
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

//==============================================================================================================================
//	 isClaimRelevantToCaller - Checks is a claim is relevant to the caller
//==============================================================================================================================
func (t *InsuranceChaincode) isClaimRelevantToCaller(stub shim.ChaincodeStubInterface, claim Claim, caller string, caller_affiliation string) (bool){
	//Super user can see everything
	if caller_affiliation == ROLE_SUPER_USER { return true }

	policy, err := RetrievePolicy(stub, claim.Relations.RelatedPolicy)
	if err != nil {fmt.Printf("Invalid policy, claim is not relevant"); return false}

	if caller_affiliation == ROLE_INSURER {
		//Check policy is associated with the insurer
		if policy.Relations.Insurer == caller {fmt.Printf("Policy insurer and caller match, claim is relevant"); return true}
	}

	//Garage checks
	if caller_affiliation == ROLE_GARAGE {
		//Is claim awaiting a report?
		if claim.Details.Status == STATE_AWAITING_GARAGE_REPORT {
			{fmt.Printf("Awaiting garage report and the caller is a garage, claim is relevant"); return true}
		} else {
			//Is claim awaiting garage work and the garage submitted the report?
			if claim.Details.Status == STATE_AWAITING_GARAGE_WORK_CONFIRMATION &&
				claim.Details.Repair.Garage == caller {
				{fmt.Printf("Awaiting garage work and garage raised report, claim is relevant"); return true}
			} else {
				{fmt.Printf("Garage has not link to claim or incorrect state, claim is not relevant"); return false}
			}
		}
	}

	//Is claim policy owned by caller?
	if caller == policy.Relations.Owner {fmt.Printf("Policy owner and caller match, claim is relevant"); return true}

	fmt.Printf("Policy owner and caller don't match, claim is not relevant")
	return false
}

//==============================================================================================================================
//	 vehicleValueOracleCallback - Callback, called by an oracle when a vehicle value has been retreived
//		args - requestId, vehicleValue
//==============================================================================================================================
func (t *InsuranceChaincode) vehicleValueOracleCallback(stub shim.ChaincodeStubInterface, caller string, callerAffiliation string, args []string) ([]byte, error) {

	if callerAffiliation != ROLE_ORACLE {
		fmt.Printf("vehicleValueCallback: Called from non oracle user")
		return nil, errors.New("vehicleValueCallback: Called from non oracle user")
	}

	vehicleValue, err := strconv.Atoi(args[1])

	if err != nil {	fmt.Printf("vehicleValueCallback: Cannot parse car value %s", err); return nil, errors.New("vehicleValueCallback: Cannot parse car value")}

	bytes, err := stub.GetState(args[0])

	if err != nil {	fmt.Printf("vehicleValueCallback: Cannot read claim id from callback state: %s", err); return nil, errors.New("vehicleValueCallback: Cannot read claim from callback state")}

	claimId := string(bytes)

	claim, err := RetrieveClaim(stub , claimId)

	return t.afterVehicleValueObtainedProcess(stub, claim, vehicleValue)
}

func (t *InsuranceChaincode) queryOracleForVehicleValue(stub shim.ChaincodeStubInterface, claim Claim, vehicle Vehicle) {
	//Store claim id against transaction id so it can be retreived during callback
	stub.PutState(stub.GetTxID(), []byte(claim.Id))
	err := RequestVehicleValuationFromOracle(stub, stub.GetTxID(), vehicle, "vehicleValueOracleCallback")

	if err != nil {
		fmt.Println("Error querying oracle for vehicle value: %s", err);
		fmt.Println("Processing with default value")
		t.afterVehicleValueObtainedProcess(stub, claim, 5555)
	}
}

//==============================================================================================================================
//	 Security Functions
//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *InsuranceChaincode)  get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error){

	user, err := t.get_username(stub)

	if err != nil { user = "dummy-user" }

	affiliation, err := t.check_affiliation(stub);

	if err != nil { affiliation = "dummy-affiliation" }

	fmt.Printf("caller_data: %s %s", user, affiliation)

	return user, affiliation, nil
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *InsuranceChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	username, err := stub.ReadCertAttribute("username");
	if err != nil {
		fmt.Printf("Couldn't get attribute 'username'. Error: %s", err)
		return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	}

	return string(username), nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *InsuranceChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
	affiliation, err := stub.ReadCertAttribute("role");
	if err != nil { return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error()) }
	return string(affiliation), nil
}

//==============================================================================================================================
//	 declareLiability - A function called by a claimant to accept or dispute liability in a multi party accident
//   args{claimId, acceptLiability}
//==============================================================================================================================
func (t *InsuranceChaincode) declareLiability(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addLiabilityDeclaration()")

	if len(args) != 2 {
		fmt.Println("addLiabilityDeclaration: Incorrect number of arguments. Expecting 2 (claimId, acceptLiability)")
		return nil, errors.New("addLiabilityDeclaration: Incorrect number of arguments. Expecting 2 (claimId, acceptLiability)")
	}

	// Get the claim
	var claimId string = args[0]

	claim, err := RetrieveClaim(stub, claimId)

	if err != nil {
		fmt.Printf("\naddLiabilityDeclaration: Failed to retrieve claim Id: %s", err);
		return nil, errors.New("addLiabilityDeclaration: Error retrieving claim with claimId = " + claimId)
	}

	//Check everything is valid
	if !t.shouldAcceptLiabilityDeclarationForClaim(stub, claim, caller) {
		return nil, errors.New("addLiabilityDeclaration: Invalid. Caller: " + caller + ", status:" + claim.Details.Status)
	}

	liable, err := strconv.ParseBool(args[1])
	if err != nil {fmt.Printf("\naddLiabilityDeclaration: Unable to parse boolean acceptLiability: %s", err); return nil, err}

	return t.processDeclareLiability(stub, claim, liable)
}

//==============================================================================================================================
//	 processDeclareLiability - Performs processing on claim when liability has been declared
//==============================================================================================================================
func (t *InsuranceChaincode) processDeclareLiability(stub shim.ChaincodeStubInterface, claim Claim, acceptLiability bool) ([]byte, error) {
	//Liability accepted
	if (acceptLiability) {
		status := STATE_AWAITING_GARAGE_REPORT
		//Update this claims status
		claim.Details.Status = status
		_, err := SaveClaim(stub, claim)
		if err != nil {fmt.Printf("\naddLiabilityDeclaration: Unable to save claim: %s", err); return nil, err}

		//Update the linked claims status'
		for _, claimId := range claim.Relations.LinkedClaims {
			otherClaim, err := RetrieveClaim(stub, claimId)
			if err != nil {fmt.Printf("\naddLiabilityDeclaration: Unable to retrieve claim with id: " + claimId + ": %s", err); return nil, err}

			otherClaim.Details.Status = status
			_, err = SaveClaim(stub, otherClaim)
			if err != nil {fmt.Printf("\naddLiabilityDeclaration: Unable to save claim: %s", err); return nil, err}
		}

	} else {
		//TODO Some kind of dispute state
		return nil, errors.New("Liability dispute is currently unsupported")
	}
	return nil, nil
}
//==============================================================================================================================
//	 shouldAcceptLiabilityDeclarationForClaim - Checks if liability can be declared for the specified claim and caller
//==============================================================================================================================
func (t *InsuranceChaincode) shouldAcceptLiabilityDeclarationForClaim(stub shim.ChaincodeStubInterface, claim Claim, caller string) (bool){
	if (claim.Details.Status != STATE_AWAITING_LIABILITY_ACCEPTANCE) {
		fmt.Println("Claim in wrong state for liability declaration: " + claim.Id)
		return false
	}

	policy, err := RetrievePolicy(stub, claim.Relations.RelatedPolicy)

	if err != nil { fmt.Printf("shouldAcceptLiabilityDeclarationForClaim: Cannot retrieve policy: %s", err); return false}

	if policy.Relations.Owner != caller {
		fmt.Println("Caller is not the owner of the claim: " + caller + " : " + policy.Relations.Owner)
		return false
	}

	return true
}

//==============================================================================================================================
//	 addGarageReport - This method adds the garage report's details into the claim
//   args{claimId, garage, estimated_cost,  writeOff, Note}
//==============================================================================================================================
func (t *InsuranceChaincode) addGarageReport(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {

	fmt.Println("running addGarageReport()")

	if len(args) != 5 {
		fmt.Println("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 5 (claimId, estimated_cost, writeOff, Note, reg_number)")
		return nil, errors.New("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 5 (claimId, estimated_cost, writeOff, Note, reg_number)")
	}

	// Does a claim exist for this vehicle?
	var theClaim Claim
	var claimId string = args[0]

	theClaim, err := RetrieveClaim(stub , claimId)

	if err != nil {	fmt.Printf("\nADD_GARAGE_REPORT: Failed to retrieve claim Id: %s", err); return nil, errors.New("ADD_GARAGE_REPORT: Error retrieving claim with claimId = " + claimId) }

	if !t.shouldAcceptGarageReportForClaim(stub, theClaim, caller, caller_affiliation) {
		return nil, errors.New("ADD_GARAGE_REPORT: Invalid. Affilliation: " + caller_affiliation + ", status:" + theClaim.Details.Status)
	}

	if !t.isVehicleValidForClaim(stub, theClaim, args[4]){
		fmt.Printf("ADD_GARAGE_REPORT: Invalid Registration: %s\n", args[4])
		return nil, errors.New("ADD_GARAGE_REPORT: Invalid Registration: " + args[4])
	}

	var report ClaimDetailsClaimGarageReport
	report, err = NewGarageReport(caller, args[1], args[2], args[3])
	theClaim.Details.Report = report

	theClaim.Details.Status = STATE_PENDING_AFTER_REPORT_DECISION

	SaveClaim(stub, theClaim)

	t.afterReportProcess(stub, theClaim)

	return nil, nil
}

func (t *InsuranceChaincode) shouldAcceptGarageReportForClaim(stub shim.ChaincodeStubInterface, claim Claim, caller string, caller_affiliation string) (bool) {
	//Super user can do anything
	if caller_affiliation == ROLE_SUPER_USER { return true }

	//Is caller a garage?
	if caller_affiliation != ROLE_GARAGE {
		fmt.Printf("Non Garage user attempting to add a garage report!!")
		return false
	}

	// Is the Claim in a valid state to receive a garage report?
	if STATE_AWAITING_GARAGE_REPORT != claim.Details.Status {
		fmt.Printf("ADD_GARAGE_REPORT: claim in invalid state for garage report, with Id: %s\n", claim.Id)
		return false
	}

	//Is the garage an approved garage?
	if !t.isApprovedGarage(stub, caller) {
		fmt.Printf("ADD_GARAGE_REPORT: Garage is not approved: %s\n", caller)
		return false
	}

	return true
}

//=========================================================================================
// This Function process the claim after a report has arrived
//=========================================================================================
func (t *InsuranceChaincode) afterReportProcess(stub shim.ChaincodeStubInterface, claim Claim) {
	//Get vehicle
	policy, _ := RetrievePolicy(stub, claim.Relations.RelatedPolicy)
	vehicle, _ := RetrieveVehicle(stub, policy.Relations.Vehicle)

	t.queryOracleForVehicleValue(stub, claim, vehicle)
}
//=========================================================================================
// This Function process the claim after a report has arrived
//=========================================================================================
func (t *InsuranceChaincode) afterVehicleValueObtainedProcess(stub shim.ChaincodeStubInterface, theClaim Claim, vehicleValue int) ([]byte, error) {
	fmt.Println("running afterVehicleValueObtainedProcess()")

	if theClaim.Details.Status == STATE_PENDING_AFTER_REPORT_DECISION{

		if theClaim.Details.Report.WriteOff || theClaim.Details.Report.Estimate > (vehicleValue * 50 / 100) {
			//process total_loss
			t.processTotalLoss(stub, theClaim, vehicleValue)
		}else{
			theClaim.Details.Status = STATE_ORDER_GARAGE_WORK
			SaveClaim(stub, theClaim)
		}
	} else {
		fmt.Println("AFTER_VALUE_PROCESS: Claim in invalid state: " + theClaim.Details.Status);
		return nil, errors.New("AFTER_VALUE_PROCESS: Claim in invalid state: " + theClaim.Id)
	}
		
	return nil, nil
}

//=========================================================================================
//
//=========================================================================================
func (t *InsuranceChaincode) processTotalLoss(stub shim.ChaincodeStubInterface, theClaim Claim, vehicleValue int) ([]byte, error) {

	fmt.Println("running processTotalLoss()")
    if &theClaim != nil{
		var settlement ClaimDetailsSettlement
		settlement.Decision = TOTAL_LOSS
		settlement.Dispute = false
		theClaim.Details.Settlement = settlement

		//theClaim.Details.Status = STATE_TOTAL_LOSS_ESTABLISHED
		theClaim.Details.Status = STATE_AWAITING_CLAIMANT_CONFIRMATION
		theClaim.Details.Settlement.TotalLoss.CarValueEstimate = vehicleValue
		SaveClaim(stub, theClaim)
	}else {fmt.Printf("PROCESS_TOTAL_LOSS: Error can not process Total loss on unexisting claim\n"); return nil, errors.New("PROCESS_TOTAL_LOSS: Error can not process Total loss on unexisting claim") }
	return nil, nil
}

//=========================================================================================
// This Function set the claimant aggreement in the payout
//=========================================================================================
func (t *InsuranceChaincode) agreePayoutAmount(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {
	fmt.Println("running agreePayoutAmount()")

    //
	//TODO - Check the Security  check that this claimant has the right to view and update this claim
	//
	if len(args) != 2 {
		fmt.Println("AGREE_PAYOUT_AMOUNT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
		return nil, errors.New("AGREE_PAYOUT_AMOUNT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
	}

    var theClaim Claim

	theClaim, err := RetrieveClaim(stub , args[0])

	if err != nil {	fmt.Printf("\nAGREE_PAYOUT_AMOUNT: Failed to retrieve claim Id: %s", err); return nil, errors.New("AGREE_PAYOUT_AMOUNT: Error retrieving claim with claimId = " + args[0]) }

	if theClaim.Details.Status != STATE_AWAITING_CLAIMANT_CONFIRMATION{
		fmt.Println("AGREE_PAYOUT_AMOUNT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
		return nil, errors.New("AGREE_PAYOUT_AMOUNT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
	}

	var acceptDeny bool
	acceptDeny, err = strconv.ParseBool(args[1])
	if err != nil {fmt.Printf("AGREE_PAYOUT_AMOUNT Error: invalid value passed for agreement: %s\n", err); return nil, errors.New("Invalid value passed for agreement")}

	if acceptDeny {
		//TODO this will only work for total loss
		theClaim.Details.Settlement.TotalLoss.CustomerAgreedValue = theClaim.Details.Settlement.TotalLoss.CarValueEstimate
		theClaim.Details.Status = STATE_SETTLED
		theClaim.Details.Settlement.Dispute = false

		policy, err := RetrievePolicy(stub, theClaim.Relations.RelatedPolicy)

		if err != nil {
			fmt.Printf("AGREE_PAYOUT_AMOUNT: Error getting policy with id %s", theClaim.Relations.RelatedPolicy);
			return nil, errors.New("Policy doesnt exist");
		}

		//Add pending payments to claim
		theClaim, err = t.addPendingPaymentsToClaim(stub, theClaim, policy)

		//TODO Assuming one linked claim for now
		var linkedClaim string
		if len(theClaim.Relations.LinkedClaims) > 0 { linkedClaim = theClaim.Relations.LinkedClaims[0]}

		event := NewClaimSettledEvent(theClaim.Id, theClaim.Relations.RelatedPolicy, linkedClaim)

		eventBytes, err := json.Marshal(event);
		if err != nil {fmt.Printf("AGREE_PAYOUT_AMOUNT Error: Unable to add pending payment: %s\n", err); return nil, err}

		if (err != nil) {
			fmt.Printf("\nAGREE_PAYOUT_AMOUNT: Unable to parse event, continuing without emitting", err)
		} else {
			//Emit claim settled event
			stub.SetEvent(event.Type, eventBytes)
		}
	} else {
		theClaim.Details.Settlement.Dispute = true
	}
	SaveClaim(stub, theClaim)

	return nil, nil
}

//=========================================================================================
// Adds pending payments to a claim
//=========================================================================================
func (t *InsuranceChaincode) addPendingPaymentsToClaim(stub shim.ChaincodeStubInterface, claim Claim, policy Policy) (Claim, error) {

	if claim.Details.Status != STATE_SETTLED{
		fmt.Println("addPendingPaymentsToClaim: Unexpected input for this STATE")
		return claim, errors.New("addPendingPaymentsToClaim: Unexpected input for this STATE")
	}

	paymentAmount := t.calculatePaymentAmount(claim, policy)

	claim, err := t.addPendingPaymentToClaimant(claim, policy, paymentAmount)
	if err != nil {fmt.Printf("addPendingPaymentsToClaim: Unable to add claimant payment: %s", err); return claim, err}

	//If we're not liable, there needs to be a pending payment added from the linked claim insurer to this parties insurer
	if (claim.Details.Incident.Type == MULTIPLE_PARTIES && !claim.Details.IsLiable) {
		claim, err = t.addPendingPaymentFromOtherPartyInsurer(stub, claim, policy, paymentAmount)
		if err != nil { return claim, err }
	}

	return claim, nil
}

func (t *InsuranceChaincode) calculatePaymentAmount(claim Claim, policy Policy) (int){
	excess := policy.Details.Excess

	if (claim.Details.Incident.Type == MULTIPLE_PARTIES) {
		//If the party is not liable in a multiple party claim, then there should be no excess
		if (!claim.Details.IsLiable) {
			excess = 0
		}
	}

	return claim.Details.Settlement.TotalLoss.CustomerAgreedValue - excess
}

//=========================================================================================
// This Function adds a pending payment from insurer to claimant
//=========================================================================================
func (t *InsuranceChaincode) addPendingPaymentToClaimant(theClaim Claim, policy Policy, amount int) (Claim, error) {
	fmt.Println("running addPendingPaymentToClaimant()")

	payment := NewClaimDetailsSettlementPayment(PAYMENT_TYPE_CLAIMANT,
		policy.Relations.Owner, PAYMENT_TYPE_INSURER, policy.Relations.Insurer, amount, STATE_NOT_PAID)

	theClaim.AddPayment(payment)

	return theClaim, nil
}

//=========================================================================================
// This Function adds a pending payment from insurer to insurer
//=========================================================================================
func (t *InsuranceChaincode) addPendingPaymentFromOtherPartyInsurer(stub shim.ChaincodeStubInterface, claim Claim, policy Policy, amount int) (Claim, error) {
	fmt.Println("running addPendingPaymentFromOtherPartyInsurer()")

	liableClaim, err := claim.GetLiableClaim(stub)
	if err != nil { return claim, err }

	liablePolicy, err := RetrievePolicy(stub, liableClaim.Relations.RelatedPolicy)
	if err != nil { return claim, err}

	payment := NewClaimDetailsSettlementPayment(PAYMENT_TYPE_INSURER,
		policy.Relations.Insurer, PAYMENT_TYPE_INSURER, liablePolicy.Relations.Insurer, amount, STATE_NOT_PAID)
	claim.AddPayment(payment)

	//Double enter the payment to the liable claim list of payments
	liableClaim.AddPayment(payment)
	_, err = SaveClaim(stub, liableClaim)

	if err != nil { return claim, err }

	//No need to save the other claim as its saved elsewhere

	return claim, nil
}

//=========================================================================================
// This Function marks the claim as paid
//=========================================================================================
func (t *InsuranceChaincode) confirmPaidOut(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error){

	fmt.Println("running confirmPaidOut()")

	if caller_affiliation != ROLE_INSURER {
		fmt.Printf("\nconfirmPaidOut: Caller is not an insurer")
		return nil, errors.New("\nconfirmPaidOut: Caller is not an insurer")
	}

	if len(args) != 2 {
		fmt.Println("CONFIRM_PAID_OUT: Incorrect number of arguments. Expecting 2 (claimId, paymentId)")
		return nil, errors.New("CONFIRM_PAID_OUT: Incorrect number of arguments. Expecting 2 (claimId, paymentId)")
	}

	theClaim, err := RetrieveClaim(stub , args[0])

	if err != nil {fmt.Println("Unable to retrieve claim with id: " + args[0]); return nil, err}

	if theClaim.Details.Status != STATE_SETTLED{
		fmt.Println("CONFIRM_PAID_OUT: Unexpected input for this STATE")
		return nil, errors.New("CONFIRM_PAID_OUT: Unexpected input for this STATE")
	}

	payment, err := theClaim.GetPayment(args[1])

	//Check that the caller is the sender
	if (payment.Sender != caller){
		fmt.Println("CONFIRM_PAID_OUT: Caller is not the sender of the payment")
		return nil, errors.New("CONFIRM_PAID_OUT: Caller is not the sender of the payment")
	}

	if err != nil { return nil, err}

	payment.Status = STATE_PAID
	theClaim.UpdatePayment(payment)

	_, err = SaveClaim(stub, theClaim)

	err = t.updateLinkedPayments(stub, theClaim, payment)
	if err != nil { fmt.Printf("Unable to update linked payments: %s", err) }
	
	return nil, err
}

//=========================================================================================
// Updates the status of any matching payments on linked claims
//=========================================================================================
func (t *InsuranceChaincode) updateLinkedPayments(stub shim.ChaincodeStubInterface, claim Claim, payment ClaimDetailsSettlementPayment) (error){

	fmt.Println("running updateLinkedPayments()")

	//If this is a payment to an insurer then update the payment on the linked claim
	if (payment.RecipientType == PAYMENT_TYPE_INSURER) {
		for _, claimId := range claim.Relations.LinkedClaims {
			linkedClaim, err := RetrieveClaim(stub, claimId)
			if err != nil{ return err }

			//Check if there is a payment with the same parameters TODO: Hacky. There should be a proper link I think
			for _, linkedPayment := range linkedClaim.Details.Settlement.Payments {
				if linkedPayment.Sender == payment.Sender &&
						linkedPayment.Recipient == payment.Recipient &&
						linkedPayment.Amount == payment.Amount {
					//We've found a matching payment...update
					linkedPayment.Status = STATE_PAID
					linkedClaim.UpdatePayment(linkedPayment)
					_, err = SaveClaim(stub, linkedClaim)

					//Fire Insurer payment paid event
					event := NewInsurerPaymentPaidEvent(linkedClaim.Id, linkedClaim.Relations.RelatedPolicy)
					eventBytes, _ := json.Marshal(event)
					stub.SetEvent(event.Type, eventBytes)

					if err != nil { return err }
				}
			}
		}
	}

	return nil
}

//===============================================================================
// This method checks whether all the payement out are sent
//===============================================================================
func (t *InsuranceChaincode) isApprovedGarage(stub shim.ChaincodeStubInterface, garage string) bool {

	approvedGarages, err := RetrieveApprovedGarages(stub)
	if err != nil {	fmt.Printf("IS_APPROVED_GARAGE: Unable to retrieve the approved garages: %s", err); return false}

    for _, appGarage := range approvedGarages.Garages {
		fmt.Printf("Approved Garage: %s\n", appGarage)
        if appGarage == garage {
            return true
        }
    }
    return false
}

//===============================================================================
// This method Closes the claim
//===============================================================================
func (t *InsuranceChaincode) closeClaim(stub shim.ChaincodeStubInterface,  caller string, caller_affiliation string, args []string) ([]byte, error){

	if len(args) != 1 {
		fmt.Println("CLOSE_CLAIM: Incorrect number of arguments. Expecting 1 (claimId)")
		return nil, errors.New("CLOSE_CLAIM: Incorrect number of arguments. Expecting 1 (claimId)")
	}

	//Only an insurer can close the claim
	if (caller_affiliation != ROLE_INSURER) {
		fmt.Println("CLOSE_CLAIM: Only an insurer can close the claim")
		return nil, errors.New("CLOSE_CLAIM: Only an insurer can close the claim")
	}

    var theClaim Claim

	theClaim, err := RetrieveClaim(stub , args[0])

	if err != nil {	fmt.Printf("\nCLOSE_CLAIM: Failed to retrieve claim Id: %s", err); return nil, errors.New("CLOSE_CLAIM: Error retrieving claim with claimId = " + args[0]) }

	if theClaim.Details.Status != STATE_SETTLED{
		fmt.Println("CLOSE_CLAIM: Incorrect STATE, Claim can not be closed if not SETTLED")
		return nil, errors.New("CLOSE_CLAIM: Incorrect STATE, Claim can not be closed if not SETTLED")
	}

	if TOTAL_LOSS == theClaim.Details.Settlement.Decision{
		return t.closeTotalLossClaim(stub, theClaim)
	}
	fmt.Printf("\nCLOSE_CLAIM: Unsuported Close Operation. Currently the System can only Close Total_Loss Claims.\n")
	return nil, errors.New("CLOSE_CLAIM: Unsuported Close Operation. Currently the System can only Close Total_Loss Claims.")
}

//============================================================================================================
//
//============================================================================================================
func (t *InsuranceChaincode) closeTotalLossClaim(stub shim.ChaincodeStubInterface,  theClaim Claim) ([]byte, error){

	if !theClaim.AreAllPaymentsPaid() {
		fmt.Println("CLOSE_CLAIM: Error: open payment out. Total_Loss Claim can not be closed with open payment out")
		return nil, errors.New("CLOSE_CLAIM: Error: open payment out. Total_Loss Claim can not be closed with open payment out")
	}

	theClaim.Details.Status = STATUS_CLOSED
	SaveClaim(stub, theClaim)

	return nil, nil
}

func (t *InsuranceChaincode) isVehicleValidForClaim(stub shim.ChaincodeStubInterface,  theClaim Claim, vehicleReg string)(bool){

	//Check policy exists
	policy, err := RetrievePolicy(stub, theClaim.Relations.RelatedPolicy);

	if err != nil {
		fmt.Printf("isVehicleValidForClaim: Error getting policy with id %s", theClaim.Relations.RelatedPolicy)
		return false
	}

	//Check policy vehicle matches the registration provided
	if policy.Relations.Vehicle == vehicleReg{
		fmt.Printf("isVehicleValidForClaim: Policy Vehicle is correct %s", vehicleReg)
		return true;
	}
	fmt.Printf("isVehicleValidForClaim: Policy Vehicle is incorrect %s", vehicleReg)
	return false
}
