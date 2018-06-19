package store

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type redisStore struct {
	Client *redis.Client
}

func (r *redisStore) GetRelation(
	ctx context.Context,
	relation Relationship,
	id ID,
	start int64,
	count int64) ([]ID, error) {
	strCmd := r.Client.ZRange(relation.GenID(id), start, count)
	strIDs, err := strCmd.Result()
	if err != nil {
		return nil, err
	}
	results := make([]ID, len(strIDs))
	for i := range strIDs {
		results[i] = ID(strIDs[i])
	}
	return results, err
}

type IDWithCursor struct {
	Id     ID
	Cursor string
}

func (r *redisStore) GetRelationWithCursor(
	ctx context.Context,
	relation Relationship,
	id ID,
	cursor *string,
	limit int64,
	isRev *bool) ([]IDWithCursor, error) {
	var cursorStr string
	if cursor == nil {
		cursorStr = "-inf"
	} else {
		cursorStr = *cursor
	}
	strCmd := r.Client.ZRangeByScoreWithScores(relation.GenID(id), redis.ZRangeBy{
		"(" + cursorStr, "+inf", 0, limit,
	})
	strIDs, err := strCmd.Result()
	if err != nil {
		return nil, err
	}
	results := make([]IDWithCursor, len(strIDs))
	for i := range strIDs {
		results[i] = IDWithCursor{ID(strIDs[i].Member.(string)), strconv.FormatFloat(strIDs[i].Score, 'f', 0, 64)}
	}
	return results, err
}

func (r *redisStore) AddRelation(
	ctx context.Context,
	relation Relationship,
	fromId ID,
	toId ID) error {
	//do not update the existing records
	intCmd := r.Client.ZAddNXCh(relation.GenID(fromId), redis.Z{score(), string(toId)})
	num, err := intCmd.Result()
	if num == 0 && err == nil {
		return errors.New("A friend with the same id has already been added.")
	}
	return err
}

//calculate score from time
func score() float64 {
	return float64(time.Now().UnixNano())
}

func (r *Relationship) GenID(fromId ID) string {
	return r.fromEntity.ShortName + ":" + string(fromId) + ":" + r.relationName
}
