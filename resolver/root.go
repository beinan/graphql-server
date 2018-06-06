package resolver

import (
	"context"
	
	graphql "github.com/graph-gophers/graphql-go"

	"github.com/beinan/graphql-server/database"
	"github.com/beinan/graphql-server/loader"
	"github.com/beinan/graphql-server/model"
)

//resolver type
type Resolver struct{
	DB database.DB
}

func (r *Resolver) GetUser(ctx context.Context, args *struct {
	Id graphql.ID
}) (*userResolver, error) {
	user,err := loader.LoadUser(ctx, string(args.Id))
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}


type userResolver struct {
	user *model.User
}

func (r *userResolver) ID() graphql.ID {
	return r.user.Id
}

func (r *userResolver) Name() string {
	return r.user.Name
}

func (r *userResolver) Gender() *string {
	return &r.user.Gender
}
