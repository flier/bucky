package bucky

import (
	"github.com/flier/bucky/core"
	"github.com/flier/bucky/http"
)

type Service core.Service

func ServerBuilder() core.ServerBuilder {
	return core.NewServerBuilder()
}

func HttpCodec() core.Codec {
	return http.NewHttpCodec()
}
