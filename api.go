package bucky

import (
	"github.com/flier/bucky/core"
	"github.com/flier/bucky/http"
	"github.com/flier/bucky/rpc"
	"github.com/flier/bucky/transport"
)

type Service core.Service

type Server core.Server

type ServerBuilder core.ServerBuilder

func (b *ServerBuilder) Build(service Service) Server {
	return ((*core.ServerBuilder)(b)).Build(service)
}

var (
	Http = http.HttpCodec
	Tcp  = transport.TcpCodec
	Udp  = transport.UdpCodec
	Unix = transport.UnixCodec

	Json       = core.JsonEncoding
	JsonPretty = core.JsonPrettyEncoding
	Xml        = core.XmlEncoding
	XmlPretty  = core.XmlPrettyEncoding
	Yaml       = core.YamlEncoding
)

func Rpc(v interface{}) core.Service {
	return rpc.NativeFactory.Build(v)
}
