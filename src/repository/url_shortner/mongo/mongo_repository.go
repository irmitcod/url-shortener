package mongo

import (
	"context"
	"time"
	"url-shortener/mongo"
	"url-shortener/src/models/url_shortener"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

func (m *mongoRepository) FindByHits(ctx context.Context) ([]url_shortener.UrlShortener, error) {

	c, err := m.Collection.Find(ctx, bson.M{"hits": bson.M{"$gte": 1000}})
	if err != nil {
		return nil, err
	}
	var urls url_shortener.UrlShortener
	err = c.All(ctx, &urls)
	if err != nil {
		return nil, err
	}
	return nil, err
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
		"key":        urlShortener.ShortUrl,
		"value":      urlShortener.OriginalURL,
		"hits":       urlShortener.Hits,
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
