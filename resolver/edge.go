package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
)

type edgesResolver struct {
	singleEdgeResolvers *[]*singleEdgeResolver
	hasMore             bool
}

type singleEdgeResolver struct {
	idNodeResolver *idNodeResolver
	cursor         string
}

type EdgesInput struct {
	Cursor   *string
	PageSize int32
	IsRev    *bool
}

type idNodeInterface interface {
	ID() graphql.ID
}

type idNodeResolver struct {
	idNodeInterface
}

func (r *idNodeResolver) ToUser() (*userResolver, bool) {
	u, ok := r.idNodeInterface.(*userResolver)
	return u, ok
}

func (r *edgesResolver) Edges() (*[]*singleEdgeResolver, error) {
	return r.singleEdgeResolvers, nil
}

func (r *edgesResolver) HasMore() bool {
	return r.hasMore
}

func (r *singleEdgeResolver) Node() (*idNodeResolver, error) {
	return r.idNodeResolver, nil
}

func (r *singleEdgeResolver) Cursor() string {
	return r.cursor
}
