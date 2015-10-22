package http

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/flier/bucky/core"
)

type httpServerCodec struct {
	*http.Server

	CertFile, KeyFile string
}

var _ = (core.ServerCodec)((*httpServerCodec)(nil))

func NewHttpServerCodec(cfg *core.ServerCodecConfig) *httpServerCodec {
	server := &http.Server{
		Addr:      cfg.Addr.String(),
		Handler:   http.NewServeMux(),
		TLSConfig: cfg.TLSConfig,
	}

	server.SetKeepAlivesEnabled(cfg.KeepAlives)

	return &httpServerCodec{Server: server, CertFile: cfg.CertFile, KeyFile: cfg.KeyFile}
}

func (c *httpServerCodec) ServerDispatcher(transport core.Transport, service core.Service) core.Server {
	return &httpServerDispatcher{
		Server:    c.Server,
		CertFile:  c.CertFile,
		KeyFile:   c.KeyFile,
		Transport: transport,
		Service:   service,
	}
}

type httpServerDispatcher struct {
	*http.Server
	CertFile, KeyFile string
	Transport         core.Transport
	Service           core.Service
}

var _ = (core.Server)((*httpServerDispatcher)(nil))

func (d *httpServerDispatcher) Serve(ctxt context.Context) error {
	if d.TLSConfig != nil {
		return d.ListenAndServeTLS(d.CertFile, d.KeyFile)
	}

	return d.ListenAndServe()
}
