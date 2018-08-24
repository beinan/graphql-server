package handler

import (
	"net/http"
)

type Filter func(next http.Handler) http.Handler

func Chain(chain ...Filter) http.Handler {
	if len(chain) == 1 {
		return chain[0](nil)
	}
	rest := chain[1:] //drop the first one: chain[0]
	return chain[0](Chain(rest...))
}
