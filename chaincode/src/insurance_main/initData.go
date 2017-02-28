package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func InitReferenceData(stub shim.ChaincodeStubInterface) (error) {
	_, err := SavePolicy(stub, NewPolicy("P1", "claimant1", "insurer1", "31/01/17", "30/01/18", 300, "BP08BRV"))
	if err != nil{ return err }
	_, err = SavePolicy(stub, NewPolicy("P2", "claimant2", "insurer2", "11/01/17", "10/01/18", 250, "DZ14TYV"))
	if err != nil{ return err }

	_, err = SaveVehicle(stub, NewVehicle("BP08BRV", "Ford", "Focus", "2008", 55000, "100924404"))
	if err != nil{ return err }
	_, err = SaveVehicle(stub, NewVehicle("DZ14TYV", "Audi", "TT", "2014", 20000, "200481078"))
	if err != nil{ return err }

	_, err = SaveUser(stub, NewUser("U1", "John", "Hancock", "john.hancock@outlook.com", "P1"))
	if err != nil{ return err }
	_, err = SaveUser(stub, NewUser("U2", "Jane", "Doe", "jane.doe@outlook.com", "P2"))
	if err != nil{ return err }

	var approvedGarages ApprovedGarages
	approvedGarages.Garages = []string{"garage1", "garage2", "garage3"}
	_, err = SaveApprovedGarages(stub, approvedGarages)

	return err
}
