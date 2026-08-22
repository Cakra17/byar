package midtrans

import (
	"fmt"

	"github.com/cakra17/byar"
	"github.com/midtrans/midtrans-go"
)

// ParseError classifies a Midtrans API error:
// 5xx (and 429 rate limiting) are transient -> safe to fail over;
// everything else (validation errors, declines) is permanent -> abort.
func ParseError(m midtrans.Error) error {
	if m.StatusCode >= 500 || m.StatusCode == 429 {
		return fmt.Errorf("%w: %s (status %d)", byar.ErrTransient, m.Message, m.StatusCode)
	}
	return fmt.Errorf("%w: %s (status %d)", byar.ErrPermanent, m.Message, m.StatusCode)
}
