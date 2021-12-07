package registries

import (
	"encoding/json"
	"fmt"
)

var (
	StreamDialerRegistry = streamDialerRegistry{
		registry: map[string]NewStreamDialerFunc{},
	}
)

type streamDialerRegistry struct {
	registry map[string]NewStreamDialerFunc
}

func (r *streamDialerRegistry) Register(name string, f NewStreamDialerFunc) {
	_, ok := r.registry[name]
	if ok {
		panic(fmt.Sprintf("stream dialer %s already registered", name))
	}
	r.registry[name] = f
}

func (r *streamDialerRegistry) New(name string, config json.RawMessage) (StreamDialer, error) {
	f, ok := r.registry[name]
	if !ok {
		return nil, fmt.Errorf("stream dialer %q not found", name)
	}
	return f(config)
}

type NewStreamDialerFunc func(config json.RawMessage) (StreamDialer, error)
