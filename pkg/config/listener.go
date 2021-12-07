package config

import (
	"encoding/json"
)

type ListenerConfig struct {
	Metadata
	Listener json.RawMessage `json:"listener"`
	Router   string          `json:"router"`
}
