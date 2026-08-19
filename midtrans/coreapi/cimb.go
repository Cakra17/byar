package coreapi

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewCimbVa(req *Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(midtrans.BankCimb).
		build()
}
