package store

import (
	"github.com/beinan/graphql-server/downstreams"
	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/utils"
)

type ID = model.ID
type IDWithCursor = model.IDWithCursor
type EntityType = model.EntityType
type Relationship = model.Relationship

func MkMongoStore() *mongoStore {
	//using dbname graphql_devel or graphql_prod
	mongoClient := downstreams.NewMongoClient("graphql_" + utils.Env)
	return &mongoStore{Client: mongoClient}
}

var RedisStore = &redisStore{Client: downstreams.RedisClient}
