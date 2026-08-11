package coreapi

import (
	md "github.com/cakra17byar/midtrans"
	"github.com/midtrans/midtrans-go"
	cr "github.com/midtrans/midtrans-go/coreapi"
)

type builder struct {
	req *cr.ChargeReq
}

func newRequest(req *md.Request) *builder {
	build := &builder{
		req: &cr.ChargeReq{
			PaymentType: cr.CoreapiPaymentType(req.PaymentType),
			Items:       &[]midtrans.ItemDetails{},
		},
	}

	return build.
		setPayment(req).
		setItem(req)
}

func (b *builder) setItem(req *md.Request) *builder {
	var items []midtrans.ItemDetails

	for _, item := range req.ItemsDetails {
		items = append(items, midtrans.ItemDetails{
			ID:           item.Id,
			Name:         item.Name,
			Price:        item.Price,
			Qty:          item.Quantity,
			Brand:        item.Brand,
			Category:     item.Category,
			MerchantName: item.MerchantName,
		})
	}
	b.req.Items = &items
	return b
}

func (b *builder) setPayment(req *md.Request) *builder {
	b.req.PaymentType = cr.CoreapiPaymentType(req.PaymentType)
	b.req.TransactionDetails = midtrans.TransactionDetails{
		OrderID: req.TransactionDetails.Orderid,
		GrossAmt: req.TransactionDetails.GrossAmount,
	}
	b.req.BankTransfer = &cr.BankTransferDetails{
		Bank: midtrans.Bank(req.Bankname),
	}

	return b
}

func (b *builder) build() *cr.ChargeReq {
	return b.req
}
