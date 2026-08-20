package byar

type PaymentType string

const (
	BcaVa  PaymentType = "bca_va"
	BriVa  PaymentType = "bri_va"
	BniVa  PaymentType = "bni_va"
	CimbVa PaymentType = "cimb_va"
	CC     PaymentType = "credit_card"
	GOPAY  PaymentType = "gopay"
)
