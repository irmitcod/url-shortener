package url_shortener

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"url-shortener/src/models/worker_result"
	"url-shortener/src/utils/rest_error"
)

type UrlShortener struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	ShortUrl    string             `bson:"key" json:"key" `
	OriginalURL string             `bson:"original_url" json:"original_url" validate:"required"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Hits        int                `bson:"hits" json:"hits"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	ExpireAt    time.Time          `bson:"expire_at" json:"expire_at"`
}

type UrlRepository interface {
	InsertOne(ctx context.Context, u *UrlShortener) (*UrlShortener, error)
	FindOne(ctx context.Context, id string) (*UrlShortener, error)
	FindByHits(ctx context.Context) ([]UrlShortener, error)
	FindOneByKey(ctx context.Context, id string) (*UrlShortener, error)
	UpdateOne(ctx context.Context, cat *UrlShortener, id string) (*UrlShortener, error)
}

type UrlUsecase interface {
	InsertOne(ctx context.Context, u *UrlShortener) (*UrlShortener, error)
	FindOne(ctx context.Context, id string) (*UrlShortener, error)
	FindOneByKey(ctx context.Context, id string) (string, error)
	UpdateOne(ctx context.Context, cat *UrlShortener, id string) (*UrlShortener, error)
	CacheUrl(url string)
	GetInvalidUrl(url string) bool
	GetUrl(url string, userid int) (buffer string, err rest_error.RestErr)
	CacheUrlWithChan(url string, file string, result chan worker_result.Result)
	EncodUrl(UserID, url string) (buffer string, error rest_error.RestErr)
}
