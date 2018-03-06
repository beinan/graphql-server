package loader

import (
	"context"
	"errors"
	"fmt"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/nicksrandall/dataloader"

	"github.com/beinan/graphql-server/database"
	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/utils"
)

const (
	UserDataLoaderKey = "user-data-loader-key"
)


//a manager class for all the data loaders
type Loaders struct {
	logger utils.Logger
	db database.DB
}

func NewLoader(db database.DB, logger utils.Logger) Loaders {
	return Loaders{ logger, db }
}

func (l Loaders) Attach(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		UserDataLoaderKey,
		dataloader.NewBatchedLoader(l.userBatchFunc),
	)
	return ctx
}

func getLoader(ctx context.Context, key string) (*dataloader.Loader, error) {
	loader, ok := ctx.Value(key).(*dataloader.Loader)

	if !ok {
		return nil, fmt.Errorf("data loader %s does not exist in request's context", key)
	}

	return loader, nil
}

func LoadUser(ctx context.Context, id string) (*model.User, error) {
	loader, loaderErr := getLoader(ctx, UserDataLoaderKey)
	if loaderErr != nil {
		return nil, loaderErr
	}
	trunk := loader.Load(ctx, dataloader.StringKey(id))
	_, err:= trunk()
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Id:     graphql.ID(id),
		Name:   "aaaa",
		Gender: "MALE",
	}
	return user,nil;
}

func (l Loaders) userBatchFunc(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result = make([]*dataloader.Result, len(keys))
	stringKeys := extractStringKeys(keys)
	records, err := l.db.UserDAO().GetUserByIds(ctx, stringKeys)
	
	if err != nil { //query failed, fill the error into each result
		for i := range stringKeys{
			results[i] = &dataloader.Result{Data: nil, Error: err}
		}
		return results
	}
	
	keyToRecordsIndex := make(map[string]int)
	//create a map for key -> records
	for i,record := range records {
		keyToRecordsIndex[string(record.Id)] = i
	}
	for i, key := range stringKeys{
		recordIndex, ok := keyToRecordsIndex[key]
		if ok {
			results[i] = &dataloader.Result{Data: records[recordIndex], Error: nil}
		} else { //Not found
			results[i] = &dataloader.Result{Data: nil, Error: errors.New("Not found")}
		}
	}
	l.logger.Debugw("User Data Loader", "keys", keys, "results", results)
	return results
}

func extractStringKeys(keys dataloader.Keys) []string {
	results := make([]string, len(keys))
	for i, key := range keys {
		results[i] = key.String()
	}
	return results
}
