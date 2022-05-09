package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"url-shortener/config"
	"url-shortener/src/controller"
	"url-shortener/src/models/urls"
	"url-shortener/src/repository/mongo"
	"url-shortener/src/utils/lfu"
	"url-shortener/src/utils/workerpool"
)

// @title    AISA API
// @version  1.0
func main() {

	//config logger
	entry := config.NewLogger()

	// Setup Configuration
	configuration := config.GetConfig()
	database := config.NewMemoryClient(configuration)

	mongoClient, err := config.NewMongoClient(configuration.MongodbUrl, configuration.MongoDB, configuration.MongoTimeout)
	if err != nil {
		panic(err)
	}
	// Setup Repository
	productRepository := mongo.NewUrlRepository(database, mongoClient)

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
	imageService := urls.NewService(&productRepository, wp, configuration.MaxHeight, configuration.MaxWidth, cache, lfuCach, entry)

	//setup controller
	handler := controller.NewHandler(imageService)

	// Start server
	router := gin.Default()
	router.Use(cors.Default())
	prefix := router.Group(configuration.Prefix)
	handler.Route(prefix)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	serverAddr := fmt.Sprintf("%s:%d", configuration.Host, configuration.Port)
	log.Panic(router.Run(serverAddr))
}
