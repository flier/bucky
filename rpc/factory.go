package rpc

import (
	"github.com/flier/bucky/core"
)

type RpcFactory interface {
	Build(v interface{}) core.Service
}
