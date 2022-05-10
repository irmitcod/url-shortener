package config

import (
	"context"
	"log"
	"time"
	"url-shortener/mongo"
)

func InitMongoDatabase() mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := App.Config.MongodbUrl

	//dbUser :="users"
	//dbPass := App.Config.GetString(`mongo.pass`)
	//dbName := App.Config.GetString(`mongo.name`)
	//mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	//if dbUser == "" || dbPass == "" {
	//	mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	//}

	client, err := mongo.NewClient(dbHost)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
