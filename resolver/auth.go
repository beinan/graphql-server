package resolver

import (
	"context"
	"errors"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/utils"
)

type AuthInput struct {
	Id       graphql.ID
	Password string
}

func (r *Resolver) SignUp(ctx context.Context, args *struct {
	Input *AuthInput
}) (*userResolver, error) {
	hashedPassword, hashErr := utils.HashPassword(args.Input.Password)
	if hashErr != nil {
		return nil, hashErr
	}
	err := r.DB.UserDAO().CreateUser(ctx, string(args.Input.Id), hashedPassword)
	user := &model.User{
		Id:     args.Input.Id,
		Name:   "",
		Gender: "MALE",
	}
	return &userResolver{user}, err
}

func (r *Resolver) SignIn(ctx context.Context, args *struct {
	Input *AuthInput
}) (string, error) {
	hashedPassword, err := r.DB.UserDAO().GetHashedPassword(ctx, string(args.Input.Id))
	if err != nil {
		return "", err
	}
	isPasswordMatched := utils.CheckPasswordHash(args.Input.Password, hashedPassword)
	if isPasswordMatched {
		return utils.GenerateJWT(string(args.Input.Id))
	} else {
		return "", errors.New("Username and password do not match")
	}
}
