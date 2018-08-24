package model

var (
	UserType    = EntityType{"user", "u"}
	AuthType    = EntityType{"authdata", "auth"}
	MessageType = EntityType{"message", "m"}

	FriendsRelation = Relationship{UserType, "friends"}
)

type EntityType struct {
	Name      string //for collection name
	ShortName string //for short collection names, such as redis key
}

type Relationship struct {
	fromEntity   EntityType
	relationName string
}

func (r *Relationship) GenID(fromId ID) string {
	return r.fromEntity.ShortName + ":" + fromId + ":" + r.relationName
}

type ID = string

type IDWithCursor struct {
	Id     ID
	Cursor string
}

type User struct {
	Id     ID `bson:"_id"`
	Name   string
	Gender string
}

type Auth struct {
	LoginName      string `bson:"_id"`
	UserId         string
	HashedPassword string
	IsAdmin        bool
	Permissions    []string
}

type Message struct {
	Id       ID `bson:"_id"`
	Text     string
	AuthorID ID
}
