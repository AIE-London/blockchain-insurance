package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func InitReferenceData(stub shim.ChaincodeStubInterface) (error) {
	SavePolicy(stub, NewPolicy("P1", "claimant1", "31/01/17", "30/01/18", 300, "BP08BRV"))
	SavePolicy(stub, NewPolicy("P2", "claimant2", "11/01/17", "10/01/18", 250, "DZ14TYV"))

	SaveVehicle(stub, NewVehicle("BP08BRV", "Ford", "Focus", "2008", 55000, "100924404"))
	SaveVehicle(stub, NewVehicle("DZ14TYV", "Audi", "TT", "2014", 20000, "200481078"))

	return nil
}

//User1 data
var UserArgs1 = []string{"John", "Hancock", "john.hancock@outlook.com"}
var UserCaller1 = "claimant1"
var UserCallerAffiliation1 = "group1"

//User2 data
var UserArgs2 = []string{"Jane", "Doe", "jane.doe@outlook.com"}
var UserCaller2 = "claimant2"
var UserCallerAffiliation2 = "group1"

//ApprovedGarage data
var ApprovedGaragesData = []string{"garage1", "garage2", "garage3"}
