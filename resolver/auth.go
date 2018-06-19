package resolver

import (
	"context"
	"errors"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/beinan/fastid"
	"github.com/beinan/graphql-server/store"
	"github.com/beinan/graphql-server/utils"
)

type AuthInput struct {
	LoginName string
	Password  string
}

type Auth struct {
	LoginName      string `bson:"_id"`
	UserId         string
	HashedPassword string
}

func (r *Resolver) SignUp(ctx context.Context, args *struct {
	Input *AuthInput
}) (*userResolver, error) {
	hashedPassword, hashErr := utils.HashPassword(args.Input.Password)
	if hashErr != nil {
		return nil, hashErr
	}
	strUserID := strconv.FormatInt(fastid.GenInt64ID(), 16)
	auth := Auth{args.Input.LoginName, strUserID, hashedPassword}
	err := store.MongoStore.AddEntity(ctx, store.AuthType, auth)
	if err != nil {
		return nil, err
	}
	user := &User{
		Id:     graphql.ID(strUserID),
		Name:   args.Input.LoginName,
		Gender: "MALE",
	}
	err = store.MongoStore.AddEntity(ctx, store.UserType, user)
	return &userResolver{user}, err
}

func (r *Resolver) SignIn(ctx context.Context, args *struct {
	Input *AuthInput
}) (string, error) {
	var auth = Auth{}
	err := store.MongoStore.GetEnitytByID(ctx, store.AuthType, store.ID(args.Input.LoginName), &auth)
	if err != nil {
		return "", err
	}
	isPasswordMatched := utils.CheckPasswordHash(args.Input.Password, auth.HashedPassword)
	if isPasswordMatched {
		return utils.GenerateJWT(string(args.Input.LoginName))
	} else {
		return "", errors.New("Username and password do not match")
	}
}
