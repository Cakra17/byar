package main

type PaymentStatus int
type BankName string

const (
	PAYMENT_SUCCESS PaymentStatus = iota
	PAYMENT_PENDING
	PAYMENT_FAILED
	PAYMENT_REFUND
	PAYMENT_CANCELED
)

const (
	BANK_BCA BankName = "bca"
	BANK_BNI BankName = "bni"
	BANK_BRI BankName = "bri"
	BANK_CIMB BankName = "cimb"
)

type PaymentGateway interface {
	CreateTransaction(req TransactionReq) (*TransactionRes, error)
	CancelTransaction(req TransactionReq) (*TransactionRes, error)
	RefundTransaction(req TransactionReq) (*TransactionRes, error)
	GetTransactionStatus(tsxId string) (*TransactionRes, error)
}

type TransactionReq struct {
	PaymentType string
	OrderId string
	Amount int64
	BankName BankName
	CustomerDetails map[string]any
}

type TransactionRes struct {
	Id string
	Status PaymentStatus
	Message string
	BankName BankName
	SignatureKey string
	FraudStatus string
	Time string
}