package http

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/flier/bucky/core"
)

type httpClientCodec struct {
	*http.Client
}

var _ = (core.ClientCodec)((*httpClientCodec)(nil))

func NewHttpClientCodec(cfg *core.ClientCodecConfig) *httpClientCodec {
	jar, err := cookiejar.New(nil)

	core.Assert(err, "fail to create cookiejar")

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: !cfg.KeepAlives,
			TLSClientConfig:   cfg.TLSConfig,
		},
		Jar: jar,
	}

	return &httpClientCodec{client}
}

func (c *httpClientCodec) ClientDispatcher(transport core.Transport) core.Service {
	return nil
}
