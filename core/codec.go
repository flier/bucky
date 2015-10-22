package core

import (
	"crypto/tls"
	"net"
	"net/url"
)

type ClientCodec interface {
	ClientDispatcher(transport Transport) Service
}

type ServerCodec interface {
	ServerDispatcher(transport Transport, service Service) Server
}

type ClientCodecConfig struct {
	Name       string
	Uri        *url.URL
	KeepAlives bool
	TLSConfig  *tls.Config
}

type ServerCodecConfig struct {
	Name              string
	Addr              net.Addr
	KeepAlives        bool
	TLSConfig         *tls.Config
	CertFile, KeyFile string
}

type CodecFactory interface {
	ClientCodec(cfg *ClientCodecConfig) ClientCodec

	ServerCodec(cfg *ServerCodecConfig) ServerCodec
}
