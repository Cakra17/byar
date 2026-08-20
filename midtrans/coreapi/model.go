package coreapi

import (
	"time"

	"github.com/cakra17/byar"
)

type Request struct {
	PaymentType byar.PaymentType
	// Bankname           Bank
	CreditCardToken    string
	TransactionDetails TransactionDetails
	ItemsDetails       []ItemsDetails
	CustomerDetails    *CustomerDetails
	SellerDetails      *SellerDetails
	CustomExpiry       *CustomExpiry
	Callback           string
}

type TransactionDetails struct {
	Orderid     string
	GrossAmount int64
}

type ItemsDetails struct {
	Id           string
	Price        int64
	Quantity     int32
	Name         string
	Brand        string
	Category     string
	MerchantName string
	Url          string
}

type CustomerDetails struct {
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

type SellerDetails struct {
	Id    string
	Name  string
	Email string
	Url   string
}

type CustomExpiry struct {
	OrderTime      time.Time
	ExpiryDuration int
	Unit           string // seconds, minutes, hours
}
