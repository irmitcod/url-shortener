package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
	"url-shortener/mongo"
)

func InitMongoDatabase() mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := App.Config.MongodbUrl

	dbUser := App.Config.MongoDBUserName
	dbPass := App.Config.MongoPassword
	dbPort := App.Config.MongodbPort
	dbName := App.Config.MongoDB
	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	}

	opt := options.Client()
	opt.SetAuth(options.Credential{
		AuthMechanism:           os.Getenv("DATABASE_MONGO_AUTH_MECHANISM"),
		AuthMechanismProperties: nil,
		Username:                os.Getenv("DATABASE_MONGO_USER_NAME"),
		Password:                os.Getenv("DATABASE_MONGO_PASSWORD"),
		PasswordSet:             true,
	})

	opt.ApplyURI(os.Getenv("DATABASE_MONGO_HOST") + ":" + os.Getenv("DATABASE_MONGO_PORT") + "/" + dbName)

	//client, err := mongo.NewClient(opt)
	log.Println("start to connect mongo db")
	log.Println("start to connect mongo db")
	log.Println("start to connect mongo db")
	log.Println("start to connect mongo db")
	log.Println("start to connect mongo db")
	log.Println(mongodbURI)

	client, err := mongo.NewClient(opt)
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
