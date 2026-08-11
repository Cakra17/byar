package coreapi

import (
	"github.com/cakra17byar/midtrans"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewBriVa(req *midtrans.Request) *coreapi.ChargeReq {
	return newRequest(req).build()
}