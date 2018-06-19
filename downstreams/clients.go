package downstreams

import (
	"context"
)

type Request struct {
	Header interface{}
	Data   interface{}
}

type Response struct {
}

type Service interface {
	Serve(context.Context, *Request) Response
}
