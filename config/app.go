package config

import (
	"database/sql"
	"url-shortener/mongo"
)

var (
	App *Application
)

type Application struct {
	Config *Config
	MySql  *sql.DB
	Mongo  mongo.Client
}

func init() {
	AppInit()
}

func AppInit() {
	App = &Application{}
	App.Config = GetConfig()
	App.Mongo = InitMongoDatabase()
}
