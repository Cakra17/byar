package byar

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

type fakeGateway struct {
	methods []PaymentType
}

func (f *fakeGateway) CreateTransaction(ctx context.Context, req *Request) (*Transaction, error) {
	return &Transaction{}, nil
}

func (f *fakeGateway) CheckStatus(ctx context.Context, orderId string) (*Transaction, error) {
	return &Transaction{}, nil
}

func (f *fakeGateway) SupportMethods() []PaymentType {
	return f.methods
}

type plainGateway struct{}

func (p *plainGateway) CreateTransaction(ctx context.Context, req *Request) (*Transaction, error) {
	return &Transaction{}, nil
}

func (p *plainGateway) CheckStatus(ctx context.Context, orderId string) (*Transaction, error) {
	return &Transaction{}, nil
}

func TestRegisterAutoDerivesChain(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("midtrans", &fakeGateway{methods: []PaymentType{BcaVa, CC}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("xendit", &fakeGateway{methods: []PaymentType{BcaVa, GOPAY}}); err != nil {
		t.Fatal(err)
	}

	chain, err := r.Chain(BcaVa)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 2 {
		t.Fatalf("BcaVa chain length = %d, want 2", len(chain))
	}
	if chain[0] != r.gateways["midtrans"] || chain[1] != r.gateways["xendit"] {
		t.Error("BcaVa chain order should follow registration order")
	}

	chain, err = r.Chain(GOPAY)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 1 || chain[0] != r.gateways["xendit"] {
		t.Error("GOPAY chain should contain only xendit")
	}
}

func TestRegisterExplicitMethods(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("plain", &plainGateway{}, BcaVa, BniVa); err != nil {
		t.Fatal(err)
	}

	chain, err := r.Chain(BniVa)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 1 {
		t.Fatalf("BniVa chain length = %d, want 1", len(chain))
	}
}

func TestRegisterNoMethods(t *testing.T) {
	t.Parallel()

	r := NewGateway()

	if err := r.Register("plain", &plainGateway{}); !errors.Is(err, ErrNoCapabilities) {
		t.Errorf("plain gateway without methods: err = %v, want ErrNoCapabilities", err)
	}
	if err := r.Register("empty", &fakeGateway{}); !errors.Is(err, ErrNoCapabilities) {
		t.Errorf("capability gateway declaring nothing: err = %v, want ErrNoCapabilities", err)
	}
}

func TestRegisterReplacesOldEntries(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("gw", &plainGateway{}, BcaVa); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("gw", &plainGateway{}, BriVa); err != nil {
		t.Fatal(err)
	}

	if _, err := r.Chain(BcaVa); !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("old method still chained: err = %v, want ErrUnsupportedMethod", err)
	}
	chain, err := r.Chain(BriVa)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 1 {
		t.Fatalf("BriVa chain length = %d, want 1", len(chain))
	}
}

func TestRegisterDuplicateMethodsDeduped(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("gw", &plainGateway{}, BcaVa, BcaVa); err != nil {
		t.Fatal(err)
	}

	chain, err := r.Chain(BcaVa)
	if err != nil {
		t.Fatal(err)
	}
	if len(chain) != 1 {
		t.Fatalf("BcaVa chain length = %d, want 1 (deduped)", len(chain))
	}
}

func TestSetPriorityReorders(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("midtrans", &fakeGateway{methods: []PaymentType{BcaVa}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("xendit", &fakeGateway{methods: []PaymentType{BcaVa}}); err != nil {
		t.Fatal(err)
	}

	if err := r.SetPriority(BcaVa, "xendit", "midtrans"); err != nil {
		t.Fatal(err)
	}

	chain, err := r.Chain(BcaVa)
	if err != nil {
		t.Fatal(err)
	}
	if chain[0] != r.gateways["xendit"] || chain[1] != r.gateways["midtrans"] {
		t.Error("SetPriority did not reorder the chain")
	}
}

func TestSetPriorityValidation(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("midtrans", &fakeGateway{methods: []PaymentType{BcaVa}}); err != nil {
		t.Fatal(err)
	}

	if err := r.SetPriority(BcaVa, "unknown"); !errors.Is(err, ErrProviderNotExist) {
		t.Errorf("unknown provider: err = %v, want ErrProviderNotExist", err)
	}
	if err := r.SetPriority(GOPAY, "midtrans"); !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("method not declared: err = %v, want ErrUnsupportedMethod", err)
	}
	if err := r.SetPriority(BcaVa); !errors.Is(err, ErrEmptyChain) {
		t.Errorf("empty chain: err = %v, want ErrEmptyChain", err)
	}
	if err := r.SetPriority(BcaVa, "midtrans", "midtrans"); err == nil {
		t.Error("duplicate provider should be rejected")
	}
}

func TestChainUnsupportedMethod(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if _, err := r.Chain(CC); !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("err = %v, want ErrUnsupportedMethod", err)
	}
}

func TestGetUnknownProvider(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if _, err := r.Get("nope"); !errors.Is(err, ErrUnsupportedProvider) {
		t.Errorf("err = %v, want ErrUnsupportedProvider", err)
	}
}

func TestRegistryConcurrent(t *testing.T) {
	t.Parallel()

	r := NewGateway()
	if err := r.Register("midtrans", &fakeGateway{methods: []PaymentType{BcaVa}}); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := range 8 {
		key := fmt.Sprintf("gw%d", i)
		wg.Add(3)
		go func() {
			defer wg.Done()
			_ = r.Register(key, &fakeGateway{methods: []PaymentType{BcaVa}})
		}()
		go func() {
			defer wg.Done()
			_, _ = r.Chain(BcaVa)
		}()
		go func() {
			defer wg.Done()
			_, _ = r.Get(key)
		}()
	}
	wg.Wait()
}
