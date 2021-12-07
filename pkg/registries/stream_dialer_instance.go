package registries

import (
	"context"
	"sync"
)

var StreamDialerInstance = &streamDialerInstance{
	instance: map[string]StreamDialer{},
}

type streamDialerInstance struct {
	mut      sync.RWMutex
	instance map[string]StreamDialer
}

func (i *streamDialerInstance) Get(name string) StreamDialer {
	return &streamDialerInstanceRef{
		name:     name,
		instance: i,
	}
}

func (i *streamDialerInstance) Put(name string, s StreamDialer) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.instance[name] = s
}

type streamDialerInstanceRef struct {
	name     string
	instance *streamDialerInstance
}

func (i *streamDialerInstanceRef) Dial(ctx context.Context) (Stream, error) {
	i.instance.mut.RLock()
	d := i.instance.instance[i.name]
	i.instance.mut.RUnlock()
	return d.Dial(ctx)
}
