package entity

import (
	"net/http"
	"time"
)

type Request struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Method      string `json:"method"`
	URL         string `json:"url"`
	Headers     string `json:"headers"`
	Body        string `json:"body"`
	Params      string `json:"params"`
	ContentType string `json:"content_type; default=application/json"`
}

type Response struct {
	StatusCode int
	Status     string
	Proto      string
	Headers    http.Header
	Body       string
	Duration   time.Duration
	Size       int
}
