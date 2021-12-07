package http_direct

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

func init() {
	registries.StreamHandlerRegistry.Register("http_direct", NewHttpDirect)
}

type Config struct {
	Code int
	Body string
}

func NewHttpDirect(config json.RawMessage) (registries.StreamHandler, error) {
	var conf Config
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return nil, err
	}
	code := conf.Code
	if code == 0 {
		code = http.StatusOK
	}

	resp := http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(conf.Body)),
	}
	buf := bytes.NewBuffer(nil)
	err = resp.Write(buf)
	if err != nil {
		return nil, err
	}
	return &HttpDirect{
		Data: buf.Bytes(),
	}, nil
}
