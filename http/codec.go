package http

import (
	"github.com/flier/bucky/core"
)

var (
	HttpCodec = (core.CodecFactory)(&httpCodecFactory{})
)

type httpCodecFactory struct {
}

var _ = (core.CodecFactory)((*httpCodecFactory)(nil))

func (f *httpCodecFactory) ClientCodec(cfg *core.ClientCodecConfig) core.ClientCodec {
	return NewHttpClientCodec(cfg)
}

func (f *httpCodecFactory) ServerCodec(cfg *core.ServerCodecConfig) core.ServerCodec {
	return NewHttpServerCodec(cfg)
}
