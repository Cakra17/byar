package byar

import (
	"context"
	"errors"
	"time"
)

// reconcileTimeout bounds the status check after a transient failure so a
// dead provider cannot stall the failover loop.
const reconcileTimeout = 5 * time.Second

// Pay executes req against the provider chain for its payment method.
//
// Failover rules:
//   - success: transaction returned immediately
//   - ErrTransient: the same gateway is asked for the transaction status
//     (the charge may have succeeded despite the error). If the transaction
//     exists and is not expired/cancelled/denied it is returned as-is,
//     preventing double charges; otherwise the next provider is tried
//   - ErrNotEnabled: provider skipped, next provider tried
//   - ErrPermanent: chain aborted, error returned
//
// If every provider fails, the joined error of all attempts is returned.
func (r *GatewayRegistry) Pay(ctx context.Context, req *Request) (*Transaction, error) {
	chain, err := r.Chain(req.PaymentType)
	if err != nil {
		return nil, err
	}

	var errs []error
	for _, gw := range chain {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		tx, err := gw.CreateTransaction(ctx, req)
		switch {
		case err == nil:
			return tx, nil

		case errors.Is(err, ErrPermanent):
			return nil, err

		case errors.Is(err, ErrNotEnabled):
			errs = append(errs, err)

		case errors.Is(err, ErrCircuitOpen):
			// request never sent — nothing to reconcile, just skip
			errs = append(errs, err)

		default:
			tx, ok := r.reconcile(ctx, gw, req.TransactionDetails.Orderid)
			if ok {
				return tx, nil
			}
			errs = append(errs, err)
		}
	}

	return nil, errors.Join(errs...)
}

// reconcile checks whether a transaction that returned a transient error was
// in fact created. ok is true only when the transaction exists in a state
// that still represents a chargeable/settled payment.
func (r *GatewayRegistry) reconcile(ctx context.Context, gw Gateway, orderId string) (*Transaction, bool) {
	rctx, cancel := context.WithTimeout(ctx, reconcileTimeout)
	defer cancel()

	tx, err := gw.CheckStatus(rctx, orderId)
	if err != nil || tx == nil {
		return nil, false
	}

	switch tx.Status {
	case StatusPending, StatusSettlement, StatusCapture:
		return tx, true
	default:
		return nil, false
	}
}
