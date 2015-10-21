package core

import (
	"net"
)

type ServerBuilder interface {
	Codec(codecs ...Codec) ServerBuilder

	BindTo(addrs ...net.Addr) ServerBuilder

	Build(services ...Service) Server
}

type serverBuilder struct {
	codecs    []Codec
	addresses []net.Addr
	services  []Service
}

func NewServerBuilder() ServerBuilder {
	return &serverBuilder{}
}

func (b *serverBuilder) Codec(codecs ...Codec) ServerBuilder {
	b.codecs = append(b.codecs, codecs...)

	return b
}

func (b *serverBuilder) BindTo(addres ...net.Addr) ServerBuilder {
	b.addresses = append(b.addresses, addres...)

	return b
}

func (b *serverBuilder) Build(services ...Service) Server {
	b.services = append(b.services, services...)

	return nil
}
