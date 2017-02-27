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

	//Should associate policy with insurer but for now, allow it
	if caller_affiliation == ROLE_INSURER { return true }

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
	policy, err := RetrievePolicy(stub, claim.Relations.RelatedPolicy)
	if err != nil {fmt.Printf("Invalid policy, claim is not relevant"); return false}

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
		fmt.Println("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
		return nil, errors.New("ADD_GARAGE_REPORT: Incorrect number of arguments. Expecting 2 (claimId, agreement)")
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

		//Add pending payment to claim
		theClaim, err = t.addPendingPayment(stub, theClaim)

		policy, _ := RetrievePolicy(stub, theClaim.Relations.RelatedPolicy)
		event := NewClaimSettledEvent(theClaim.Id, theClaim.Relations.RelatedPolicy, policy.Relations.Owner)

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
// This Function marks the claim as paid
//=========================================================================================
func (t *InsuranceChaincode) addPendingPayment(stub shim.ChaincodeStubInterface, theClaim Claim) (Claim, error) {
	fmt.Println("running addPendingPayment()")

	if theClaim.Details.Status != STATE_SETTLED{
		fmt.Println("APPROVE_PAYMENT_OUT: Unexpected input for this STATE")
		return theClaim, errors.New("APPROVE_PAYMENT_OUT: Unexpected input for this STATE")
	}

	policy, err := RetrievePolicy(stub, theClaim.Relations.RelatedPolicy);

	if err != nil {
		fmt.Printf("APPROVE_PAYMENT_OUT: Error getting policy with id %s", theClaim.Relations.RelatedPolicy);
		return theClaim, errors.New("Policy doesnt exist");
	}

	var payment ClaimDetailsSettlementPayment
	
	var Amount = theClaim.Details.Settlement.TotalLoss.CustomerAgreedValue - policy.Details.Excess

	payment = NewClaimDetailsSettlementPayment(RECIPIENT_TYPE_CLAIMANT, policy.Relations.Owner, Amount, STATE_NOT_PAID)

	theClaim.Details.Settlement.Payments = make([]ClaimDetailsSettlementPayment, 1)
	theClaim.Details.Settlement.Payments[0] = payment

	return theClaim, nil
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

	if len(args) != 1 {
		fmt.Println("CONFIRM_PAID_OUT: Incorrect number of arguments. Expecting 1 (claimId)")
		return nil, errors.New("CONFIRM_PAID_OUT: Incorrect number of arguments. Expecting 1 (claimId)")
	}

	theClaim, err := RetrieveClaim(stub , args[0])

	if err != nil {fmt.Println("Unable to retrieve claim with id: " + args[0]); return nil, err}

	if theClaim.Details.Status != STATE_SETTLED{
		fmt.Println("CONFIRM_PAID_OUT: Unexpected input for this STATE")
		return nil, errors.New("CONFIRM_PAID_OUT: Unexpected input for this STATE")
	}
	
	// Assumption: there should be only one payment
	theClaim.Details.Settlement.Payments[0].Status = STATE_PAID
	_, err = SaveClaim(stub, theClaim)
	
	return nil, err
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

	if !t.isAllPaidout(theClaim){
		fmt.Println("CLOSE_CLAIM: Error: open payment out. Total_Loss Claim can not be closed with open payment out")
		return nil, errors.New("CLOSE_CLAIM: Error: open payment out. Total_Loss Claim can not be closed with open payment out")
	}

	theClaim.Details.Status = STATUS_CLOSED
	SaveClaim(stub, theClaim)

	return nil, nil
}

//===========================================================================================================
//  This function checks whether all the payment out are sent
//===========================================================================================================
func (t *InsuranceChaincode) isAllPaidout(theClaim Claim) bool {

	if theClaim.Details.Settlement.Payments != nil{
		// Assumption: there should be only one payment out
		if STATE_PAID == theClaim.Details.Settlement.Payments[0].Status{
			return true;
		}
	}

	return false
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
