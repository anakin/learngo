package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr          *cronexpr.Expression
		err           error
		now, nextTime time.Time
	)
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
	}
	now = time.Now()
	nextTime = expr.Next(now)
	//fmt.Println(nextTime)

	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("here", now)
	})

	time.Sleep(20 * time.Second)
}
