package types

import "time"

type Record struct {
	ID         int       `json:"id"`
	InsertedAt time.Time `json:"inserted_at"`
}

type Key struct {
	ID         int       `json:"id"`
	Value      string    `json:"value"`
	IsInUse    bool      `json:"is_in_use"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type Signature struct {
	ID         int       `json:"id"`
	RecordID   int       `json:"record_id"`
	KeyID      int       `json:"key_id"`
	Value      string    `json:"value"`
	InsertedAt time.Time `json:"inserted_at"`
}
