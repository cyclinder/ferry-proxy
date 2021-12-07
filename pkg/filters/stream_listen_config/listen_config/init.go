package ferry_proxy

import (
	"encoding/json"
	"fmt"

	registries2 "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

func init() {
	registries2.StreamListenConfigRegistry.Register("listen_config", NewListenConfig)
}

type Config struct {
	Network string `json:"network"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

func NewListenConfig(config json.RawMessage) (registries2.StreamListenConfig, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}

	network := conf.Network
	if network == "" {
		network = "tcp"
	}
	d := &ListenConfig{
		Network: network,
		Address: fmt.Sprintf("%s:%d", conf.Address, conf.Port),
	}
	return d, nil
}
