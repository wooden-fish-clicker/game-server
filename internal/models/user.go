package models

type User struct {
	ID         string   `bson:"_id,omitempty"`
	Account    string   `bson:"account"`
	Email      string   `bson:"email"`
	Password   string   `bson:"password"`
	UserInfo   UserInfo `bson:"user_info"`
	CreatedAt  int64    `bson:"created_at"`
	ModifiedAt int64    `bson:"modified_at"`
}

type UserInfo struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Points  int64  `bson:"points"`
	Hp      int32  `bson:"hp"`
}
