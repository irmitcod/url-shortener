package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host            string `env:"HOST" env-default:"0.0.0.0"`
	MongodbUrl      string `env:"DATABASE_MONGO_HOST" env-default:"localhost"`
	MongodbPort     string `env:"DATABASE_MONGO_PORT" env-default:"27017"`
	BaseUrl         string `env:"BASE_URL" env-default:"http://localhost:8080/"`
	MongoPassword   string `env:"MONGO_INITDB_ROOT_PASSWORD" env-default:"admin"`
	MongoDB         string `env:"MONGO_INITDB_DATABASE" env-default:"url_shortener"`
	MongoDBUserName string `env:"MONGO_INITDB_ROOT_USERNAME" env-default:"admin"`
	MongoTimeout    int    `env:"MONGODB_TIMEOUT" env-default:"30"`
	Port            int    `env:"PORT" env-default:"8080"`
	Prefix          string `env:"PREFIX" env-default:"/url-shortener/"`
	RedisAddress    string `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	RedisPassword   string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB         int    `env:"REDIS_DB" env-default:"0"`
	MaxWidth        int    `env:"MAX_WIDTH" env-default:"120"`
	Secret          string `env:"SECRET" env-default:"aslkjda;lkshdaslk;dhaskljd"`
	MaxHeight       int    `env:"MAX_HEIGHT" env-default:"120"`
	UpperBound      int    `env:"UPPER_BOUND" env-default:"200"`
	LowerBound      int    `env:"LOWER_BOUND" env-default:"100"`
	Lifetime        string `env:"LOWER_BOUND" env-default:"100"`
	Username        string `env:"ADMIN_PUSERNAME" env-default:"master"`
	Password        string `env:"ADMIN_PASSWORd" env-default:"123456789"`
}

func GetConfig() *Config {
	var cfg Config
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
