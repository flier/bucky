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
	bucky.Service

	Uppercase(string) (string, error)

	Count(string) int
}

type stringService struct {
	StringService
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

	bucky.ServerBuilder().
		Codec(bucky.HttpCodec()).
		BindTo(&net.TCPAddr{Port: 8080}).
		Build(&stringService{}).
		Serve(ctxt)
}
