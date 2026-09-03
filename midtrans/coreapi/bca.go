package coreapi

import (
	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBcaVa(req *byar.Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankBca).
		build()
}
