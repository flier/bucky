package core

import (
	"crypto/tls"
	"net"
)

type ServiceBuilder interface {
	Build(service Service) Server
}

type ServerBuilder struct {
	Name              string
	Addr              net.Addr
	Backlog           int
	Daemon            bool
	Encoding          Encoding
	Codec             ServerCodec
	CodecFactory      CodecFactory
	TLSConfig         *tls.Config
	CertFile, KeyFile string
	Transport         Transport
}

func (b *ServerBuilder) Build(service Service) Server {
	Assert(b.Name, "No Name was specified")
	Assert(b.Addr, "No Addr was specified")

	if b.Codec == nil && b.CodecFactory != nil {
		b.Codec = b.CodecFactory.ServerCodec(&ServerCodecConfig{
			Name:      b.Name,
			Addr:      b.Addr,
			TLSConfig: b.TLSConfig,
			CertFile:  b.CertFile,
			KeyFile:   b.KeyFile,
		})
	}

	Assert(b.Codec, "No Codec was specified")

	return b.Codec.ServerDispatcher(b.Transport, service)
}
