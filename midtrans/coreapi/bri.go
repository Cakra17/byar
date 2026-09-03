package coreapi

import (
	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBriVa(req *byar.Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankBri).
		build()
}
