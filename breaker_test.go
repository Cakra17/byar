package byar

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type bgFake struct {
	name string

	createCalls atomic.Int64
	statusCalls atomic.Int64

	createErr error
	statusErr error
	status    TransactionStatus
}

func (f *bgFake) CreateTransaction(ctx context.Context, req *Request) (*Transaction, error) {
	f.createCalls.Add(1)
	if f.createErr != nil {
		return nil, f.createErr
	}
	return &Transaction{Provider: f.name, OrderId: req.TransactionDetails.Orderid, Status: StatusPending}, nil
}

func (f *bgFake) CheckStatus(ctx context.Context, orderId string) (*Transaction, error) {
	f.statusCalls.Add(1)
	if f.statusErr != nil {
		return nil, f.statusErr
	}
	if f.status == "" {
		return nil, fmt.Errorf("transaction not found")
	}
	return &Transaction{Provider: f.name, OrderId: orderId, Status: f.status}, nil
}

func TestBreakerStateMachine(t *testing.T) {
	b := newBreaker(3, 10*time.Millisecond)

	b.recordFailure()
	b.recordFailure()
	if !b.allow() {
		t.Fatal("2 failures < threshold 3, breaker must stay closed")
	}
	b.recordFailure()
	if b.allow() {
		t.Fatal("breaker must open at threshold")
	}

	time.Sleep(12 * time.Millisecond)
	if !b.allow() {
		t.Fatal("cooldown elapsed, probe must be granted")
	}
	if b.allow() {
		t.Fatal("second caller in half-open must be blocked")
	}

	b.recordSuccess()
	if !b.allow() {
		t.Fatal("probe success must close the breaker")
	}
}

func TestBreakerReopensOnProbeFailure(t *testing.T) {
	b := newBreaker(2, 10*time.Millisecond)

	b.recordFailure()
	b.recordFailure()
	time.Sleep(12 * time.Millisecond)

	if !b.allow() {
		t.Fatal("probe expected after cooldown")
	}
	b.recordFailure()

	if b.allow() {
		t.Fatal("failed probe must reopen the breaker")
	}

	time.Sleep(12 * time.Millisecond)
	if !b.allow() {
		t.Fatal("fresh cooldown must grant a new probe")
	}
}

func TestBreakerSuccessResetsCount(t *testing.T) {
	b := newBreaker(3, time.Hour)

	b.recordFailure()
	b.recordFailure()
	b.recordSuccess()
	b.recordFailure()
	b.recordFailure()
	if !b.allow() {
		t.Fatal("counting must restart after success")
	}
}

func TestBreakerGatewayBlocksWritesNotReads(t *testing.T) {
	inner := &bgFake{name: "a", createErr: ErrTransient, statusErr: errors.New("down")}
	gw := &breakerGateway{Gateway: inner, b: newBreaker(2, time.Hour)}
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		if _, err := gw.CreateTransaction(ctx, req); !errors.Is(err, ErrTransient) {
			t.Fatalf("call %d: err = %v, want ErrTransient", i, err)
		}
	}

	if _, err := gw.CreateTransaction(ctx, req); !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("err = %v, want ErrCircuitOpen", err)
	}
	if got := inner.createCalls.Load(); got != 2 {
		t.Errorf("inner create calls = %d, want 2 (third blocked)", got)
	}

	if _, err := gw.CheckStatus(ctx, "o1"); errors.Is(err, ErrCircuitOpen) {
		t.Error("reads must bypass the open breaker")
	}
	if got := inner.statusCalls.Load(); got != 1 {
		t.Errorf("inner status calls = %d, want 1", got)
	}
}

