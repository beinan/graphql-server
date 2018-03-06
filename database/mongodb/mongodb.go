package mongodb

import(
	"context"
	
	mgo "gopkg.in/mgo.v2"

	"github.com/beinan/graphql-server/database"
	"github.com/beinan/graphql-server/utils"
)

const (
	mongoDBSessionKey = "mongodb_session_key"
)

func NewDB(logger utils.Logger) database.DB {
	uri, dbname := "localhost", "devdb"
	logger.Infof("Dialing mongodb: %v", uri)
	session,err := mgo.Dial(uri)
	if err != nil {
		logger.Info(err)
	}
	return &MongoDB{uri, dbname, logger, session}
}

type MongoDB struct {
	uri string
	dbname string
	logger utils.Logger
	session *mgo.Session //original session
}

func (db *MongoDB) UserDAO() database.UserDAO {
	return &UserDAO{db}
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

//Close the mongodb session attached in the context
func (db *MongoDB ) Close(ctx context.Context) {
	getSession(ctx).Close()
}

func getSession(ctx context.Context) *mgo.Session{
	return ctx.Value(mongoDBSessionKey).(*mgo.Session)
}


