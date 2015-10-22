package main

import (
	"errors"
	"net"
	"strings"

	"golang.org/x/net/context"

	"github.com/flier/bucky"
)

var ErrEmpty = errors.New("empty string")

type StringService interface {
	Uppercase(string) (string, error)

	Count(string) int
}

type stringService struct {
}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func main() {
	ctxt := context.Background()

	builder := &bucky.ServerBuilder{
		Name:         "stringsvc1",
		Addr:         &net.TCPAddr{Port: 8080},
		CodecFactory: bucky.Http,
		Encoding:     bucky.Json,
	}

	builder.Build(bucky.Rpc(&stringService{})).Serve(ctxt)
}
