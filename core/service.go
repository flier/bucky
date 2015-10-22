package core

import (
	"github.com/fanliao/go-promise"
	"golang.org/x/net/context"
)

type Request interface {
}

type Response interface {
}

type Service interface {
	Apply(ctxt context.Context, req Request) *promise.Future
}
