package tui

import (
	"encoding/json"
	"strings"
)

// tryPrettyJSON parses s as JSON and returns indented form + true on success.
func tryPrettyJSON(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return "", false
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", false
	}
	return string(out), true
}

// prettyJSON returns indented JSON or "" if s is invalid JSON.
func prettyJSON(s string) string {
	out, _ := tryPrettyJSON(s)
	return out
}
