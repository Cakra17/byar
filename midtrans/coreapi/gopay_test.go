package coreapi

import (
	"testing"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
)

func TestGopay(t *testing.T) {

	tests := []struct {
		name  string
		input *Request
		out   *coreapi.ChargeReq
		err   error
	}{
		{
			name: "Create new Gopay",
			input: &Request{
				TransactionDetails: TransactionDetails{
					Orderid:     "gopay-123",
					GrossAmount: 100000,
				},
				Callback: "someapps://callback",
			},
			out: &coreapi.ChargeReq{
				PaymentType: coreapi.PaymentTypeGopay,
				TransactionDetails: midtrans.TransactionDetails{
					OrderID:  "gopay-123",
					GrossAmt: 100000,
				},
				Gopay: &coreapi.GopayDetails{
					EnableCallback: true,
					CallbackUrl:    "someapps://callback",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			va := NewGopay(test.input)
			assert.Equal(t, va.PaymentType, test.out.PaymentType)
			assert.Equal(t, va.Gopay.EnableCallback, test.out.Gopay.EnableCallback)
			assert.Equal(t, va.Gopay.CallbackUrl, test.out.Gopay.CallbackUrl)
			assert.Equal(t, va.TransactionDetails.OrderID, test.out.TransactionDetails.OrderID)
			assert.Equal(t, va.TransactionDetails.GrossAmt, test.out.TransactionDetails.GrossAmt)
		})
	}
}
