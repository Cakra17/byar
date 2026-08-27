package byar

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	defaultBreakerThreshold = 5
	defaultBreakerCooldown  = 30 * time.Second
)

type breakerState int

const (
	stateClosed breakerState = iota
	stateOpen
	stateHalfOpen
)

type breaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration

	state    breakerState
	failures int
	openedAt time.Time
}

func newBreaker(threshold int, cooldown time.Duration) *breaker {
	return &breaker{
		threshold: threshold,
		cooldown:  cooldown,
		state:     stateClosed,
	}
}

// allow reports whether a state-changing call may proceed. When the cooldown
// after tripping has elapsed it transitions to half-open and grants exactly
// one probe slot; concurrent callers are denied until the probe resolves.
func (b *breaker) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case stateClosed:
		return true
	case stateOpen:
		if time.Since(b.openedAt) < b.cooldown {
			return false
		}
		b.state = stateHalfOpen
		return true
	default: // stateHalfOpen: probe already in flight
		return false
	}
}

func (b *breaker) recordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.state = stateClosed
	b.failures = 0
}

func (b *breaker) recordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case stateClosed:
		b.failures++
		if b.failures >= b.threshold {
			b.trip()
		}
	case stateHalfOpen:
		b.trip()
	}
}

func (b *breaker) trip() {
	b.state = stateOpen
	b.failures = 0
	b.openedAt = time.Now()
}

// breakerGateway wraps a Gateway with a circuit breaker.
//
// Write calls (CreateTransaction) are guarded by the breaker. Read calls
// (CheckStatus) bypass it: status checks are the tool for resolving
// ambiguous failures, so blocking them would risk double charges. Their
// outcome still feeds the breaker — an answered status check proves the
// provider is alive and closes an open circuit.
type breakerGateway struct {
	Gateway
	b *breaker
}

func (g *breakerGateway) CreateTransaction(ctx context.Context, req *Request) (*Transaction, error) {
	if !g.b.allow() {
		return nil, fmt.Errorf("%w: request not sent", ErrCircuitOpen)
	}

	tx, err := g.Gateway.CreateTransaction(ctx, req)
	g.record(err)
	return tx, err
}

func (g *breakerGateway) CheckStatus(ctx context.Context, orderId string) (*Transaction, error) {
	tx, err := g.Gateway.CheckStatus(ctx, orderId)
	g.record(err)
	return tx, err
}

func (g *breakerGateway) record(err error) {
	switch {
	case err == nil:
		g.b.recordSuccess()
	case errors.Is(err, ErrPermanent), errors.Is(err, ErrNotEnabled):
		// provider answered; declines and config errors say nothing about health
	default:
		g.b.recordFailure()
	}
}
