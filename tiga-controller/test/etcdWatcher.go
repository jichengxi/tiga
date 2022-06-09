package test

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "log"
    "time"
)

func EtcdClient() {
    var (
        client *clientv3.Client
        err error
        kv clientv3.KV
        watcher clientv3.Watcher
        getResp *clientv3.GetResponse
        watchStartRevision int64
        watchRespChan <-chan clientv3.WatchResponse
        watchResp clientv3.WatchResponse
        event *clientv3.Event
    )


    client, err = clientv3.New(clientv3.Config{
        Endpoints:   []string{"120.24.98.32:2379"},
        DialTimeout: time.Second * 1,
    })

    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    kv = clientv3.NewKV(client)

    // 模拟etcd中KV的变化
    //go func() {
    //    for {
    //        kv.Put(context.TODO(), "/test01", "i am test01")
    //
    //        kv.Delete(context.TODO(), "/test01")
    //
    //        time.Sleep(1 * time.Second)
    //    }
    //}()

    if getResp, err = kv.Get(context.TODO(), "/test01"); err != nil {
        fmt.Println("有问题")
        log.Fatal(err)
    }


    if len(getResp.Kvs) != 0 {
        fmt.Println("当前值:", string(getResp.Kvs[0].Value))
    } else {
        fmt.Println("当前没有值")
    }

    // 当前etcd集群事务ID, 单调递增的（监听/cron/jobs/job7后续的变化,也就是通过监听版本变化）
    fmt.Println("事务ID：", getResp.Header.Revision)
    watchStartRevision = getResp.Header.Revision + 1

    // 创建一个watcher(监听器)
    watcher = clientv3.NewWatcher(client)

    // 启动监听
    //fmt.Println("从该版本向后监听:", watchStartRevision)

    ctx, cancel := context.WithCancel(context.TODO())
    //5秒钟后取消
    //time.AfterFunc(10 * time.Second, func() {
    //    cancel()
    //})
    defer cancel()

    //这里ctx感知到cancel则会关闭watcher
    watchRespChan = watcher.Watch(ctx, "/test01" , clientv3.WithRev(watchStartRevision))

    // 处理kv变化事件
    for watchResp = range watchRespChan {
        for _, event = range watchResp.Events {
            log.Println(event.Type)
            switch event.Type.String() {
            case "PUT":
                fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
            case "DELETE":
                fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
            }
        }
    }



}

