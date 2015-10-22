package transport

import (
	"github.com/flier/bucky/core"
)

var (
	UnixCodec = &unixCodecFactory{}
)

type unixCodecFactory struct {
}

var _ = (core.CodecFactory)((*unixCodecFactory)(nil))

func (f *unixCodecFactory) ClientCodec(cfg *core.ClientCodecConfig) core.ClientCodec {
	return nil
}

func (f *unixCodecFactory) ServerCodec(cfg *core.ServerCodecConfig) core.ServerCodec {
	return nil
}
