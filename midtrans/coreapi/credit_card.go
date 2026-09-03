package coreapi

import (
	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewCreditCard(req *byar.Request) *coreapi.ChargeReq {
	return newRequest(req).SetCCPayment(req.CreditCardToken).build()
}
