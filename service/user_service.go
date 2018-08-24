package service

import (
	"context"

	"github.com/beinan/graphql-server/model"
)

type UserService interface {
	GetByID(context.Context, ID) (*User, error)
	GetByIDs(context.Context, []ID) ([]User, error)
	Create(context.Context, User) error
}

type UserDAO struct {
	Reader EntityDataReader
	Writer EntityDataWriter
}

func (dao *UserDAO) GetByID(ctx context.Context, id ID) (*User, error) {
	user := &User{}
	err := dao.Reader.GetEnitytByID(ctx, model.UserType, id, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) GetByIDs(ctx context.Context, ids []ID) ([]User, error) {
	var users []User
	err := dao.Reader.GetEntitiesByIDs(ctx, model.UserType, ids, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAO) Create(ctx context.Context, user User) error {
	return dao.Writer.AddEntity(ctx, model.UserType, user)
}
