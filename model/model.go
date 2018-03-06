package model

import(	
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	Id     graphql.ID
	Name   string
	Gender string
}
