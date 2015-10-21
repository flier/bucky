package http

import (
	"github.com/flier/bucky/core"
)

type HttpCodec interface {
	core.Codec
}

type httpCodec struct {
	HttpCodec
}

func NewHttpCodec() HttpCodec {
	return &httpCodec{}
}
