package database

import(
	"context"

	"github.com/beinan/graphql-server/model"
)

type DB interface {
	//Attach a db session to the context
	Attach(ctx context.Context) context.Context
	//Close the db session in the context 
	Close(ctx context.Context)

	UserDAO() UserDAO
}

type UserDAO interface {
	GetUserByIds(ctx context.Context, ids []string) ([]model.User, error)
	CreateUser(ctx context.Context, id string, password string)  error
	GetHashedPassword(ctx context.Context, id string) (string, error)
}

