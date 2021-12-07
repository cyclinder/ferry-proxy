package registries

import (
	"encoding/json"
	"fmt"
)

var (
	StreamListenConfigRegistry = streamListenConfigRegistry{
		registry: map[string]NewStreamListenConfigFunc{},
	}
)

type streamListenConfigRegistry struct {
	registry map[string]NewStreamListenConfigFunc
}

func (r *streamListenConfigRegistry) Register(name string, f NewStreamListenConfigFunc) {
	_, ok := r.registry[name]
	if ok {
		panic(fmt.Sprintf("stream listen config %s already registered", name))
	}
	r.registry[name] = f
}

func (r *streamListenConfigRegistry) New(name string, config json.RawMessage) (StreamListenConfig, error) {
	f, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("stream listen config %q not found", name)
	}
	return f(config)
}

type NewStreamListenConfigFunc func(config json.RawMessage) (StreamListenConfig, error)
