package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}
type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command""`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint""`
}

func main() {

	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		result     *mongo.InsertOneResult
		docId      objectid.ObjectID
	)
	if client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017", clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}
	database = client.Database("cron")
	collection = database.Collection("job")
	record = &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}

	docId = result.InsertedID.(objectid.ObjectID)
	fmt.Println("insert ID:", docId.Hex())
}
