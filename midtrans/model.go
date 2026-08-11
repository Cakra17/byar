package midtrans

import "time"

type Bank string
type Payment string

const (
	BANKBCA  Bank = "bca"
	BANKBNI  Bank = "bni"
	BANKBRI  Bank = "bri"
	BANKCIMB Bank = "cimb"
)

const (
	BANK_TRANSFER Payment = "bank_transfer"
	CREDIT_CARD   Payment = "credit_card"
	QRIS          Payment = "qris"
)

type Request struct {
	PaymentType        Payment
	Bankname           Bank
	TransactionDetails TransactionDetails
	ItemsDetails       []ItemsDetails
	CustomerDetails    *CustomerDetails
	SellerDetails      *SellerDetails
	CustomExpiry       *CustomExpiry
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
