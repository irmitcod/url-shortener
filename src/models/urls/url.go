package urls

import (
	log "github.com/sirupsen/logrus"
	"url-shortener/config"
	redis2 "url-shortener/src/repository/redis"
	"url-shortener/src/utils/lfu"
	"url-shortener/src/utils/workerpool"
)

type UrlStruct struct {
	cache               config.LocalCache
	lfuCache            *lfu.Cache
	repository          redis2.UrlRepository
	wp                  workerpool.WorkerPool
	MaxWidth, MaxHeight int
	entry               *log.Entry

	UserId      int64  `json:"user_id" bson:"user_id"`
	OriginalURL string `json:"original_url" bson:"original_url"`
	EncodedURL  string `json:"encoded_url" bson:"encoded_url"`
	CreatedAt   int64  `bson:"created_at" msgpack:"created_at" json:"created_at"`
	ExpireAt    int64  `bson:"expire_at" msgpack:"expire_at" json:"expire_at"`
}
