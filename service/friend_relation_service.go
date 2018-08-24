package service

import (
	"context"

	"github.com/beinan/graphql-server/model"
)

type FriendRelationService interface {
	Add(context.Context, ID, ID) error
	Get(ctx context.Context, fromId ID, start int64, count int64) ([]ID, error)
	GetWithCursor(ctx context.Context, fromId ID, cursor string, count int64, isRev bool) ([]IDWithCursor, error)
}

type FriendRelationDAO struct {
	Reader RelationDataReader
	Writer RelationDataWriter
}

func (s *FriendRelationDAO) Add(ctx context.Context, fromId ID, toId ID) error {
	return s.Writer.AddRelation(ctx, model.FriendsRelation, fromId, toId)
}

func (s *FriendRelationDAO) Get(
	ctx context.Context,
	fromId ID,
	start int64,
	count int64,
) ([]ID, error) {
	return s.Reader.GetRelation(ctx, model.FriendsRelation, fromId, start, count)
}

func (s *FriendRelationDAO) GetWithCursor(
	ctx context.Context,
	fromId ID,
	cursor string,
	count int64,
	isRev bool,
) ([]IDWithCursor, error) {
	return s.Reader.GetRelationWithCursor(ctx, model.FriendsRelation, fromId, cursor, count, isRev)
}
