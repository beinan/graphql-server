package service

import (
	"context"

	"github.com/beinan/graphql-server/model"
)

type ID = model.ID
type IDWithCursor = model.IDWithCursor

type User = model.User
type Auth = model.Auth

type EntityType = model.EntityType
type Relationship = model.Relationship

type Services struct {
	UserService           UserService
	AuthService           AuthService
	FriendRelationService FriendRelationService
}

type EntityDataReader interface {
	GetEnitytByID(context.Context, EntityType, ID, interface{}) error
	GetEntitiesByIDs(context.Context, EntityType, []ID, interface{}) error
}

type RelationDataReader interface {
	GetRelationWithCursor(
		ctx context.Context,
		relation Relationship,
		id ID,
		cursor string,
		limit int64,
		isRev bool,
	) ([]IDWithCursor, error)
	GetRelation(
		ctx context.Context,
		relation Relationship,
		id ID,
		start int64,
		count int64,
	) ([]ID, error)
}

type EntityDataWriter interface {
	AddEntity(context.Context, EntityType, interface{}) error
	//UpdateEntity(context.Context, EntityType, interface{}) (interface{}, error)
}

type RelationDataWriter interface {
	AddRelation(context.Context, Relationship, ID, ID) error
}
