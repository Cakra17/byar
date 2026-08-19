package coreapi

import "github.com/midtrans/midtrans-go/coreapi"

func NewCreditCard(req *Request) *coreapi.ChargeReq {
	return newRequest(req).SetCCPayment(req.CreditCardToken).build()
}
