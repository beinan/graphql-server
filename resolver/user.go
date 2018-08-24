package resolver

import (
	"context"
	"errors"
	"fmt"

	"github.com/beinan/graphql-server/model"
	"github.com/beinan/graphql-server/service"
	"github.com/beinan/graphql-server/store"
	"github.com/beinan/graphql-server/utils"
	graphql "github.com/graph-gophers/graphql-go"
)

var logger = utils.DefaultLogger

type User = model.User

func (r *Resolver) GetUser(ctx context.Context, args *struct {
	Id graphql.ID
}) (*userResolver, error) {
	user, err := r.services.UserService.GetByID(ctx, model.ID(args.Id))
	logger.Debugf("Got user(id:%v) %v err:%v", args.Id, user, err)
	if err != nil {
		return nil, err
	}
	return &userResolver{user, r.services}, nil
}

func (r *Resolver) AddFriend(ctx context.Context, args *struct {
	FromId graphql.ID
	ToId   graphql.ID
}) (bool, error) {
	auth := utils.GetAuthObject(ctx)
	if auth.UserId != string(args.FromId) && auth.IsAdmin != true {
		return false, fmt.Errorf("You have no permission to add friend for user %v", args.FromId)
	}
	if args.FromId == args.ToId {
		return false, errors.New("You cannot add yourself as a friend")
	}

	//get "TO" user
	_, err := r.services.UserService.GetByID(ctx, model.ID(args.ToId))
	if err != nil {
		return false, err
	}
	err = r.services.FriendRelationService.Add(ctx, model.ID(args.FromId), model.ID(args.ToId))
	if err != nil {
		return false, err
	}
	return true, nil
}

type userResolver struct {
	user     *User
	services *service.Services
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(r.user.Id)
}

func (r *userResolver) Name() string {
	return r.user.Name
}

func (r *userResolver) Gender() *string {
	return &r.user.Gender
}

func (r *userResolver) Friends(ctx context.Context, args *struct {
	PageNum  int32
	PageSize int32
}) (*[]*userResolver, error) {
	start := int64(args.PageNum * args.PageSize)
	stop := start + int64(args.PageSize) - 1
	userIds, err := store.RedisStore.GetRelation(ctx, model.FriendsRelation, model.ID(string(r.user.Id)), start, stop)
	if err != nil {
		return nil, err
	}
	users, err := r.services.UserService.GetByIDs(ctx, userIds)
	if err != nil {
		return nil, err
	}
	var resolvers []*userResolver = make([]*userResolver, len(users))
	for i := range users {
		resolvers[i] = &userResolver{&users[i], r.services}
	}
	return &resolvers, nil
}

func (r *userResolver) FriendEdges(ctx context.Context, args *struct {
	Input *EdgesInput
}) (*edgesResolver, error) {
	var cursor string
	if args.Input.Cursor == nil {
		cursor = *args.Input.Cursor
	}
	var isRev bool
	if args.Input.IsRev != nil {
		isRev = *args.Input.IsRev
	}
	userIdsWithCursor, err := store.RedisStore.GetRelationWithCursor(ctx, model.FriendsRelation, model.ID(string(r.user.Id)),
		cursor, int64(args.Input.PageSize)+1, isRev)
	if err != nil {
		return nil, err
	}
	var hasMore = false
	if len(userIdsWithCursor) > int(args.Input.PageSize) {
		hasMore = true
		userIdsWithCursor = userIdsWithCursor[:args.Input.PageSize]
	}

	var userIds = make([]store.ID, len(userIdsWithCursor))
	for i := range userIdsWithCursor {
		userIds[i] = userIdsWithCursor[i].Id
	}
	users, err := r.services.UserService.GetByIDs(ctx, userIds)
	if err != nil {
		return nil, err
	}
	var resolvers []*singleEdgeResolver = make([]*singleEdgeResolver, len(users))
	for i := range users {
		resolvers[i] = &singleEdgeResolver{&idNodeResolver{&userResolver{&users[i], r.services}}, userIdsWithCursor[i].Cursor} //todo: get cursor safer, use a map?
	}

	return &edgesResolver{&resolvers, hasMore}, nil
}
