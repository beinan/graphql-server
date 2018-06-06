package handler

import (
	"encoding/json"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	//"github.com/beinan/graphql_go_boilerplate/loader"
	"github.com/beinan/graphql-server/utils"
)


type GraphQLHandler struct {
	schema *graphql.Schema
	logger utils.Logger
	next http.Handler
}


func HandleGraphQL(schema *graphql.Schema, logger utils.Logger) Filter {
	return func(next http.Handler) http.Handler {
		return &GraphQLHandler {schema, logger, next}			
	}
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	h.logger.Debugw("GraphQL query:",
		"query", params.Query,
		"operation", params.OperationName,
		"vars", params.Variables, 
	)
				
	response := h.schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
		
