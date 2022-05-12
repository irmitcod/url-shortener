package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	ApiKeys   struct {
		Key      string `json:"key" bson:"key"`
		ExpireAt int64  `json:"expire_at" bson:"expire_at"`
		CreateAt int64  `json:"create_at" bson:"create_at"`
	} `json:"api_keys" bson:"api_keys"`
}
type UserRepository interface {
	InsertOne(ctx context.Context, u *User) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	UpdateOne(ctx context.Context, user *User, id string) (*User, error)
	GetByCredential(ctx context.Context, username string, password string) (*User, error)
}

type UserUsecase interface {
	InsertOne(ctx context.Context, u *User) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	UpdateOne(ctx context.Context, user *User, id string) (*User, error)
}
