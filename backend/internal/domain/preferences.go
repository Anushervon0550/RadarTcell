package domain

import "encoding/json"

type Preferences struct {
	UserID   string          `json:"user_id"`
	Settings json.RawMessage `json:"settings"`
}
