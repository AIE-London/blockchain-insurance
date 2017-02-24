package main

//==============================================================================================================================
//	ClaimSettledEvent - Defines the structure for a claim settled event.
//==============================================================================================================================
type ClaimSettledEvent struct {
	Type		string				`json:"eventType"`
	ClaimId		string				`json:"claimId"`
	PolicyId	string				`json:"policyId"`
	User		string				`json:"user"`
}

//==============================================================================================================================
//	 Event type codes
//==============================================================================================================================
const EVENT_TYPE_CLAIM_SETTLED = "ClaimSettled";

//=================================================================================================================================
//	 NewClaimSettledEvent	-	Constructs a new ClaimSettledEvent
//=================================================================================================================================
func NewClaimSettledEvent(claimId string, policyId string, user string) (ClaimSettledEvent) {
	var event ClaimSettledEvent

	event.Type = EVENT_TYPE_CLAIM_SETTLED
	event.ClaimId = claimId
	event.PolicyId = policyId
	event.User = user

	return event
}
