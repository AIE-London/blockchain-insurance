package main

//==============================================================================================================================
//	Vehicle - Defines the structure for a Vehicle object.
//==============================================================================================================================
type Vehicle struct {
	Id		string			`json:"id"`
	Type		string			`json:"type"`
	Details		VehicleDetails		`json:"details"`
}

//==============================================================================================================================
//	VehicleDetails - Defines the structure for a VehicleDetails object.
//==============================================================================================================================
type VehicleDetails struct {
	Make			string		`json:"make"`
	Model			string		`json:"model"`
	Registration		string		`json:"registration"`
	Year			string		`json:"year"`
	Mileage			int		`json:"mileage"`
}

//=================================================================================================================================
//	 New Vehicle	-	Constructs a new vehicle
//=================================================================================================================================
func NewVehicle(id string, make string, model string, registration string, year string, mileage int) (Vehicle) {
	var vehicle Vehicle

	vehicle.Id = id
	vehicle.Type = "vehicle"

	vehicle.Details.Make = make
	vehicle.Details.Model = model
	vehicle.Details.Registration = registration
	vehicle.Details.Year = year
	vehicle.Details.Mileage = mileage

	return vehicle
}
