package loader

import (
	"context"

	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/service"
	"github.com/nicksrandall/dataloader"
)

type UserLoader struct {
	UserDAO service.UserService
}

func (ul *UserLoader) GetByID(ctx context.Context, id ID) (*model.User, error) {
	loader, loaderErr := getLoader(ctx, UserDataLoaderKey)
	if loaderErr != nil {
		return nil, loaderErr
	}
	trunk := loader.Load(ctx, dataloader.StringKey(id))
	_, err := trunk()
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Id:     id,
		Name:   "aaaa",
		Gender: "MALE",
	}
	return user, nil
}
