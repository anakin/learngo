package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
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
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		cond       *FindByJobName
		cursor     mongo.Cursor
		record     *LogRecord
	)
	if client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017", clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}
	database = client.Database("cron")
	collection = database.Collection("job")

	cond = &FindByJobName{JobName: "job10"}

	if cursor, err = collection.Find(context.TODO(), cond, findopt.Skip(0), findopt.Limit(2)); err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}
}
