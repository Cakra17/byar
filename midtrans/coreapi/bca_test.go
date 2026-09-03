package coreapi

import (
	"testing"

	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
)

func TestBcaVa(t *testing.T) {
	tests := []struct {
		name  string
		input *byar.Request
		out   *coreapi.ChargeReq
		err   error
	}{
		{
			name: "Create new BCA VA",
			input: &byar.Request{
				TransactionDetails: byar.TransactionDetails{
					Orderid: "bca-123",
					GrossAmount: 100000,
				},
			},
			out: &coreapi.ChargeReq{
				PaymentType: coreapi.PaymentTypeBankTransfer,
				TransactionDetails: midtrans.TransactionDetails{
					OrderID: "bca-123",
					GrossAmt: 100000,
				},
				BankTransfer: &coreapi.BankTransferDetails{
					Bank: midtrans.BankBca,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			va := NewBcaVa(test.input)
			assert.Equal(t, va.PaymentType, test.out.PaymentType)
			assert.Equal(t, va.BankTransfer.Bank, test.out.BankTransfer.Bank)
			assert.Equal(t, va.TransactionDetails.OrderID, test.out.TransactionDetails.OrderID)
			assert.Equal(t, va.TransactionDetails.GrossAmt, test.out.TransactionDetails.GrossAmt)
		})
	}
}
