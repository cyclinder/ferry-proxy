package config

import (
	"encoding/json"
)

type Metadata struct {
	Name string `json:"name"`
}

type Address struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

type Filter struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}
