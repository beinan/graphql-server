package handler

import (
	"net/http"
	"time"
	
	"github.com/beinan/graphql-server/utils"
)

func LatencyStat(logger utils.Logger) Filter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(res, req)
			logger.Debugw("Http request latency:", "latency", time.Since(start), "url", req.URL.Path)
		})
	}
}
