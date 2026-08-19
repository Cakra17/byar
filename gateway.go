package byar

import "fmt"

type Gateway interface{
	CreateTransaction()
}

type GatewayRegistry struct {
	gateways map[string]*Gateway
}

func (r *GatewayRegistry) Get(provider string) (*Gateway, error) {
	gw, ok := r.gateways[provider]
	if !ok {
		return nil, fmt.Errorf("Unsuported Payment Provider")
	}
	return gw, nil
}

func (r *GatewayRegistry) Register(provider string, pg *Gateway) {
	r.gateways[provider] = pg
}
