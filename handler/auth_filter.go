package handler

import (
	"net/http"
	"strings"

	"github.com/beinan/graphql-server/utils"
)

func AuthFilter(logger utils.Logger) Filter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			logger.Debugf("Token in header %v", token)
			//Attach auth object to context
			ctx := utils.AuthAttach(req.Context(), token)
			//Using the request with the changed context
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
