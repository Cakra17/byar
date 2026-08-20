package coreapi

import (
	"testing"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
)

func TestBriVa(t *testing.T) {

	tests := []struct {
		name  string
		input *Request
		out   *coreapi.ChargeReq
		err   error
	}{
		{
			name: "Create new BRI VA",
			input: &Request{
				TransactionDetails: TransactionDetails{
					Orderid: "bri-123",
					GrossAmount: 100000,
				},
			},
			out: &coreapi.ChargeReq{
				PaymentType: coreapi.PaymentTypeBankTransfer,
				TransactionDetails: midtrans.TransactionDetails{
					OrderID: "bri-123",
					GrossAmt: 100000,
				},
				BankTransfer: &coreapi.BankTransferDetails{
					Bank: midtrans.BankBri,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			va := NewBriVa(test.input)
			assert.Equal(t, va.PaymentType, test.out.PaymentType)
			assert.Equal(t, va.BankTransfer.Bank, test.out.BankTransfer.Bank)
			assert.Equal(t, va.TransactionDetails.OrderID, test.out.TransactionDetails.OrderID)
			assert.Equal(t, va.TransactionDetails.GrossAmt, test.out.TransactionDetails.GrossAmt)
		})
	}
}
