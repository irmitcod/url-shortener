package usecase

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"url-shortener/config"
	"url-shortener/src/models/url_shortener"
	"url-shortener/src/models/worker_result"
	redis2 "url-shortener/src/repository/redis"
	"url-shortener/src/utils/base58"
	"url-shortener/src/utils/lfu"
	"url-shortener/src/utils/rest_error"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type urlUsecase struct {
	urlRepo        url_shortener.UrlRepository
	contextTimeout time.Duration
	redisRepo      redis2.UrlRepository
	lfuCache       *lfu.Cache
	cache          config.LocalCache
}

func (uu *urlUsecase) GetUrl(url string, userid int) (buffer string, err rest_error.RestErr) {
	panic("implement me")
}

func (uu *urlUsecase) CacheUrl(url string) {
	panic("implement me")
}

func (uu *urlUsecase) GetInvalidUrl(url string) bool {
	ctx := context.Background()
	_, err := uu.redisRepo.GetInvalidUrl(ctx, url)
	if err != nil {
		return false
	}
	return true
}

func (i *urlUsecase) EncodUrl(UserID, url string) (buffer string, error rest_error.RestErr) {
	//get urlStruct form lfu cache
	key := base58.GenerateShortLink(url, UserID)
	url = fmt.Sprintf("%s::%s", url, UserID)
	cache := i.lfuCache.Get(key)
	if cache != nil {
		//i.entry.Infof("This  %s is already in lfu cache\n", url)
		buffer = cache.(string)
		return
	}

	//check url is invalid or not
	b := i.GetInvalidUrl(key)
	if b {
		error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		return
	}

	ctx := context.Background()
	bufferStr, err := i.redisRepo.GetUrl(ctx, key)
	if err == redis.Nil {

		//create chan for return date or error from download
		result := make(chan worker_result.Result)
		defer close(result)

		//add download and save to redis task with worker
		//cache urlStruct data with config and urlStruct format
		//we return from result chan to return as bytes for user requested
		go i.CacheUrlWithChan(url, (key), result)

		res := <-result
		if res.Status == 200 {
			buffer = res.Value
		} else {
			error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		}
		return
	}
	//create byte from redis

	buffer = (bufferStr)
	return
}
func (i *urlUsecase) CacheUrlWithChan(url string, file string, result chan worker_result.Result) {

	//buffer, _ := i.Translate(maxWidth, maxHeight, format)
	i.redisRepo.CacheUrl(url, file)
	i.lfuCache.Set(url, file)
	result <- worker_result.Result{
		Rrr:    nil,
		Status: 200,
		Value:  file,
	}
}

func NewUrlUsecase(u url_shortener.UrlRepository, to time.Duration, repository redis2.UrlRepository, cache config.LocalCache, l *lfu.Cache) url_shortener.UrlUsecase {
	return &urlUsecase{
		urlRepo:        u,
		contextTimeout: to,
		redisRepo:      repository,
		lfuCache:       l,
		cache:          cache,
	}
}

func (uu *urlUsecase) InsertOne(c context.Context, m *url_shortener.UrlShortener) (*url_shortener.UrlShortener, error) {

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	key := base58.GenerateShortLink(m.OriginalURL, m.UserID.String())

	cache := uu.lfuCache.Get(key)
	if cache != nil {
		//i.entry.Infof("This  %s is already in lfu cache\n", url)

		return &url_shortener.UrlShortener{
			ShortUrl:    key,
			OriginalURL: m.OriginalURL,
		}, nil
	}

	//check url is invalid or not
	b, _ := uu.redisRepo.GetUrl(ctx, key)
	if len(b) > 0 {
		return &url_shortener.UrlShortener{
			ShortUrl:    key,
			OriginalURL: m.OriginalURL,
		}, nil
	}

	shortner, err := uu.urlRepo.FindOneByKey(c, key)
	if err == nil {
		uu.lfuCache.Set(key, m.OriginalURL)
		uu.redisRepo.CacheUrl(key, m.OriginalURL)
		return shortner, nil
	}

	m.ID = primitive.NewObjectID()
	m.ShortUrl = key
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := uu.urlRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

	uu.lfuCache.Set(key, m.OriginalURL)
	uu.redisRepo.CacheUrl(key, m.OriginalURL)

	return res, nil
}

func (uu *urlUsecase) FindOne(c context.Context, id string) (*url_shortener.UrlShortener, error) {

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	res, err := uu.urlRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
func (uu *urlUsecase) FindOneByKey(c context.Context, id string) (string, error) {

	cache := uu.lfuCache.Get(id)
	if cache != nil {
		//i.entry.Infof("This  %s is already in lfu cache\n", url)
		return cache.(string), nil
	}

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	res, err := uu.urlRepo.FindOneByKey(ctx, id)
	if err != nil {
		return "", err
	}
	uu.lfuCache.Set(id, res.OriginalURL)
	uu.redisRepo.CacheUrl(id, res.OriginalURL)

	return res.OriginalURL, nil
}

func (uu *urlUsecase) GetAllWithPage(c context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]url_shortener.UrlShortener, int64, error) {

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	res, count, err := uu.urlRepo.GetAllWithPage(ctx, rp, p, filter, setsort)
	if err != nil {
		return res, count, err
	}

	return res, count, nil
}

func (uu *urlUsecase) UpdateOne(c context.Context, m *url_shortener.UrlShortener, id string) (*url_shortener.UrlShortener, error) {

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	res, err := uu.urlRepo.UpdateOne(ctx, m, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
