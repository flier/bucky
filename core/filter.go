package core

type Filter interface {
	Handler

	Next() Handler
}
