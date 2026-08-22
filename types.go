package byar

import "time"

type Config struct {
	ClientKey string
	ServerKey string
}

type PaymentType string

const (
	BcaVa  PaymentType = "bca_va"
	BriVa  PaymentType = "bri_va"
	BniVa  PaymentType = "bni_va"
	CimbVa PaymentType = "cimb_va"
	CC     PaymentType = "credit_card"
	GOPAY  PaymentType = "gopay"
)

type Request struct {
	Provider           string
	PaymentType        PaymentType
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

type Transaction struct {
	Id         string
	OrderId    string
	VaNumber   string
	PaymentUrl string
	Status     string
	Expiry     time.Time
	Service    string
}
