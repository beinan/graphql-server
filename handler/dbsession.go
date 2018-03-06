package handler

import (
	"net/http"

	"github.com/beinan/graphql-server/database"
	"github.com/beinan/graphql-server/utils"
)

func DatabaseHandler(db database.DB, logger utils.Logger) Filter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			//Attach database session to context
			ctx := db.Attach(req.Context())
			//close database session
			defer db.Close(ctx)
			//Using the request with the changed context
			next.ServeHTTP(res, req.WithContext(ctx))			
		})
	}
}

