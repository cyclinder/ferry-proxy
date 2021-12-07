package forward

import (
	"encoding/json"

	registries2 "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

func init() {
	registries2.StreamHandlerRegistry.Register("forward", NewForward)
}

type Config struct {
	Endpoint string `json:"endpoint"`
}

func NewForward(config json.RawMessage) (registries2.StreamHandler, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}

	target := registries2.StreamDialerInstance.Get(conf.Endpoint)
	return &Forward{
		Target: target,
	}, nil
}
