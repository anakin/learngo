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
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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

	lease = clientv3.NewLease(client)
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	leaseId = leaseGrantResp.ID

	ctx, cancelFunc = context.WithCancel(context.TODO())

	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("keep alive fail")
					goto END
				} else {
					fmt.Println("receive:", keepResp.ID)
				}
			}
		}
	END:
	}()

	kv = clientv3.NewKV(client)
	txn = kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/jobs/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/jobs/job9"))
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	if !txnResp.Succeeded {
		fmt.Println("occupied", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}
	fmt.Println("doing sth")
	time.Sleep(5 * time.Second)
}
