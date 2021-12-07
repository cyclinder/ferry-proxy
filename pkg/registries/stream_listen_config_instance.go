package registries

import (
	"context"
	"sync"
)

var StreamListenConfigInstance = &streamListenConfigInstance{
	instance: map[string]StreamListenConfig{},
}

type streamListenConfigInstance struct {
	mut      sync.RWMutex
	instance map[string]StreamListenConfig
}

func (i *streamListenConfigInstance) Get(name string) StreamListenConfig {
	return &streamListenConfigInstanceRef{
		name:     name,
		instance: i,
	}
}

func (i *streamListenConfigInstance) Put(name string, s StreamListenConfig) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.instance[name] = s
}

type streamListenConfigInstanceRef struct {
	name     string
	instance *streamListenConfigInstance
}

func (i *streamListenConfigInstanceRef) Listen(ctx context.Context) (StreamListener, error) {
	i.instance.mut.RLock()
	d := i.instance.instance[i.name]
	i.instance.mut.RUnlock()
	return d.Listen(ctx)
}
