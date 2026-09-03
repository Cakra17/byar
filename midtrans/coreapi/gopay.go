package coreapi

import (
	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewGopay(req *byar.Request) *coreapi.ChargeReq {
	return newRequest(req).
		SetGopayPayment(req.Callback).
		build()
}
