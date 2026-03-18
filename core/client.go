// Package httpclient is a thin HTTP layer. It knows nothing about bubbletea —
// callers are responsible for wrapping Execute() in a tea.Cmd.
package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request carries everything needed to fire one HTTP round-trip.
type Request struct {
	Method  string
	URL     string
	Body    string
	Headers string // raw "Key: Value\n" lines
	Params  string // raw "key=value\n" lines — appended as query string
}

// Response carries the result of a completed round-trip.
type Response struct {
	Status   int
	Proto    string
	Headers  http.Header
	Body     string
	Duration time.Duration
	Size     int
}

// Execute performs the HTTP request and returns a Response or an error.
// Run this inside a tea.Cmd goroutine — it blocks until the server replies.
func Execute(r Request) (*Response, error) {
	urlStr, err := buildURL(r.URL, r.Params)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if strings.TrimSpace(r.Body) != "" {
		bodyReader = bytes.NewBufferString(r.Body)
	}

	req, err := http.NewRequest(r.Method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	applyHeaders(req, r.Headers)
	if bodyReader != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "BLINK/1.0 TUI-HTTP-Client")

	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	dur := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return &Response{
		Status:   resp.StatusCode,
		Proto:    resp.Proto,
		Headers:  resp.Header,
		Body:     string(bodyBytes),
		Duration: dur,
		Size:     len(bodyBytes),
	}, nil
}

// buildURL normalises the URL and appends query params.
func buildURL(rawURL, params string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	if q := strings.TrimSpace(params); q != "" {
		sep := "?"
		if strings.Contains(rawURL, "?") {
			sep = "&"
		}
		rawURL += sep + strings.ReplaceAll(q, "\n", "&")
	}
	return rawURL, nil
}

// applyHeaders parses raw "Key: Value\n" lines into the request.
func applyHeaders(req *http.Request, raw string) {
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}
