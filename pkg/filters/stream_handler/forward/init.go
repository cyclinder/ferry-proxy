package forward

import (
	"encoding/json"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

func init() {
	registries.StreamHandlerRegistry.Register("forward", NewForward)
}

type Config struct {
	Endpoint string `json:"endpoint"`
}

func NewForward(config json.RawMessage) (registries.StreamHandler, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}

	target := registries.StreamDialerInstance.Get(conf.Endpoint)
	return &Forward{
		Target: target,
	}, nil
}
