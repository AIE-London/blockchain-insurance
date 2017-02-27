package main

import (
	"fmt"
	"strconv"
	"net/http"
	"log"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
)

//Stores the oracle host address
const ORACLE_HOST_KEY = "oracleHost"

func InitOracleService(stub shim.ChaincodeStubInterface, args []string) {
	if (len(args) >= 2) {
		stub.PutState(ORACLE_HOST_KEY, []byte(args[1]))
	} else {
		stub.DelState(ORACLE_HOST_KEY)
	}
}

func RequestVehicleValuationFromOracle(stub shim.ChaincodeStubInterface, txId string, vehicle Vehicle, callbackFunction string) (error){
	bytes, err := stub.GetState(ORACLE_HOST_KEY)

	oracleHost := string(bytes)

	if (err != nil || len(oracleHost) == 0) { log.Println("Unable to obain oracle host: ", err); return errors.New("Unable to obain oracle host") }

	fmt.Println("running requestVehicleValuationFromOracle()")

	url := fmt.Sprintf("http://" + oracleHost +
			"/component/oracle/vehicle/%[1]s/value?mileage=%[2]s&requestId=%[3]s&callbackFunction=%[4]s",
			vehicle.Details.StyleId, strconv.Itoa(vehicle.Details.Mileage), txId, callbackFunction);

	err = sendHttpRequest(url)

	if (err != nil) {
		log.Println("Unable to request vehicle valuation: ", err)
		log.Println("Not erroring though as theres no guarantee that all other peers failed. Hopefully one got through!", err)
	}

	return nil
}

func sendHttpRequest(url string) (error){
	fmt.Println("url: " + url);
	// Build the request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil { log.Println("NewRequest: ", err); return errors.New("Error contacting oracle") }

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil { log.Println("Do: ", err); return errors.New("Error contacting oracle") }

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	return nil
}
