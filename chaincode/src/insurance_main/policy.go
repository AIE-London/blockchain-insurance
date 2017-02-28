package main

//==============================================================================================================================
//	Policy - Defines the structure for a policy object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Policy struct {
	Id		string					`json:"id"`
	Type		string				`json:"type"`
	Details		PolicyDetails		`json:"details"`
	Relations 	PolicyRelations		`json:"relations"`
}

//==============================================================================================================================
//	PolicyDetails - Defines the structure for the PolicyDetails object.
//==============================================================================================================================


type PolicyDetails struct {
	StartDate	string			`json:"startDate"`
	EndDate		string			`json:"endDate"`
	Excess		int				`json:"excess"`
}

//==============================================================================================================================
//	PolicyRelations - Defines the structure for a PolicyRelations object.
//==============================================================================================================================
type PolicyRelations struct {
	Owner		string			`json:"owner"`
	Insurer		string			`json:"insurer"`
	Vehicle		string			`json:"vehicle"`
	Claims		[]string		`json:"claims"`

}

//=================================================================================================================================
//	 newPolicy	-	Constructs a new policy
//=================================================================================================================================
func NewPolicy(id string, owner string, insurer string, startDate string, endDate string, excess int, vehicleReg string) (Policy) {
	var policy Policy

	policy.Type = "policy"
	policy.Id = id

	policy.Details.StartDate = startDate
	policy.Details.EndDate = endDate
	policy.Details.Excess = excess

	policy.Relations.Owner = owner
	policy.Relations.Insurer = insurer
	policy.Relations.Vehicle = vehicleReg

	return policy
}
