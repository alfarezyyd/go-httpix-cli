package utils

import (
	"encoding/json"
	"strings"
)

// tryPrettyJSON parses s as JSON and returns indented form + true on success.
func TryPrettyJSON(s string) (string, bool) {
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
func PrettyJSON(s string) string {
	out, _ := TryPrettyJSON(s)
	return out
}
