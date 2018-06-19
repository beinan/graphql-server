package resolver

import (
	"context"
	"errors"
	"github.com/beinan/graphql-server/store"
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	Id     graphql.ID `bson:"_id"`
	Name   string
	Gender string
	//Friends *[]User
}

func (r *Resolver) GetUser(ctx context.Context, args *struct {
	Id graphql.ID
}) (*userResolver, error) {
	user := &User{}
	err := store.MongoStore.GetEnitytByID(ctx, store.UserType, store.ID(args.Id), user)
	if err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}

func (r *Resolver) AddFriend(ctx context.Context, args *struct {
	FromId graphql.ID
	ToId   graphql.ID
}) (bool, error) {
	if args.FromId == args.ToId {
		return false, errors.New("You cannot add yourself as a friend")
	}
	//todo: check permission
	user := &User{}
	//get "TO" user
	err := store.MongoStore.GetEnitytByID(ctx, store.UserType, store.ID(args.ToId), user)
	if err != nil {
		return false, err
	}
	err = store.RedisStore.AddRelation(ctx, store.FriendsRelation, store.ID(args.FromId), store.ID(args.ToId))
	if err != nil {
		return false, err
	}
	return true, nil
}

type userResolver struct {
	user *User
}

func (r *userResolver) ID() graphql.ID {
	return r.user.Id
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
	userIds, err := store.RedisStore.GetRelation(ctx, store.FriendsRelation, store.ID(string(r.user.Id)), start, stop)
	if err != nil {
		return nil, err
	}
	var users []User
	err = store.MongoStore.GetEntitiesByIDs(ctx, store.UserType, userIds, &users)
	if err != nil {
		return nil, err
	}
	var resolvers []*userResolver = make([]*userResolver, len(users))
	for i := range users {
		resolvers[i] = &userResolver{&users[i]}
	}
	return &resolvers, nil
}

func (r *userResolver) FriendEdges(ctx context.Context, args *struct {
	Input *EdgesInput
}) (*edgesResolver, error) {
	userIdsWithCursor, err := store.RedisStore.GetRelationWithCursor(ctx, store.FriendsRelation, store.ID(string(r.user.Id)),
		args.Input.Cursor, int64(args.Input.PageSize)+1, args.Input.IsRev)
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
	var users []User
	err = store.MongoStore.GetEntitiesByIDs(ctx, store.UserType, userIds, &users)
	if err != nil {
		return nil, err
	}
	var resolvers []*singleEdgeResolver = make([]*singleEdgeResolver, len(users))
	for i := range users {
		resolvers[i] = &singleEdgeResolver{&idNodeResolver{&userResolver{&users[i]}}, userIdsWithCursor[i].Cursor} //todo: get cursor safer, use a map?
	}

	return &edgesResolver{&resolvers, hasMore}, nil
}
