package mongo

import (
	"context"
	"time"
	"url-shortener/mongo"
	"url-shortener/src/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "user"
)

func NewMongoRepository(DB mongo.Database) users.UserRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, user *users.User) (*users.User, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string) (*users.User, error) {
	var (
		user users.User
		err  error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *mongoRepository) UpdateOne(ctx context.Context, user *users.User, id string) (*users.User, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": bson.M{
		"name":       user.Name,
		"username":   user.Username,
		"password":   user.Password,
		"updated_at": time.Now(),
	}}

	_, err = m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *mongoRepository) GetByCredential(ctx context.Context, username string, password string) (*users.User, error) {
	var (
		user users.User
		err  error
	)

	credential := bson.M{
		"username": username,
		"password": password,
	}

	err = m.Collection.FindOne(ctx, credential).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
