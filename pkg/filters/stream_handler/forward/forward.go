package forward

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/utils"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

type Forward struct {
	Target registries.StreamDialer
}

func (f *Forward) ServeStream(ctx context.Context, stm registries.Stream) {
	target, err := f.Target.Dial(ctx)
	if err != nil {
		log.Println(err)
		stm.Close()
		return
	}

	buf1 := bytesPool.Get().([]byte)
	buf2 := bytesPool.Get().([]byte)
	defer func() {
		bytesPool.Put(buf1)
		bytesPool.Put(buf2)
	}()
	err = tunnel(ctx, stm, target, buf1, buf2)
	if err != nil && !utils.IsClosedConnError(err) {
		log.Println(err)
		return
	}
}

// tunnel create tunnels for two io.ReadWriteCloser
func tunnel(ctx context.Context, c1, c2 io.ReadWriteCloser, buf1, buf2 []byte) error {
	ctx, cancel := context.WithCancel(ctx)
	var errs tunnelErr
	go func() {
		_, errs[0] = io.CopyBuffer(c1, c2, buf1)
		cancel()
	}()
	go func() {
		_, errs[1] = io.CopyBuffer(c2, c1, buf2)
		cancel()
	}()
	<-ctx.Done()
	errs[2] = c1.Close()
	errs[3] = c2.Close()
	errs[4] = ctx.Err()
	if errs[4] == context.Canceled {
		errs[4] = nil
	}
	return errs.FirstError()
}

type tunnelErr [5]error

func (t tunnelErr) FirstError() error {
	for _, err := range t {
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	bytesPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 32*1024)
		},
	}
)
