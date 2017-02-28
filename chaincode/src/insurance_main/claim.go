package main

import (
	"errors"
	"fmt"
	"strconv"
)

//==============================================================================================================================
//	Claim - Defines the structure for a Claim object.
//==============================================================================================================================
type Claim struct {
	Id		string				`json:"id"`
	Type		string			`json:"type"`
	Details		ClaimDetails	`json:"details"`
	Relations 	ClaimRelations	`json:"relations"`
}

//==============================================================================================================================
//	ClaimDetails - Defines the structure for a ClaimDetails object.
//==============================================================================================================================
type ClaimDetails struct {
	Status		string							`json:"status"`
	Description	string							`json:"description"`
	Incident	ClaimDetailsIncident			`json:"incident"`
	Report		ClaimDetailsClaimGarageReport	`json:"report"`
	Repair		RepairWorkOrder             	`json:"repair"`
	Settlement	ClaimDetailsSettlement			`json:"settlement"`
	IsLiable	bool							`json:"liable"`
}

//==============================================================================================================================
//	ClaimDetailsIncident - Defines the structure for a ClaimDetailsIncident object.
//==============================================================================================================================
type ClaimDetailsIncident struct {
	Date	string	`json:"date"`
	Type	string	`json:"type"`
}

//==============================================================================================================================
//	ClaimDetailsClaimGarageReport - Defines the structure for a ClaimDetailsClaimGarageReport object.
//==============================================================================================================================
type ClaimDetailsClaimGarageReport struct {
	Garage		string	`json:"garage"`
	Estimate	int		`json:"estimate"`
	WriteOff	bool	`json:"writeOff"`
	Notes		string	`json:"notes"`
}

//==============================================================================================================================
//	ClaimDetailsSettlement - Defines the structure for a ClaimDetailsSettlement object.
//==============================================================================================================================
type ClaimDetailsSettlement struct {
	Decision	string							`json:"decision"`
	Dispute		bool							`json:"dispute"`
	TotalLoss	ClaimDetailsSettlementTotalLoss	`json:"totalLoss"`
	Payments	[]ClaimDetailsSettlementPayment	`json:"payments"`
}

//==============================================================================================================================
//	ClaimDetailsSettlementTotalLoss - Defines the structure for a ClaimDetailsSettlementTotalLoss object.
//==============================================================================================================================
type ClaimDetailsSettlementTotalLoss struct {
	CarValueEstimate	int	`json:"carValueEstimate"`
	CustomerAgreedValue	int	`json:"customerAgreedValue"`
}

//==============================================================================================================================
//	ClaimDetailsSettlementPayment - Defines the structure for a ClaimDetailsSettlementPayment object.
//==============================================================================================================================
type ClaimDetailsSettlementPayment struct {
	RecipientType	string	`json:"recipientType"`
	Recipient	string		`json:"recipient"`
	Amount		int			`json:"amount"`
	Status		string		`json:"status"`
	SenderType	string		`json:"senderType"`
	Sender		string		`json:"sender"`
}

//==============================================================================================================================
//	ClaimRelations - Defines the structure for a ClaimRelations object.
//==============================================================================================================================
type ClaimRelations struct {
	RelatedPolicy	string		`json:"relatedPolicy"`
	OtherPartyReg	string		`json:"otherPartyRegistration"`
	LinkedClaims	[]string	`json:"linkedClaims"`
}

//=============================================================================================================================
//  RepairWorkOrder - Defines the structure for a RepairWorkOrder object.
//     RepairWorkOrder can be in the following workstate : STATUS_OPEN - STATUS_CLOSED - STATE_IN_PROGRESS
//==============================================================================================================================
type RepairWorkOrder struct {
	Garage  	           string	`json:"garages"`
	WorkOrderSeq           string   `json:"workOrderSeq"`
	Description            string   `json:"description"`
	StartDate              string   `json:"startDate"`
	EndDate                string   `json:"endDate"`
	EstimatedRepairCost    int      `json:"estimatedRepairCost"`
	ActualRepairCost       int      `json:"actualRepairCost"`
	WorkStatus             string   `json:"workStatus"`
}

