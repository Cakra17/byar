package coreapi

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBniVa(req *Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankBni).
		build()
}
