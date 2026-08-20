package coreapi

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBriVa(req *Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankBri).
		build()
}
