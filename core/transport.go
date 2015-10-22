package core

import (
	"net"
	"time"

	"github.com/fanliao/go-promise"
	"golang.org/x/net/context"
)

type Transport interface {
	Read() *promise.Future // []byte

	Write() *promise.Future // int

	Close(deadline time.Time) *promise.Future // nil

	LocalAddr() net.Addr

	RemoteAddr() net.Addr
}

type NetTransport struct {
	ctxt context.Context
	conn net.Conn
}

func NewNetTransport(ctxt context.Context, conn net.Conn) *NetTransport {
	return &NetTransport{ctxt, conn}
}

func (t *NetTransport) Read() *promise.Future {
	result := promise.NewPromise()

	go func() {
		var buf [4096]byte

		read, err := t.conn.Read(buf[:])

		if err != nil {
			result.Reject(err)
		} else {
			result.Resolve(buf[:read])
		}
	}()

	go func() {
		<-t.ctxt.Done()

		result.Cancel()
	}()

	return result.Future
}

func (t *NetTransport) Write(data []byte) *promise.Future {
	result := promise.NewPromise()

	go func() {
		wrote, err := t.conn.Write(data)

		if err != nil {
			result.Reject(err)
		} else {
			result.Resolve(wrote)
		}
	}()

	go func() {
		<-t.ctxt.Done()

		result.Cancel()
	}()

	return result.Future
}

func (t *NetTransport) Close(deadline time.Time) *promise.Future {
	result := promise.NewPromise()

	if err := t.conn.Close(); err != nil {
		result.Reject(err)
	} else {
		result.Resolve(nil)
	}

	return result.Future
}

func (t *NetTransport) LocalAddr() net.Addr { return t.conn.LocalAddr() }

func (t *NetTransport) RemoteAddr() net.Addr { return t.conn.RemoteAddr() }
