package dialer

import (
	"context"
	"net"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

var dialer = net.Dialer{}

type Dialer struct {
	Network string
	Address string
}

func (d *Dialer) Dial(ctx context.Context) (registries.Stream, error) {
	return dialer.DialContext(ctx, d.Network, d.Address)
}
