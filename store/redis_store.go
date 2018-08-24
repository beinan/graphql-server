package store

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
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

func (r *redisStore) GetRelationWithCursor(
	ctx context.Context,
	relation Relationship,
	id ID,
	cursor string,
	limit int64,
	isRev bool) ([]IDWithCursor, error) {
	var cursorStr string
	if cursor == "" { //if cursor is Zero string
		cursorStr = "-inf"
	} else {
		cursorStr = cursor
	}
	strCmd := r.Client.ZRangeByScoreWithScores(relation.GenID(id), redis.ZRangeBy{
		Min:    "(" + cursorStr,
		Max:    "+inf",
		Offset: 0,
		Count:  limit,
	})
	strIDs, err := strCmd.Result()
	if err != nil {
		return nil, err
	}
	results := make([]IDWithCursor, len(strIDs))
	for i := range strIDs {
		results[i] = IDWithCursor{
			Id: ID(strIDs[i].Member.(string)),
			//format score (float value) into string
			Cursor: strconv.FormatFloat(strIDs[i].Score, 'f', 0, 64),
		}
	}
	return results, err
}

func (r *redisStore) AddRelation(
	ctx context.Context,
	relation Relationship,
	fromId ID,
	toId ID) error {
	//do not update the existing records
	intCmd := r.Client.ZAddNXCh(relation.GenID(fromId),
		redis.Z{
			Score:  score(),
			Member: toId,
		})
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
