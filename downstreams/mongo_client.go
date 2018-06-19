package downstreams

import (
	"context"
	mgo "gopkg.in/mgo.v2"

	"github.com/beinan/graphql-server/utils"
)

const (
	mongoDBSessionKey = "mongodb_session_key"
)

var logger = utils.DefaultLogger

//using dbname graphql_devel or graphql_prod
var MongoClient = newClient("graphql_" + utils.Env)

func newClient(dbname string) *MongoDB {
	uri, dbname := "mongo", dbname
	logger.Infof("Dialing mongodb: %s", uri)
	session, err := mgo.Dial(uri)
	if err != nil {
		logger.Info(err)
		panic("mongodb dialing failed!")
	}
	return &MongoDB{session, dbname}
}

type MongoDB struct {
	session *mgo.Session //original session
	dbname  string
}

func (db *MongoDB) Attach(ctx context.Context) context.Context {
	//reuses the same socket as the original session
	//for long operations, pls use session.Copy() instead
	session := db.session.Clone()
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	session.SetMode(mgo.Monotonic, true)
	return context.WithValue(
		ctx,
		mongoDBSessionKey,
		session,
	)
}

func (db *MongoDB) GetSession(ctx context.Context) *mgo.Session {
	return ctx.Value(mongoDBSessionKey).(*mgo.Session)
}

func (db *MongoDB) GetCollection(ctx context.Context, name string) *mgo.Collection {
	return db.GetSession(ctx).DB(db.dbname).C(name)
}

//Close the mongodb session attached in the context
func (db *MongoDB) Close(ctx context.Context) {
	db.GetSession(ctx).Close()
}
