package domain

import "time"

type Alert struct {
	ID         int       `json:"id"`
	NodeID     uint16    `json:"node_id"`
	EpochTime  time.Time `json:"epoch_time"`
	SequenceNum	uint16	`json:"sequence_num"`
	ThreatClass uint8    `json:"threat_class"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
}
