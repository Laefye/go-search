package dto

import "time"

type SearchQueryEvent struct {
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
}
