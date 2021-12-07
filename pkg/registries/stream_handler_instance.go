package registries

import (
	"context"
	"sync"
)

var StreamHandlerInstance = &streamHandlerInstance{
	instance: map[string]StreamHandler{},
}

type streamHandlerInstance struct {
	mut      sync.RWMutex
	instance map[string]StreamHandler
}

func (i *streamHandlerInstance) Get(name string) StreamHandler {
	return &streamHandlerInstanceRef{
		name:     name,
		instance: i,
	}
}

func (i *streamHandlerInstance) Put(name string, s StreamHandler) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.instance[name] = s
}

type streamHandlerInstanceRef struct {
	name     string
	instance *streamHandlerInstance
}

func (i *streamHandlerInstanceRef) ServeStream(ctx context.Context, stm Stream) {
	i.instance.mut.RLock()
	d := i.instance.instance[i.name]
	i.instance.mut.RUnlock()
	d.ServeStream(ctx, stm)
}
