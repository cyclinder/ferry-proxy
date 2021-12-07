package config

import (
	"encoding/json"
)

type EndpointConfig struct {
	Metadata
	Endpoint json.RawMessage `json:"endpoint"`
}
