package byar

import (
	"context"
	"fmt"
	"sync"
)

type Gateway interface {
	CreateTransaction(ctx context.Context, req *Request) (*Transaction, error)
	CheckStatus(ctx context.Context, orderId string) (*Transaction, error)
}

type CapabilityProvider interface {
	SupportMethods() []PaymentType
}

type GatewayRegistry struct {
	mu       sync.RWMutex
	gateways map[string]Gateway
	caps     map[string]map[PaymentType]struct{}
	matrix   map[PaymentType][]string
}

func NewGateway() *GatewayRegistry {
	return &GatewayRegistry{
		gateways: make(map[string]Gateway),
		caps:     make(map[string]map[PaymentType]struct{}),
		matrix:   make(map[PaymentType][]string),
	}
}

func (r *GatewayRegistry) Get(provider string) (Gateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gw, ok := r.gateways[provider]
	if !ok {
		return nil, ErrUnsupportedProvider
	}
	return gw, nil
}

func (r *GatewayRegistry) Register(key string, gw Gateway, methods ...PaymentType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	supported := methods
	if cp, ok := gw.(CapabilityProvider); ok {
		supported = cp.SupportMethods()
	}
	if len(supported) == 0 {
		return ErrNoCapabilities
	}

	r.remove(key)

	caps := make(map[PaymentType]struct{}, len(supported))
	for _, pt := range supported {
		if _, dup := caps[pt]; dup {
			continue
		}
		caps[pt] = struct{}{}
		r.matrix[pt] = append(r.matrix[pt], key)
	}

	r.gateways[key] = gw
	r.caps[key] = caps
	return nil
}

func (r *GatewayRegistry) SetPriority(method PaymentType, providers ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(providers) == 0 {
		return ErrEmptyChain
	}

	seen := make(map[string]struct{}, len(providers))
	for _, key := range providers {
		if _, dup := seen[key]; dup {
			return fmt.Errorf("duplicate provider %q in chain", key)
		}
		if _, ok := r.gateways[key]; !ok {
			return fmt.Errorf("provider %q: %w", key, ErrProviderNotExist)
		}
		if _, ok := r.caps[key][method]; !ok {
			return fmt.Errorf("provider %q: %w", key, ErrUnsupportedMethod)
		}
		seen[key] = struct{}{}
	}

	chain := make([]string, len(providers))
	copy(chain, providers)
	r.matrix[method] = chain
	return nil
}

func (r *GatewayRegistry) Chain(method PaymentType) ([]Gateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := r.matrix[method]
	if len(keys) == 0 {
		return nil, ErrUnsupportedMethod
	}

	chain := make([]Gateway, 0, len(keys))
	for _, key := range keys {
		chain = append(chain, r.gateways[key])
	}
	return chain, nil
}

func (r *GatewayRegistry) remove(key string) {
	for pt := range r.caps[key] {
		r.matrix[pt] = removeKey(r.matrix[pt], key)
		if len(r.matrix[pt]) == 0 {
			delete(r.matrix, pt)
		}
	}
	delete(r.gateways, key)
	delete(r.caps, key)
}

func removeKey(chain []string, key string) []string {
	out := make([]string, 0, len(chain)-1) //chain[:0]
	for _, k := range chain {
		if k != key {
			out = append(out, k)
		}
	}
	return out
}
