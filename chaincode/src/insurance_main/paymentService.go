package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Stored the chaincodeId of the CRUD chaincode
const SETTLEMENTS_CHAINCODE_ID_KEY = "SETTLEMENTS_CHAINCODE_ID"

const ADD_PAYMENT_FUNCTION = "addPayment"

type PaymentService struct {}

func InitPaymentService(stub shim.ChaincodeStubInterface, settlementsChaincodeId string) {
	stub.PutState(SETTLEMENTS_CHAINCODE_ID_KEY, []byte(settlementsChaincodeId))
}

func AddInsurerPayment(stub shim.ChaincodeStubInterface, payment ClaimDetailsSettlementPayment, linkedClaimId string) (error) {

	chaincodeId := getSettlementsChaincodeId(stub)

	if len(chaincodeId) > 0 {
		args := getPaymentArgs(payment, linkedClaimId)

		err := invoke(stub, getSettlementsChaincodeId(stub), ADD_PAYMENT_FUNCTION, args)

		return err
	}

	return nil
}

func getPaymentArgs(payment ClaimDetailsSettlementPayment, linkedClaimId string) ([]string) {
	return []string{"toInsurerForClaim", payment.Sender, payment.Recipient, strconv.Itoa(payment.Amount), linkedClaimId}
}

func getSettlementsChaincodeId(stub shim.ChaincodeStubInterface) (string) {
	bytes, _ := stub.GetState(SETTLEMENTS_CHAINCODE_ID_KEY)
	return string(bytes)
}
