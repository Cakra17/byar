package coreapi

import (
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewGopay(req *Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetGopayPayment(req.Callback).
		build()
}
