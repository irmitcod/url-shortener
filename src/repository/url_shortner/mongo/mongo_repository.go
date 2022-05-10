package mongo

import (
	"context"
	"fmt"
	"time"
	"url-shortener/mongo"
	"url-shortener/src/models/url_shortener"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "urlShortener"
)

func NewMongoRepository(DB mongo.Database) url_shortener.UrlRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, urlShortener *url_shortener.UrlShortener) (*url_shortener.UrlShortener, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, urlShortener)
	if err != nil {
		return urlShortener, err
	}

	return urlShortener, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string) (*url_shortener.UrlShortener, error) {
	var (
		urlShortener url_shortener.UrlShortener
		err          error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &urlShortener, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&urlShortener)
	if err != nil {
		return &urlShortener, err
	}

	return &urlShortener, nil
}

func (m *mongoRepository) FindOneByKey(ctx context.Context, id string) (*url_shortener.UrlShortener, error) {
	var (
		urlShortener url_shortener.UrlShortener
		err          error
	)

	err = m.Collection.FindOne(ctx, bson.M{"key": id}).Decode(&urlShortener)
	if err != nil {
		return &urlShortener, err
	}

	return &urlShortener, nil
}

func (m *mongoRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]url_shortener.UrlShortener, int64, error) {

	var (
		urlShortener []url_shortener.UrlShortener
		skip         int64
		opts         *options.FindOptions
	)

	skip = (p * rp) - rp
	if setsort != nil {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
			options.Find().SetSort(setsort),
		)
	} else {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
		)
	}

	cursor, err := m.Collection.Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, err
	}
	if cursor == nil {
		return nil, 0, fmt.Errorf("nil cursor value")
	}
	err = cursor.All(ctx, &urlShortener)
	if err != nil {
		return nil, 0, err
	}

	count, err := m.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return urlShortener, 0, err
	}

	return urlShortener, count, err
}

func (m *mongoRepository) UpdateOne(ctx context.Context, urlShortener *url_shortener.UrlShortener, id string) (*url_shortener.UrlShortener, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return urlShortener, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": bson.M{
		"key":        urlShortener.Key,
		"value":      urlShortener.OriginalURL,
		"updated_at": time.Now(),
	}}

	_, err = m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return urlShortener, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(urlShortener)
	if err != nil {
		return urlShortener, err
	}
	return urlShortener, nil
}
