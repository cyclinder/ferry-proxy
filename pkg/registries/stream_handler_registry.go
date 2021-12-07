package registries

import (
	"encoding/json"
	"fmt"
)

var (
	StreamHandlerRegistry = streamHandlerRegistry{
		registry: map[string]NewStreamHandlerFunc{},
	}
)

type streamHandlerRegistry struct {
	registry map[string]NewStreamHandlerFunc
}

func (r *streamHandlerRegistry) Register(name string, f NewStreamHandlerFunc) {
	_, ok := r.registry[name]
	if ok {
		panic(fmt.Sprintf("stream handler %s already registered", name))
	}
	r.registry[name] = f
}

func (r *streamHandlerRegistry) New(name string, config json.RawMessage) (StreamHandler, error) {
	f, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("stream handler %q not found", name)
	}
	return f(config)
}

type NewStreamHandlerFunc func(config json.RawMessage) (StreamHandler, error)
