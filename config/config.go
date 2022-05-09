package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Host          string `env:"HOST" env-default:"0.0.0.0"`
	MongodbUrl    string `env:"MONGODB_URL" env-default:"mongodb://localhost:27017"`
	MongoPassword string `env:"MONGODB_PASSWORD" env-default:""`
	MongoDB       string `env:"MONGODB_DB" env-default:"url_shortener"`
	MongoTimeout  int    `env:"MONGODB_TIMEOUT" env-default:"30"`
	Port          int    `env:"PORT" env-default:"8080"`
	Prefix        string `env:"PREFIX" env-default:"/url-shortener/"`
	RedisAddress  string `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB       int    `env:"REDIS_DB" env-default:"0"`
	MaxWidth      int    `env:"MAX_WIDTH" env-default:"120"`
	MaxHeight     int    `env:"MAX_HEIGHT" env-default:"120"`
	UpperBound    int    `env:"UPPER_BOUND" env-default:"200"`
	LowerBound    int    `env:"LOWER_BOUND" env-default:"100"`
}

func GetConfig() *config {
	var cfg config
	cleanenv.ReadEnv(&cfg)
	return &cfg
}