//==============================================================================================================================
//	 Claim Status types
//==============================================================================================================================
const   STATE_AWAITING_POLICE_REPORT                = "awaiting_police_report"
const	STATE_AWAITING_LIABILITY_ACCEPTANCE			= "awaiting_liability_acceptance"
const   STATE_AWAITING_GARAGE_REPORT                = "awaiting_garage_report"
const   STATE_PENDING_AFTER_REPORT_DECISION         = "pending_decision"
const   STATE_TOTAL_LOSS_ESTABLISHED                = "total_loss_established"
const   STATE_ORDER_GARAGE_WORK                     = "garage_work_ordered"
const   STATE_AWAITING_GARAGE_WORK_CONFIRMATION     = "awaiting_garage_work"
const   STATE_AWAITING_CLAIMANT_CONFIRMATION        = "awaiting_claimant_confirmation"
const   STATE_SETTLED  			                    = "settled"
const	STATUS_OPEN									= "open"
const	STATUS_CLOSED								= "closed"
const   STATE_IN_PROGRESS                           = "in_progress"
const   STATE_NOT_PAID                              = "not_paid"
const   STATE_PAID                      = "paid"

//==============================================================================================================================
//	 Claim Type types - TODO Flesh these out. TODO Following IBM sample, but should these be enums?
//==============================================================================================================================
const   SINGLE_PARTY                =  "single_party"
const   MULTIPLE_PARTIES  			=  "multiple_parties"

//==============================================================================================================================
//	 Decision types - 
//==============================================================================================================================
const   TOTAL_LOSS          =  "total_loss"
const   LIABILITY  			=  "liability"
const   NOT_AT_FAULT        =  "not_at_fault"

//
//
//
const PAYMENT_TYPE_CLAIMANT    = "claimant"
const PAYMENT_TYPE_THIRDPARTY  = "third_party"
const PAYMENT_TYPE_GARAGE      = "garage"
const PAYMENT_TYPE_INSURER     = "insurer"

//=================================================================================================================================
//	 newClaim	-	Constructs a new claim
//=================================================================================================================================
func NewClaim(id string, relatedPolicy string, description string, incidentDate string, incidentType string) (Claim) {

	var claim Claim

	claim.Id = id
	claim.Type = "claim"

	claim.Relations.RelatedPolicy = relatedPolicy
	claim.Relations.LinkedClaims = []string{}
	claim.Details.Description = description
	claim.Details.Incident.Date = incidentDate
	claim.Details.Incident.Type = incidentType

	return claim
}

//================================================================================================================================================
// newGarageReport  create a new Garage Report
//================================================================================================================================================
func NewGarageReport(Garage string, EstimateStr string,  WriteOffStr string, Notes string) (ClaimDetailsClaimGarageReport, error) {
	 
	var report ClaimDetailsClaimGarageReport

	var Estimate int
	Estimate, err := strconv.Atoi(EstimateStr)
	if err != nil {fmt.Printf("\nNewGarageReport Error: invalid value passed for Estimate: %s", err); return report, errors.New("Invalid value passed for Estimate")}
	
	var WriteOff bool
	WriteOff, err = strconv.ParseBool(WriteOffStr)
	if err != nil {fmt.Printf("\nNewGarageReport Error: invalid value passed for WriteOff: %s", err); return report, errors.New("Invalid value passed for WriteOff")}
	
	report.Garage    = Garage
	report.Estimate	 = Estimate
	report.WriteOff	 = WriteOff
	report.Notes     = Notes
	
	return report, nil
}


//================================================================================================================================================
//  NewClaimDetailsSettlementPayment: creates a new NewClaimDetailsSettlementPayment
//================================================================================================================================================
func NewClaimDetailsSettlementPayment(RecipientType	string, Recipient string, SenderType string, Sender string, Amount int, Status string)(ClaimDetailsSettlementPayment){

	var payment ClaimDetailsSettlementPayment 
	
	payment.RecipientType = RecipientType
	payment.Recipient     = Recipient
	payment.SenderType	  = SenderType
	payment.Sender        = Sender
	payment.Amount        = Amount
	payment.Status        = Status

	return payment
}

