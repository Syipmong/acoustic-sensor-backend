package domain

import "time"

type Alert struct {
	ID         int       `json:"id"`
	NodeID     string    `json:"node_id"`
	ThreatType string    `json:"threat_type"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Timestamp  time.Time `json:"timestamp"`
}
