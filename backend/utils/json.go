package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

func JSONMarshal(v any, indentSize int) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", strings.Repeat(" ", indentSize))
	err := enc.Encode(v)
	if err != nil {
		return nil, err // nolint:wrapcheck
	}
	return buf.Bytes(), nil
}
