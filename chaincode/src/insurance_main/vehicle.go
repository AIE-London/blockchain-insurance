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
	StyleId			string		`json:"styleId"`
}

//=================================================================================================================================
//	 New Vehicle	-	Constructs a new vehicle
//=================================================================================================================================
func NewVehicle(registration string, make string, model string, year string, mileage int, styleId string) (Vehicle) {
	var vehicle Vehicle

	vehicle.Id = registration
	vehicle.Type = "vehicle"

	vehicle.Details.Make = make
	vehicle.Details.Model = model
	vehicle.Details.Registration = registration
	vehicle.Details.Year = year
	vehicle.Details.Mileage = mileage
	vehicle.Details.StyleId = styleId

	return vehicle
}
