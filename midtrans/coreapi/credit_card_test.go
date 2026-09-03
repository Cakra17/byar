package coreapi

import (
	"testing"

	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
)

func TestCC(t *testing.T) {

	tests := []struct {
		name  string
		input *byar.Request
		out   *coreapi.ChargeReq
		err   error
	}{
		{
			name: "Create new Credit Card",
			input: &byar.Request{
				TransactionDetails: byar.TransactionDetails{
					Orderid:     "cc-123",
					GrossAmount: 100000,
				},
				CreditCardToken: "cc-lksdawjdlkdwa",
			},
			out: &coreapi.ChargeReq{
				PaymentType: coreapi.PaymentTypeCreditCard,
				TransactionDetails: midtrans.TransactionDetails{
					OrderID:  "cc-123",
					GrossAmt: 100000,
				},
				CreditCard: &coreapi.CreditCardDetails{
					TokenID: "cc-lksdawjdlkdwa",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			va := NewCreditCard(test.input)
			assert.Equal(t, va.PaymentType, test.out.PaymentType)
			assert.Equal(t, va.CreditCard.TokenID, test.out.CreditCard.TokenID)
			assert.Equal(t, va.TransactionDetails.OrderID, test.out.TransactionDetails.OrderID)
			assert.Equal(t, va.TransactionDetails.GrossAmt, test.out.TransactionDetails.GrossAmt)
		})
	}
}
