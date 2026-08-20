package midtrans

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type Config struct {
	client_key string
	server_key string
}

type Service struct {
	client *coreapi.Client
}

func NewService(cfg Config) *Service {
	var c *coreapi.Client

	switch os.Getenv("ENVIRONMENT") {
	case "production":
		c.New(cfg.server_key, midtrans.Production)
	default:
		c.New(cfg.server_key, midtrans.Sandbox)
	}

	return &Service{ client: c }
}
