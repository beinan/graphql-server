package mongodb

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/beinan/graphql-server/model"
)

type UserDAO struct {
	db *MongoDB
}

//User model for mongodb
type User struct {
	Id       string `bson:"_id"`
	Password string `bson:"password"` //hashed password
	Name     string `bson:"name"`
	Gender   string `bson:"gender"`
}

func (dao *UserDAO) getCollection(ctx context.Context) *mgo.Collection {
	return getSession(ctx).DB(dao.db.dbname).C("users")
}

func (dao *UserDAO) GetHashedPassword(ctx context.Context, id string) (string, error) {
	col := dao.getCollection(ctx)
	var record User
	err := col.Find(bson.M{"_id": id}).One(&record)
	if err != nil {
		return "", err
	}
	return record.Password, nil
}

func (dao *UserDAO) GetUserByIds(ctx context.Context, ids []string) ([]model.User, error) {
	col := dao.getCollection(ctx)
	var records []User //records in db
	err := col.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&records)
	if err != nil {
		dao.db.logger.Debugw("Get User by ids failed:", "err", err)
		return nil, err
	}
	var results []model.User = make([]model.User, len(records))
	//convert db model into graphql model
	for i, record := range records {
		results[i] = model.User{
			Id:     record.Id,
			Name:   record.Name,
			Gender: record.Gender,
		}
	}
	dao.db.logger.Debugw("GetUserByIds:", "ids", ids, "results", results)
	return results, nil
}

func (dao *UserDAO) CreateUser(ctx context.Context, id string, password string) error {
	col := dao.getCollection(ctx)
	dao.db.logger.Debugf("Creating  user: %v", id)
	return col.Insert(User{id, password, "", ""})
}
