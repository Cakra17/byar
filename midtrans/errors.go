package midtrans

import (
	"fmt"

	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
)

// ParseError classifies a Midtrans API error:
// 5xx (and 429 rate limiting) are transient -> safe to fail over;
// 401/403 access denied means the method is not activated for the merchant
// -> skip this provider but try the next one;
// everything else (validation errors, declines) is permanent -> abort.
func ParseError(m midtrans.Error) error {
	switch {
	case m.StatusCode >= 500 || m.StatusCode == 429:
		return fmt.Errorf("%w: %s (status %d)", byar.ErrTransient, m.Message, m.StatusCode)
	case m.StatusCode == 401 || m.StatusCode == 403:
		return fmt.Errorf("%w: %s (status %d)", byar.ErrNotEnabled, m.Message, m.StatusCode)
	default:
		return fmt.Errorf("%w: %s (status %d)", byar.ErrPermanent, m.Message, m.StatusCode)
	}
}
