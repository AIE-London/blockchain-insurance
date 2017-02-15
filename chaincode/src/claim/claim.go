package claim

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
	Repair		ClaimDetailsClaimGarageReport	`json:"repair"`
	Settlement	ClaimDetailsSettlement			`json:"settlement"`
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
	Actual		int		`json:"actual"`
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
}

//==============================================================================================================================
//	ClaimRelations - Defines the structure for a ClaimRelations object.
//==============================================================================================================================
type ClaimRelations struct {
	RelatedPolicy	string	`json:"relatedPolicy"`
}

//==============================================================================================================================
//	 Claim Status types - TODO Flesh these out. TODO Following IBM sample, but should/could these be enums?
//==============================================================================================================================
const   STATE_AWAITING_POLICE_REPORT                = "awaiting_police_report"
const   STATE_AWAITING_GARAGE_REPORT                = "awaiting_garage_report"
const   STATE_AWAITING_GARAGE_WORK_CONFIRMATION     = "awaiting_garage_work"
const   STATE_SETTLED  			                    = "settled"
const	STATUS_OPEN									= "open"
const	STATUS_CLOSED								= "closed"

//==============================================================================================================================
//	 Claim Type types - TODO Flesh these out. TODO Following IBM sample, but should these be enums?
//==============================================================================================================================
const   SINGLE_PARTY                =  "single_party"
const   MULTIPLE_PARTIES  			=  "multiple_parties"

//=================================================================================================================================
//	 newClaim	-	Constructs a new claim
//=================================================================================================================================
func New(id string, relatedPolicy string, description string, incidentDate string, incidentType string) (Claim) {

	var claim Claim

	claim.Id = id
	claim.Type = "claim"

	claim.Relations.RelatedPolicy = relatedPolicy
	claim.Details.Description = description
	claim.Details.Incident.Date = incidentDate
	claim.Details.Incident.Type = incidentType

	claim.Details.Status = STATE_AWAITING_GARAGE_REPORT

	return claim
}

//================================================================================================================================================
// newGarageReport  create a new Garage Report
//================================================================================================================================================
func NewGarageReport(Garage string, EstimateStr string, ActualStr string, WriteOffStr string, Notes string) (ClaimDetailsClaimGarageReport, error) {
	 
	var report ClaimDetailsClaimGarageReport

	var Estimate int
	Estimate, err := strconv.Atoi(EstimateStr)
	if err != nil {fmt.Printf("\nevaluate Repair Error: invalid value passed for Estimate: %s", err); return report, errors.New("Invalid value passed for Estimate")}
	
	var Actual int
	Actual, err	= strconv.Atoi(ActualStr)
	if err != nil {fmt.Printf("\nevaluate Repair Error: invalid value passed for Actual: %s", err); return report, errors.New("Invalid value passed for Actual")}
	
	var WriteOff bool
	WriteOff, err = strconv.ParseBool(WriteOffStr)
	if err != nil {fmt.Printf("\nevaluate Repair Error: invalid value passed for WriteOff: %s", err); return report, errors.New("Invalid value passed for WriteOff")}
	
	report.Garage    = Garage
	report.Estimate	 = Estimate
	report.Actual	 = Actual
	report.WriteOff	 = WriteOff
	report.Notes     = Notes
	
	return report, nil
}
