package rpc

import (
	"reflect"

	"github.com/fanliao/go-promise"
	"github.com/flier/bucky/core"
	"golang.org/x/net/context"
)

var (
	NativeFactory = &nativeFactory{}
)

type nativeFactory struct {
}

var _ = (RpcFactory)((*nativeFactory)(nil))

func (f *nativeFactory) Build(v interface{}) core.Service {
	return NewNativeDispatcher(v)
}

type nativeDispatcher struct {
	metadata Metadata
	target   interface{}
}

func NewNativeDispatcher(v interface{}) *nativeDispatcher {
	return &nativeDispatcher{
		metadata: NewNativeMetadata(reflect.TypeOf(v)),
		target:   v,
	}
}

func (d *nativeDispatcher) Apply(ctxt context.Context, req core.Request) *promise.Future {
	result := promise.NewPromise()

	return result.Future
}

type nativeMetadata struct {
	t reflect.Type
}

func NewNativeMetadata(t reflect.Type) *nativeMetadata {
	return &nativeMetadata{t}
}
