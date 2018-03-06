package handler

import (
	"net/http"
	
	"github.com/beinan/graphql-server/loader"
)

func AttachLoader(loaders loader.Loaders) Filter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//Attach dataloader to context
			ctx := loaders.Attach(req.Context())
			//Using the request with the changed context
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
