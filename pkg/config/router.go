package config

import (
	"encoding/json"
)

type RouterConfig struct {
	Metadata
	Router json.RawMessage `json:"router"`
}
