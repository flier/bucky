package transport

import (
	"github.com/flier/bucky/core"
)

var (
	TcpCodec = &tcpCodecFactory{}
)

type tcpCodecFactory struct {
}

var _ = (core.CodecFactory)((*tcpCodecFactory)(nil))

func (f *tcpCodecFactory) ClientCodec(cfg *core.ClientCodecConfig) core.ClientCodec {
	return nil
}

func (f *tcpCodecFactory) ServerCodec(cfg *core.ServerCodecConfig) core.ServerCodec {
	return nil
}
