package registries

import (
	"context"
	"net"
)

type (
	// Stream is alias for net.Conn
	Stream = net.Conn

	// StreamListener is alias for net.Listener
	StreamListener = net.Listener

	// StreamHandler is a stream handler
	StreamHandler interface {
		ServeStream(ctx context.Context, stm Stream)
	}

	// StreamListenConfig is a stream listener
	StreamListenConfig interface {
		Listen(ctx context.Context) (StreamListener, error)
	}

	// StreamDialer is a stream dialer
	StreamDialer interface {
		Dial(ctx context.Context) (Stream, error)
	}
)
