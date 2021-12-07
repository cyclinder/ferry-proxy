package http_direct

import (
	"context"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

type HttpDirect struct {
	Data []byte
}

func (h *HttpDirect) ServeStream(ctx context.Context, stm registries.Stream) {
	stm.Write(h.Data)
}
