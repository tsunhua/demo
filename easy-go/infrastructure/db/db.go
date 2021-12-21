package db

import (
	"app/infrastructure/config"
	"app/infrastructure/log"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var dbConfig = config.Get().DB

var mongoInstance = struct {
	client *mongo.Client
	err    error
	sync.Mutex
}{}

func init() {
	go initMongo()
}

func initMongo() {
	if mongoInstance.client == nil {
		mongoInstance.Lock()
		defer mongoInstance.Unlock()
		if mongoInstance.client == nil {
			log.Info("start connect mongo")
			ctx, _ := context.WithTimeout(context.Background(), dbConfig.ConnectionTimeout)
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.URI))
			if err != nil {
				mongoInstance.err = err
				mongoInstance.client = nil
			} else {
				mongoInstance.client = client
			}
			log.Info(fmt.Sprintf("finish connect mongo, err:%v", mongoInstance.err))
			return
		}
	}
}

func mongoClient() (*mongo.Client, error) {
	initMongo()
	return mongoInstance.client, mongoInstance.err
}

func getTable(name string) (*mongo.Collection, error) {
	client, err := mongoClient()
	if err != nil {
		return nil, err
	}
	return client.Database(dbConfig.Name).Collection(name), nil
}
