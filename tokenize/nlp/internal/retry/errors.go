package retry

import "errors"

var errRetriesExhausted = errors.New("max retries reached")
