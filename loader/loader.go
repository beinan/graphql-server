package loader

import (
	"context"
	"fmt"

	"github.com/nicksrandall/dataloader"
)

const (
	UserDataLoaderKey = "user-data-loader-key"
)

type Clients interface {
	//cache
	//db
}

type Loaders struct {
	key         string
	downstreams Clients
}

func LoaderFactory(key string, downstreams Clients) Clients {
	return Loaders{
		key:         key,
		downstreams: downstreams,
	}
}

func (l Loaders) batchFunc(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	// do some aync work to get data for specified keys
	// append to this list resolved values
	return results
}

func (l Loaders) Attach(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		l.key,
		dataloader.NewBatchedLoader(l.batchFunc),
	)
}

func GetLoader(key string, ctx context.Context) (*dataloader.Loader, error) {
	loader, ok := ctx.Value(key).(*dataloader.Loader)

	if !ok {
		return nil, fmt.Errorf("data loader %s does not exist in request's context", key)
	}

	return loader, nil
}
