package midtrans

import (
	"github.com/cakra17/byar"
	cp "github.com/cakra17/byar/midtrans/coreapi"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewPayment(req *cp.Request) *coreapi.ChargeReq {
	switch req.PaymentType {
	case byar.BcaVa:
		return cp.NewBcaVa(req)
	case byar.BniVa:
		return cp.NewBniVa(req)
	case byar.BriVa:
		return cp.NewBriVa(req)
	case byar.CimbVa:
		return cp.NewCimbVa(req)
	case byar.CC:
		return cp.NewCreditCard(req)
	case byar.GOPAY:
		return cp.NewGopay(req)
	default:
		return nil
	}
}
