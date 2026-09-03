package midtrans

import (
	"time"

	"github.com/cakra17/byar"
	cp "github.com/cakra17/byar/midtrans/coreapi"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewPayment(req *byar.Request) *coreapi.ChargeReq {
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

func NewResponse(provider string, tp byar.PaymentType, res *coreapi.ChargeResponse) *byar.Transaction {
	var tx byar.Transaction	
	switch tp {
	case byar.BcaVa, byar.BniVa, byar.BriVa, byar.CimbVa:
		tx = byar.Transaction{
			Provider: provider,
			Id:       res.ID,
			VaNumber: res.VaNumbers[0].VANumber,
			OrderId:  res.OrderID,
			Status:   byar.TransactionStatus(res.TransactionStatus),
			Amount:   res.GrossAmount,
			Message:  res.StatusMessage,
			Time:     res.TransactionTime,
			Currency: res.Currency,
			Expiry:   parseTime(res.ExpiryTime),
		}
		case byar.CC:
		tx = byar.Transaction{
			Provider: provider,
			Id:       res.ID,
			OrderId:  res.OrderID,
			Status:   byar.TransactionStatus(res.TransactionStatus),
			Amount:   res.GrossAmount,
			Message:  res.StatusMessage,
			Time:     res.TransactionTime,
			Currency: res.Currency,
			Expiry:   parseTime(res.ExpiryTime),
		}
	case byar.GOPAY:
		tx = byar.Transaction{
			Provider: provider,
			Id:       res.ID,
			OrderId:  res.OrderID,
			Status:   byar.TransactionStatus(res.TransactionStatus),
			Amount:   res.GrossAmount,
			Message:  res.StatusMessage,
			Time:     res.TransactionTime,
			Currency: res.Currency,
			Expiry:   parseTime(res.ExpiryTime),
		}
	}
	return &tx
}

func parseTime(tm string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", tm)
	if err != nil {
		return time.Now().Add(time.Duration(t.Day()))
	}
	return t
}