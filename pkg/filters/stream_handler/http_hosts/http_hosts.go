package http_hosts

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"

	"github.com/wzshiming/cmux"
	"github.com/wzshiming/hostmatcher"
	"github.com/wzshiming/sni"
)

type HttpHosts struct {
	Hosts   map[string]registries.StreamHandler
	Matches []HttpHostMatches
	Default registries.StreamHandler
}

type HttpHostMatches struct {
	Match   hostmatcher.Matcher
	Handler registries.StreamHandler
}

func (h *HttpHosts) ServeStream(ctx context.Context, stm registries.Stream) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	host, err := sni.HTTPHost(io.TeeReader(io.LimitReader(stm, preloadSize), buf))
	if err != nil {
		h.Default.ServeStream(ctx, stm)
		return
	}
	stm = cmux.UnreadConn(stm, buf.Bytes())
	if handler, ok := h.Hosts[host]; ok {
		handler.ServeStream(ctx, stm)
	} else {
		matched := false
		for _, match := range h.Matches {
			if match.Match.Match(host) {
				match.Handler.ServeStream(ctx, stm)
				matched = true
				break
			}
		}
		if !matched {
			h.Default.ServeStream(ctx, stm)
		}
	}
}

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, preloadSize))
		},
	}
)

const (
	preloadSize = 4 * 1024
)
