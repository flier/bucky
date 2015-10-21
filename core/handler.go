package core

import (
	"github.com/fanliao/go-promise"
)

type Request interface {
}

type Response interface {
}

type Handler interface {
	Apply(req Request) *promise.Future
}
