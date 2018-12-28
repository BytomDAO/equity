package equity

import (
	"bytes"
	"encoding/json"
)

// JSONMarshal escapes the special characters from JSON results
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "  ")

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
