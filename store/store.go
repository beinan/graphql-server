package store

import (
	"context"
	"github.com/beinan/graphql-server/downstreams"
)

var (
	UserType = EntityType{"user", "u"}
	AuthType = EntityType{"authdata", "auth"}

	FriendsRelation = Relationship{UserType, "friends"}
)

type EntityType struct {
	Name      string //for collection name
	ShortName string //for redis key
}

type ID string

type Relationship struct {
	fromEntity   EntityType
	relationName string
}

var MongoStore = &mongoStore{Client: downstreams.MongoClient}
var RedisStore = &redisStore{Client: downstreams.RedisClient}

type DataReader interface {
	GetEnitytByID(context.Context, EntityType, ID, *interface{}) error
	GetEntitiesByIDs(context.Context, EntityType, []ID, *[]interface{}) error
	GetRelation(context.Context, Relationship, ID, *[]ID) error
}

type DataWriter interface {
	AddEntity(context.Context, EntityType, interface{}) error
	UpdateEntity(context.Context, EntityType, interface{}) (interface{}, error)
	AddRelation(context.Context, Relationship, ID, ID) error
}
