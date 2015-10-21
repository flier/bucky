package core

import (
	"golang.org/x/net/context"
)

type Server interface {
	Serve(ctxt context.Context) error
}
