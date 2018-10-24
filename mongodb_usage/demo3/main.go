package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCon     *DeleteCond
		delResult  *mongo.DeleteResult
	)
	if client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017", clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}
	database = client.Database("cron")
	collection = database.Collection("job")
	delCon = &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	if delResult, err = collection.DeleteMany(context.TODO(), delCon); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("deleted:", delResult.DeletedCount)
}
