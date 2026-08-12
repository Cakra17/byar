package coreapi

import (
	md "github.com/cakra17/byar/midtrans"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewCimbVa(req *md.Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(coreapi.PaymentTypeBankTransfer, midtrans.BankCimb).
		build()
}
