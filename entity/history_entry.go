package entity

import "time"

// HistoryEntry records one completed request shown in the sidebar.
type HistoryEntry struct {
	Method   string
	URL      string
	Status   int
	Duration time.Duration
	At       time.Time
}
