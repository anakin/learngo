package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		kv             clientv3.KV
		ctx            context.Context
		cancel         context.CancelFunc
		getResp        *clientv3.GetResponse
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	putResp, err = kv.Put(ctx, "test", "1111")
	cancel()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("write success", putResp.Header.Revision)

	lease = clientv3.NewLease(client)
	//fmt.Println(lease)
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("grant fail,", err)
		return
	} else {
		fmt.Println("grant success")
	}
	fmt.Println("hereee")

	leaseId = leaseGrantResp.ID
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功", putResp.Header.Revision)
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("过期了")
			break
		}
		fmt.Println("没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
