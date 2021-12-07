package config

import (
	"encoding/json"
)

type ClusterConfig struct {
	Metadata
	Cluster []json.RawMessage `json:"cluster"`
}
