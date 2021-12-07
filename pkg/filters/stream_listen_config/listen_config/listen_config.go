package ferry_proxy

import (
	"context"
	"net"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"

	"github.com/gogf/greuse"
)

var (
	listenConfig = net.ListenConfig{
		Control: greuse.Control,
	}
)

type ListenConfig struct {
	Network string
	Address string
}

func (d *ListenConfig) Listen(ctx context.Context) (registries.StreamListener, error) {
	return listenConfig.Listen(ctx, d.Network, d.Address)
}
