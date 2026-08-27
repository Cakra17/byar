package byar

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

type payGateway struct {
	name string

	mu        atomic.Int64
	createErr error // returned by CreateTransaction
	statusErr error // returned by CheckStatus
	status    TransactionStatus
}

func (g *payGateway) CreateTransaction(ctx context.Context, req *Request) (*Transaction, error) {
	g.mu.Add(1)
	if g.createErr != nil {
		return nil, g.createErr
	}
	return &Transaction{Provider: g.name, OrderId: req.TransactionDetails.Orderid, Status: StatusPending}, nil
}

func (g *payGateway) CheckStatus(ctx context.Context, orderId string) (*Transaction, error) {
	if g.statusErr != nil {
		return nil, g.statusErr
	}
	if g.status == "" {
		return nil, fmt.Errorf("transaction not found")
	}
	return &Transaction{Provider: g.name, OrderId: orderId, Status: g.status}, nil
}

func (g *payGateway) calls() int64 { return g.mu.Load() }

func newTestRegistry(t *testing.T, gws ...*payGateway) *GatewayRegistry {
	t.Helper()
	r := NewGateway()
	for _, g := range gws {
		if err := r.Register(g.name, g, BcaVa); err != nil {
			t.Fatal(err)
		}
	}
	return r
}

func testReq() *Request {
	return &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "order-1"}}
}

func TestPayFirstProviderSucceeds(t *testing.T) {
	primary := &payGateway{name: "a"}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "a" {
		t.Errorf("provider = %q, want %q", tx.Provider, "a")
	}
	if fallback.calls() != 0 {
		t.Error("fallback should not be called when primary succeeds")
	}
}

func TestPayTransientReconcileFound(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrTransient, status: StatusPending}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "a" || tx.Status != StatusPending {
		t.Errorf("reconciled tx = %+v, want pending tx from a", tx)
	}
	if fallback.calls() != 0 {
		t.Error("must not fail over when reconcile finds the transaction")
	}
}

func TestPayTransientReconcileNotFound(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrTransient}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "b" {
		t.Errorf("provider = %q, want failover to b", tx.Provider)
	}
}

func TestPayTransientReconcileStatusErr(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrTransient, statusErr: errors.New("down")}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "b" {
		t.Errorf("provider = %q, want failover when status check fails", tx.Provider)
	}
}

func TestPayReconcileTerminalStatusFailsOver(t *testing.T) {
	for _, s := range []TransactionStatus{StatusExpire, StatusCancel, StatusDeny} {
		primary := &payGateway{name: "a", createErr: ErrTransient, status: s}
		fallback := &payGateway{name: "b"}
		r := newTestRegistry(t, primary, fallback)

		tx, err := r.Pay(context.Background(), testReq())
		if err != nil {
			t.Fatalf("status %s: %v", s, err)
		}
		if tx.Provider != "b" {
			t.Errorf("status %s: provider = %q, want failover to b", s, tx.Provider)
		}
	}
}

func TestPayPermanentAbortsChain(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrPermanent}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	_, err := r.Pay(context.Background(), testReq())
	if !errors.Is(err, ErrPermanent) {
		t.Errorf("err = %v, want ErrPermanent", err)
	}
	if fallback.calls() != 0 {
		t.Error("must not fail over on permanent error")
	}
}

func TestPayNotEnabledSkipsProvider(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrNotEnabled}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "b" {
		t.Errorf("provider = %q, want skip to b", tx.Provider)
	}
}

func TestPayNotEnabledDoesNotReconcile(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrNotEnabled, status: StatusPending}
	r := newTestRegistry(t, primary)

	_, err := r.Pay(context.Background(), testReq())
	if !errors.Is(err, ErrNotEnabled) {
		t.Errorf("err = %v, want ErrNotEnabled", err)
	}
}

func TestPayAllFailJoinsErrors(t *testing.T) {
	a := &payGateway{name: "a", createErr: fmt.Errorf("%w: 503", ErrTransient)}
	b := &payGateway{name: "b", createErr: fmt.Errorf("%w: 503", ErrTransient)}
	r := newTestRegistry(t, a, b)

	_, err := r.Pay(context.Background(), testReq())
	if !errors.Is(err, ErrTransient) {
		t.Errorf("err = %v, want joined ErrTransient", err)
	}
	if a.calls() != 1 || b.calls() != 1 {
		t.Errorf("calls = a:%d b:%d, want each tried once", a.calls(), b.calls())
	}
}

func TestPayUnclassifiedErrorFailsOver(t *testing.T) {
	primary := &payGateway{name: "a", createErr: errors.New("some raw SDK error")}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	tx, err := r.Pay(context.Background(), testReq())
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "b" {
		t.Errorf("provider = %q, want failover on unclassified error", tx.Provider)
	}
}

func TestPayUnsupportedMethod(t *testing.T) {
	r := newTestRegistry(t, &payGateway{name: "a"})

	req := &Request{PaymentType: GOPAY, TransactionDetails: TransactionDetails{Orderid: "order-1"}}
	if _, err := r.Pay(context.Background(), req); !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("err = %v, want ErrUnsupportedMethod", err)
	}
}

func TestPayContextCancelled(t *testing.T) {
	primary := &payGateway{name: "a", createErr: ErrTransient}
	fallback := &payGateway{name: "b"}
	r := newTestRegistry(t, primary, fallback)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := r.Pay(ctx, testReq()); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}
