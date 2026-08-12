package coreapi

import (
	"github.com/midtrans/midtrans-go"
	md "github.com/cakra17/byar/midtrans"
	cr "github.com/midtrans/midtrans-go/coreapi"
)

type builder struct {
	req *cr.ChargeReq
}

func newRequest(req *md.Request) *builder {
	build := &builder{
		req: &cr.ChargeReq{
			Items:       &[]midtrans.ItemDetails{},
		},
	}

	return build.
		setCustomer(req).
		setTransaction(req).
		setItem(req).
		setExpiry(req)
}

func (b *builder) setItem(req *md.Request) *builder {
	if (len(req.ItemsDetails) > 0) {	
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
	}
	return b
}

func (b *builder) setTransaction(req *md.Request) *builder {
	b.req.TransactionDetails = midtrans.TransactionDetails{
		OrderID: req.TransactionDetails.Orderid,
		GrossAmt: req.TransactionDetails.GrossAmount,
	}	

	return b
}

func (b *builder) SetBankPayment(typ cr.CoreapiPaymentType, bank midtrans.Bank) *builder {
	b.req.BankTransfer = &cr.BankTransferDetails{
		Bank: bank,
	}
	b.req.PaymentType = typ

	return b
}

func (b *builder) setCustomer(req *md.Request) *builder {
	if req.CustomerDetails != nil {
		b.req.CustomerDetails = &midtrans.CustomerDetails{
			FName: req.CustomerDetails.FirstName,
			LName: req.CustomerDetails.LastName,
			Email: req.CustomerDetails.Email,
			Phone: req.CustomerDetails.Phone,
		}
	}
	return b
}

func (b *builder) setExpiry(req *md.Request) *builder {
	if req.CustomExpiry != nil {
		b.req.CustomExpiry = &cr.CustomExpiry{
			OrderTime: req.CustomExpiry.OrderTime.String(),
			ExpiryDuration: req.CustomExpiry.ExpiryDuration,
			Unit: req.CustomExpiry.Unit,
		}
	}
	return b
}

func (b *builder) build() *cr.ChargeReq {
	return b.req
}
