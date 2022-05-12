package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"time"
	"url-shortener/config"
	"url-shortener/src/controller"
	_jwt "url-shortener/src/models/jwt/usecase"
	_loginUsecase "url-shortener/src/models/login/usecase"
	_urlusecase "url-shortener/src/models/url_shortener/usecase"
	_usecase "url-shortener/src/models/users/usecase"
	"url-shortener/src/repository/redis"
	_urlShortenerRepo "url-shortener/src/repository/url_shortner/mongo"
	"url-shortener/src/repository/users/mongo"
	"url-shortener/src/utils/lfu"
	"url-shortener/src/utils/workerpool"
)

// @title    AISA API
// @version  1.0
func main() {
	e := echo.New()

	// Setup Configuration
	configuration := config.GetConfig()

	// setup lfu cache
	lfuCach := lfu.New()
	//set upper bound for lfu cache if len oc caceh reached to the UpperBouund
	//evict elements
	lfuCach.UpperBound = configuration.UpperBound
	lfuCach.LowerBound = configuration.LowerBound

	// setup worker pool
	totalWorker := 5
	wp := workerpool.NewWorkerPool(totalWorker)
	wp.Run()

	//setup localache
	cache := config.NewCache()
	// Setup Service

	//time out for connec
	timeoutContext := time.Duration(configuration.MongoTimeout) * time.Second
	database := config.NewMemoryClient(configuration)
	redisRepository := redis.NewUrlRepository(database)

	//get monogdb
	mongoDatabase := config.App.Mongo.Database(configuration.MongoDB)
	//this repository for mogno repository
	userRepo := mongo.NewMongoRepository(mongoDatabase)
	usrUsecase := _usecase.NewUserUsecase(userRepo, timeoutContext, redisRepository, cache)

	//set jwt middleware
	jwt := _jwt.NewJwtUsecase(userRepo, 4380*time.Hour, configuration)
	userJwt := e.Group("")

	//set api key middleware for user
	userJwt.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
			return true, nil
		},
	}))
	jwt.SetJwtUser(userJwt)
	adminJwt := e.Group("")
	jwt.SetJwtAdmin(adminJwt)
	generalJwt := e.Group("")
	jwt.SetJwtGeneral(generalJwt)

	//init url shortner repository
	urlShortenerRepo := _urlShortenerRepo.NewMongoRepository(mongoDatabase)

	//Handle For urls shortener endpoint
	urlUsecase := _urlusecase.NewUrlUsecase(urlShortenerRepo, timeoutContext, redisRepository, cache, lfuCach)
	controller.NewUserHandler(e, userJwt, usrUsecase, urlUsecase, configuration)

	//Handle For login endpoint
	loginUsecase := _loginUsecase.NewLoginUsecase(userRepo, timeoutContext)
	controller.NewLoginHandler(e, loginUsecase, configuration)

	//Handle For urls shortener endpoint
	controller.NewUrlHandler(userJwt, urlUsecase, database, configuration)

	appPort := fmt.Sprintf(":%d", configuration.Port)
	log.Fatal(e.Start(appPort))

}
