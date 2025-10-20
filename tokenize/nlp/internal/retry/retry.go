package retry

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/genproto/googleapis/api/error_reason"
)

func Do(ctx context.Context, req func() error) error {
	const (
		// Retry request up to retries times if the rate limit is exceeded with
		// an increasing backoff.
		retries = 6
		backoff = 180 // In seconds
	)

	var (
		try   = 0
		delay = 1 // In seconds
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Retrier exhausted
		if try >= retries {
			return errRetriesExhausted
		}

		// Request
		err := req()
		var e *apierror.APIError
		if errors.As(err, &e) {
			// Check for rate limit exceeded error.
			if e.Reason() != error_reason.ErrorReason_RATE_LIMIT_EXCEEDED.String() {
				// Other error
				return err
			}

			// Exponentially back-off
			time.Sleep(time.Second * time.Duration(delay))

			// Increase delay and try.
			delay = min(
				backoff,
				delay*delay+rand.IntN(10-1)+1, // delay² + jitter[1, 10)
			)
			try++

			continue
		}

		return err
	}
}
