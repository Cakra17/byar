package midtrans

import (
	"os"

	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type Service struct {
	client *coreapi.Client
}

func NewService(cfg byar.Config) *Service {
	env := midtrans.Sandbox
	if os.Getenv("ENVIRONMENT") == "production" {
		env = midtrans.Production
	}

	var c coreapi.Client
	c.New(cfg.ServerKey, env)

	return &Service{client: &c}
}

func (s Service) SupportMethods() []byar.PaymentType {
	return []byar.PaymentType{
		byar.BcaVa,
		byar.BniVa,
		byar.BriVa,
		byar.CimbVa,
		byar.CC,
		byar.GOPAY,
	}
}
