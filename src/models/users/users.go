package users

type User struct {
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	LastLogin string `json:"last_login" bson:"lastLogin"`
	ApiKeys   struct {
		Key      string `json:"key" bson:"key"`
		ExpireAt int64  `json:"expire_at" bson:"expire_at"`
		CreateAt int64  `json:"create_at" bson:"create_at"`
	}
}
