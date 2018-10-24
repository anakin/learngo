package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config          clientv3.Config
		client          *clientv3.Client
		err             error
		getResp         *clientv3.GetResponse
		kv              clientv3.KV
		watcher         clientv3.Watcher
		startRevision   int64
		watcherRespChan <-chan clientv3.WatchResponse
		watchResp       clientv3.WatchResponse
		event           *clientv3.Event
	)

	config = clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("connected")
	}
	kv = clientv3.NewKV(client)
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}
	}()
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println("current:", string(getResp.Kvs[0].Value))
	}
	startRevision = getResp.Header.Revision + 1
	watcher = clientv3.Watcher(client)
	fmt.Println("start revision:", startRevision)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watcherRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(startRevision))
	for watchResp = range watcherRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("modify to:", string(event.Kv.Value), "rev:", event.Kv.ModRevision)

			case mvccpb.DELETE:
				fmt.Println("delete rev:", event.Kv.ModRevision)
			}

		}
	}
}
