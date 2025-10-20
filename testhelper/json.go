package testhelper

import (
	"encoding/json"
	"testing"
)

func MarshalJSON(t *testing.T, v any) (b []byte) {
	t.Helper()
	b, _ = json.MarshalIndent(v, "", "  ")
	return
}
