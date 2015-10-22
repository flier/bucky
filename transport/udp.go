package transport

import (
	"github.com/flier/bucky/core"
)

var (
	UdpCodec = &udpCodecFactory{}
)

type udpCodecFactory struct {
}

var _ = (core.CodecFactory)((*udpCodecFactory)(nil))

func (f *udpCodecFactory) ClientCodec(cfg *core.ClientCodecConfig) core.ClientCodec {
	return nil
}

func (f *udpCodecFactory) ServerCodec(cfg *core.ServerCodecConfig) core.ServerCodec {
	return nil
}
