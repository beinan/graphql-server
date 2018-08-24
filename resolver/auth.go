package resolver

import (
	"context"
	"errors"
	"strconv"

	"github.com/beinan/fastid"
	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/utils"
)

type AuthInput struct {
	LoginName string
	Password  string
}

type Auth = model.Auth

func (r *Resolver) SignUp(ctx context.Context, args *struct {
	Input *AuthInput
}) (*userResolver, error) {
	hashedPassword, hashErr := utils.HashPassword(args.Input.Password)
	if hashErr != nil {
		return nil, hashErr
	}
	strUserID := strconv.FormatInt(fastid.GenInt64ID(), 16)
	auth := Auth{
		LoginName:      args.Input.LoginName,
		UserId:         strUserID,
		HashedPassword: hashedPassword,
		IsAdmin:        false,
		Permissions:    make([]string, 0),
	}
	err := r.services.AuthService.Create(ctx, auth)
	if err != nil {
		return nil, err
	}
	user := User{
		Id:     strUserID,
		Name:   args.Input.LoginName,
		Gender: "MALE",
	}
	err = r.services.UserService.Create(ctx, user)
	return &userResolver{&user, r.services}, err
}

func (r *Resolver) SignIn(ctx context.Context, args *struct {
	Input *AuthInput
}) (string, error) {
	auth, err := r.services.AuthService.GetByID(ctx, model.ID(args.Input.LoginName))
	if err != nil {
		return "", err
	}
	isPasswordMatched := utils.CheckPasswordHash(args.Input.Password, auth.HashedPassword)
	if isPasswordMatched {
		return utils.GenerateJWT(utils.AppClaims{
			UserId:      auth.UserId,
			IsAdmin:     auth.IsAdmin,
			Permissions: auth.Permissions,
		})
	} else {
		return "", errors.New("Username and password do not match")
	}
}