func TestBreakerGatewayPermanentDoesNotTrip(t *testing.T) {
	inner := &bgFake{name: "a", createErr: ErrPermanent}
	gw := &breakerGateway{Gateway: inner, b: newBreaker(2, time.Hour)}
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		if _, err := gw.CreateTransaction(ctx, req); !errors.Is(err, ErrPermanent) {
			t.Fatalf("call %d: err = %v, want ErrPermanent", i, err)
		}
	}
	if got := inner.createCalls.Load(); got != 5 {
		t.Errorf("declines must not trip the breaker, calls = %d", got)
	}
}

func TestBreakerGatewayNotEnabledDoesNotTrip(t *testing.T) {
	inner := &bgFake{name: "a", createErr: ErrNotEnabled}
	gw := &breakerGateway{Gateway: inner, b: newBreaker(2, time.Hour)}
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		if _, err := gw.CreateTransaction(ctx, req); !errors.Is(err, ErrNotEnabled) {
			t.Fatalf("call %d: err = %v, want ErrNotEnabled", i, err)
		}
	}
	if got := inner.createCalls.Load(); got != 5 {
		t.Errorf("not-enabled must not trip the breaker, calls = %d", got)
	}
}

func TestPaySkipsOpenBreaker(t *testing.T) {
	primary := &bgFake{name: "a", createErr: ErrTransient, statusErr: errors.New("down")}
	fallback := &bgFake{name: "b"}

	r := NewGateway(WithBreakerThreshold(2), WithBreakerCooldown(time.Hour))
	if err := r.Register("a", primary, BcaVa); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("b", fallback, BcaVa); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}

	// Pay #1: create fails + reconcile status fails -> breaker trips mid-pay,
	// failover to b.
	tx, err := r.Pay(ctx, req)
	if err != nil || tx.Provider != "b" {
		t.Fatalf("tx = %v, err = %v, want failover to b", tx, err)
	}

	// Pay #2: a's breaker is open, skipped without touching the inner gateway.
	tx, err = r.Pay(ctx, req)
	if err != nil || tx.Provider != "b" {
		t.Fatalf("tx = %v, err = %v, want b", tx, err)
	}
	if got := primary.createCalls.Load(); got != 1 {
		t.Errorf("primary create calls = %d, want 1 (skipped while open)", got)
	}
	if got := fallback.createCalls.Load(); got != 2 {
		t.Errorf("fallback create calls = %d, want 2", got)
	}
}

func TestPayBreakerRecoversAfterCooldown(t *testing.T) {
	primary := &bgFake{name: "a", createErr: ErrTransient, statusErr: errors.New("down")}
	fallback := &bgFake{name: "b"}

	r := NewGateway(WithBreakerThreshold(2), WithBreakerCooldown(20*time.Millisecond))
	if err := r.Register("a", primary, BcaVa); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("b", fallback, BcaVa); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}

	if _, err := r.Pay(ctx, req); err != nil {
		t.Fatal(err)
	}

	time.Sleep(25 * time.Millisecond)
	primary.createErr = nil
	primary.statusErr = nil

	tx, err := r.Pay(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if tx.Provider != "a" {
		t.Errorf("provider = %q, want probe to succeed on a", tx.Provider)
	}
	if got := primary.createCalls.Load(); got != 2 {
		t.Errorf("primary create calls = %d, want 2 (initial + probe)", got)
	}
}

func TestPayAllBreakersOpenReturnsJoinedError(t *testing.T) {
	primary := &bgFake{name: "a", createErr: ErrTransient, statusErr: errors.New("down")}

	r := NewGateway(WithBreakerThreshold(2), WithBreakerCooldown(time.Hour))
	if err := r.Register("a", primary, BcaVa); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := &Request{PaymentType: BcaVa, TransactionDetails: TransactionDetails{Orderid: "o1"}}

	// Pay #1 fails: single provider, transient error, no failover target.
	if _, err := r.Pay(ctx, req); !errors.Is(err, ErrTransient) {
		t.Fatalf("err = %v, want ErrTransient", err)
	}

	_, err := r.Pay(ctx, req)
	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("err = %v, want joined error containing ErrCircuitOpen", err)
	}
}
