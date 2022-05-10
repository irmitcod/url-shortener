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
	if err == redis.Nil {
		return false
	}
	return true
}

func (i *urlUsecase) EncodUrl(UserID, url string) (buffer string, error rest_error.RestErr) {
	//get urlStruct form lfu cache

	url = fmt.Sprintf("%s::%s", url, UserID)
	cache := i.lfuCache.Get(url)
	if cache != nil {
		//i.entry.Infof("This  %s is already in lfu cache\n", url)
		buffer = cache.(string)
		return
	}

	//check url is invalid or not
	b := i.GetInvalidUrl(url)
	if b {
		error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		return
	}

	ctx := context.Background()
	bufferStr, err := i.redisRepo.GetUrl(ctx, url)
	if err == redis.Nil {

		//create chan for return date or error from download
		result := make(chan worker_result.Result)
		defer close(result)

		//add download and save to redis task with worker

		response := base58.GenerateShortLink(url, UserID)

		//cache urlStruct data with config and urlStruct format
		//we return from result chan to return as bytes for user requested
		go i.CacheUrlWithChan(url, (response), result)

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
	encodeURl, errREs := uu.EncodUrl(m.UserID.String(), m.OriginalURL)
	if errREs != nil {
		return nil, errREs
	}
	m.ID = primitive.NewObjectID()
	m.Key = encodeURl
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := uu.urlRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

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
func (uu *urlUsecase) FindOneByKey(c context.Context, id string) (*url_shortener.UrlShortener, error) {

	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	res, err := uu.urlRepo.FindOneByKey(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
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
