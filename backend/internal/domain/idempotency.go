package domain

import (
	"time"
)

type IdempotencyRecord struct {
	Key        string    `json:"key"`
	FirstSeen  time.Time `json:"first_seen_at"`
	ResultRef  string    `json:"result_ref"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func NewIdempotencyRecord(key, resultRef string, ttl time.Duration) IdempotencyRecord {
	now := time.Now()
	return IdempotencyRecord{
		Key:       key,
		FirstSeen: now,
		ResultRef: resultRef,
		ExpiresAt: now.Add(ttl),
	}
}
