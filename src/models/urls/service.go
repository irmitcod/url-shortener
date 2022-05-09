package urls

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
	"url-shortener/config"
	"url-shortener/src/models/worker_result"
	"url-shortener/src/repository"
	"url-shortener/src/utils/base58"
	"url-shortener/src/utils/lfu"
	"url-shortener/src/utils/rest_error"
	"url-shortener/src/utils/workerpool"
)

var (
	MaxRequest = 250
)

var (
	ErrNon200      = errors.New("received non 200 response code")
	ErrURLNotFound = errors.New("urlStruct not found")
	requestCounter = 0
)

type Service interface {
	CacheUrl2(url string, file []byte)
	CacheUrl(url string)
	GetInvalidUrl(url string) bool
	GetUrl(url string) (buffer []byte, err rest_error.RestErr)
	EncodUrl(url string) rest_error.RestErr
}
type urlStruct struct {
	cache               config.LocalCache
	lfuCache            *lfu.Cache
	repository          repository.UrlRepository
	wp                  workerpool.WorkerPool
	MaxWidth, MaxHeight int
	entry               *log.Entry
}

func (i urlStruct) GetInvalidUrl(url string) bool {
	ctx := context.Background()
	_, err := i.repository.GetInvalidUrl(ctx, url)
	if err == redis.Nil {
		return false
	}
	return true
}

func (i urlStruct) CacheUrl(url string) {
	i.repository.CacheInvalidUrl(url)
}

func (i urlStruct) EncodUrl(url string) (error rest_error.RestErr) {
	//get urlStruct form lfu cache
	cache := i.lfuCache.Get(url)
	if cache != nil {
		i.entry.Infof("This  %s is already in lfu cache\n", url)
		return
	}

	//check url is invalid or not
	b := i.GetInvalidUrl(url)
	if b {
		error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		return
	}

	ctx := context.Background() //todo add timeout
	_, err := i.repository.GetUrl(ctx, url)

	if err == redis.Nil {
		//create name spaces for local cache for is urk in download progress or not
		_, ok := i.cache.Get(url)
		if ok {
			i.entry.Infof("This  %s in download proggress\n", url)
			return nil
		}

		//set url to cache for check is url in download
		//progress and after 30 second is going to remove url form local cache
		i.cache.SetWithTTL(url, 1, 0, 30*time.Second)

		//add download and save to redis task with worker
		//create chan for return date or error from download
		result := make(chan worker_result.Result)
		defer close(result)

		//add download and save to redis task with worker
		i.wp.AddTask(func() {
			encoded := ""
			go i.CacheUrlWithChan(url, []byte(encoded), result)
		})
		res := <-result
		if res.Status == 404 {
			error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		}
		return
	}
	return nil
}

func (i urlStruct) evictUrl() {
	if requestCounter == MaxRequest {
		i.lfuCache.Evict(1)
		requestCounter = 0
	}
	requestCounter += 1
}

func (i urlStruct) GetUrl(url string) (buffer []byte, error rest_error.RestErr) {
	//get urlStruct form lfu cache
	cache := i.lfuCache.Get(url)
	if cache != nil {
		i.entry.Infof("This  %s is already in lfu cache\n", url)
		buffer = cache.([]byte)
		return
	}

	//check url is invalid or not
	b := i.GetInvalidUrl(url)
	if b {
		error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		return
	}

	ctx := context.Background()
	bufferStr, err := i.repository.GetUrl(ctx, url)
	if err == redis.Nil {
		//create name spaces for local cache for is urk in download progress or not
		_, ok := i.cache.Get(url)
		if ok {
			i.entry.Infof("This  %s in download proggress\n", url)
			return nil, nil
		}
		//set url to cache for check is url in download
		//progress and after 30 second is going to remove url form local cache
		i.cache.SetWithTTL(url, 1, 0, 30*time.Second)

		//create chan for return date or error from download
		result := make(chan worker_result.Result)
		defer close(result)

		//add download and save to redis task with worker
		i.wp.AddTask(func() {
			response := base58.GenerateShortLink(url, "200")

			//cache urlStruct data with config and urlStruct format
			//we return from result chan to return as bytes for user requested
			go i.CacheUrlWithChan(url, []byte(response), result)
		})
		res := <-result
		if res.Status == 200 {
			buffer = res.Value
		} else {
			error = rest_error.NewNotFoundError("this url is not valid and couldn't fine any urlStruct from this url")
		}
		return
	}
	//create byte from redis
	i.entry.Infof("This  %s is already in redis cache\n", url)
	buffer = []byte(bufferStr)
	return
}

func (i urlStruct) CacheUrl2(url string, file []byte) {

	//buffer, _ := i.Translate(maxWidth, maxHeight, format)
	i.lfuCache.Set(url, file)
	i.repository.CacheUrl(url, file)
}
func (i urlStruct) CacheUrlWithChan(url string, file []byte, result chan worker_result.Result) {

	//buffer, _ := i.Translate(maxWidth, maxHeight, format)
	i.repository.CacheUrl(url, file)
	i.lfuCache.Set(url, file)
	result <- worker_result.Result{
		Rrr:    nil,
		Status: 0,
		Value:  file,
	}
}
func NewService(repo *repository.UrlRepository, wp workerpool.WorkerPool, maxWidth, maxHeight int, cache config.LocalCache, lf *lfu.Cache, entry *log.Entry) Service {
	return &urlStruct{
		lfuCache:   lf,
		cache:      cache,
		wp:         wp,
		repository: *repo,
		MaxWidth:   maxWidth,
		MaxHeight:  maxHeight,
		entry:      entry,
	}
}
