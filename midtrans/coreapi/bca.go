package coreapi

import (
	md "github.com/cakra17/byar/midtrans"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBcaVa(req *md.Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetBankPayment(coreapi.PaymentTypeBankTransfer, midtrans.BankBca).
		build()
}
