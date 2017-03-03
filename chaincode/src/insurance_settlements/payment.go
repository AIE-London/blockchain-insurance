package main

//==============================================================================================================================
//	Policy - Defines the structure for a policy object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Payment struct {
	Id					string				`json:"id"`
	Status				string				`json:"status"`
	Type				string				`json:"type"`
	From				string				`json:"from"`
	To					string				`json:"to"`
	Amount	 			int					`json:"amount"`
	LinkedId			string				`json:"linkedId"`
	InvoiceId			string				`json:"invoiceId"`
}

const PAYMENT_STATUS_PENDING = "pending"

const PAYMENT_TYPE_TO_INSURER_FOR_CLAIM = "toInsurerForClaim"

//=================================================================================================================================
//	 newPayment	-	Constructs a new payment
//=================================================================================================================================
func NewPayment(id string, paymentType string, from string, to string, amount int, linkedId string) (Payment) {
	var payment Payment

	payment.Id = id
	payment.Status = PAYMENT_STATUS_PENDING
	payment.Type = paymentType
	payment.From = from
	payment.To = to
	payment.Amount = amount
	payment.LinkedId = linkedId

	return payment
}
