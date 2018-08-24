package service

import (
	"context"

	"github.com/beinan/graphql-server/model"
)

type AuthService interface {
	GetByID(context.Context, ID) (*Auth, error)
	Create(context.Context, Auth) error
}

type AuthDAO struct {
	Reader EntityDataReader
	Writer EntityDataWriter
}

func (dao *AuthDAO) GetByID(ctx context.Context, id ID) (*Auth, error) {
	auth := &Auth{}
	err := dao.Reader.GetEnitytByID(ctx, model.AuthType, id, auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (dao *AuthDAO) Create(ctx context.Context, auth Auth) error {
	return dao.Writer.AddEntity(ctx, model.AuthType, auth)
}
