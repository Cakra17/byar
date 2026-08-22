package byar

import "errors"

var (
	ErrUnsupportedProvider = errors.New("unsupported payment provider")
	ErrProviderNotExist    = errors.New("provider not exist")
	ErrNoCapabilities      = errors.New("gateway declares no payment methods")
	ErrUnsupportedMethod   = errors.New("payment method not supported")
	ErrEmptyChain          = errors.New("provider chain must not be empty")

	ErrTransient = errors.New("transient error, will retry with another provider")
	ErrPermanent = errors.New("permanent error, the process will abort")
)
