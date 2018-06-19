package store

import (
	"context"
	"github.com/beinan/graphql-server/downstreams"
	//	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoStore struct {
	Client *downstreams.MongoDB
}

func (m *mongoStore) GetEnitytByID(
	ctx context.Context,
	entityType EntityType,
	id ID,
	result interface{},
) error {
	col := m.Client.GetCollection(ctx, entityType.Name)
	err := col.Find(bson.M{"_id": id}).One(result)
	return err
}

func (m *mongoStore) GetEntitiesByIDs(
	ctx context.Context,
	entityType EntityType,
	ids []ID,
	results interface{}, //might be an array of entities
) error {
	col := m.Client.GetCollection(ctx, entityType.Name)
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(results)
	return err
}

func (m *mongoStore) GetEnitytByQuery(
	ctx context.Context,
	entityType EntityType,
	query bson.M, //map[string]interface{}
	result interface{},
) error {
	col := m.Client.GetCollection(ctx, entityType.Name)
	err := col.Find(query).One(result)
	return err
}

func (m *mongoStore) AddEntity(
	ctx context.Context,
	entityType EntityType,
	entity interface{},
) error {
	col := m.Client.GetCollection(ctx, entityType.Name)
	err := col.Insert(entity)
	return err
}

/**
func (m *MongoStore) GetEntityByIDs(context.Context, []ID) []interface{} {

}

func (m *MongoStore) GetRelation(context.Context, EntityType, ID, RelationName) []String {

}
*/
