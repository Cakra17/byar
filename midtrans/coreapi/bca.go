package coreapi

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBcaVa(req *Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankBca).
		build()
}
