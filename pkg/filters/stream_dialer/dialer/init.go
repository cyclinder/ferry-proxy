package dialer

import (
	"encoding/json"
	"fmt"

	registries2 "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

func init() {
	registries2.StreamDialerRegistry.Register("dialer", NewDialer)
}

type Config struct {
	Network string `json:"network"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

func NewDialer(config json.RawMessage) (registries2.StreamDialer, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}

	network := conf.Network
	if network == "" {
		network = "tcp"
	}
	d := &Dialer{
		Network: network,
		Address: fmt.Sprintf("%s:%d", conf.Address, conf.Port),
	}
	return d, nil
}
