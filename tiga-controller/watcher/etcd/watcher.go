package etcd

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/common/log"
    "go.etcd.io/etcd/clientv3"
    "strconv"
    "time"
)

type ClientSet struct {
    Client *clientv3.Client
    KV clientv3.KV
    Watch clientv3.Watcher
    respConfig struct {
        Pod map[string]int `json:"pod"`
        Node map[string]int `json:"node"`
    }
}

func (t *ClientSet)NewClient (etcdList []string) {
    etcdConfig := clientv3.Config{
        Endpoints:   etcdList,
        DialTimeout: time.Second * 1,
    }
    client, err := clientv3.New(etcdConfig)

    if err != nil {
        log.Fatal(err)
    }

    // etcdClient现在的版本就算是连不上也不会报错，可以通过下面的方法来判断，以后再实现吧，代码嘛，先能跑就行
    //timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 500)
    //defer cancel()
    //for _, v := range etcdConfig.Endpoints {
    //    _, err = client.Status(timeoutCtx, v)
    //    if err != nil {
    //        return
    //    }
    //}

    t.Client = client
    t.KV = clientv3.NewKV(client)
    t.Watch = clientv3.NewWatcher(client)
}

func (t *ClientSet)watcherFunc(watchKey string, metricsMap *map[string]int)  {
    tempMap := make(map[string]int)
    getResp, err := t.KV.Get(context.TODO(), watchKey)
    if err != nil {
        log.Fatal(watchKey, err)
    }
    watchStartRevision := getResp.Header.Revision + 1
    if len(getResp.Kvs) != 0 {
        log.Info(string(getResp.Kvs[0].Value), tempMap)

        err = json.Unmarshal(getResp.Kvs[0].Value, &tempMap)
        if err != nil {
            log.Error(watchKey, err)
            tempMap = make(map[string]int)
        }
        *metricsMap = tempMap
    } else {
        *metricsMap = tempMap
    }
    log.Info(watchKey, getResp.Header.Revision, *metricsMap)

    ctx, cancel := context.WithCancel(context.TODO())
    defer cancel()

    watchRespChan := t.Watch.Watch(ctx, watchKey , clientv3.WithRev(watchStartRevision))
    for watchResp := range watchRespChan {
        tempMap = make(map[string]int)
        for _, event := range watchResp.Events {
            switch event.Type.String() {
            case "PUT":
                err = json.Unmarshal(event.Kv.Value, &tempMap)
                if err != nil {
                    log.Error(string(event.Kv.Value), err)
                    tempMap = make(map[string]int)
                }
                *metricsMap = tempMap
            case "DELETE":
                fmt.Println("删除了", string(event.Kv.Value))
                *metricsMap = make(map[string]int)
            }
        }
    }
}

func (t *ClientSet)Watcher(port int)  {
    go func() {
        t.watcherFunc("/tiga-controller/daemonSet/pod", &t.respConfig.Pod)
    }()
    go func() {
       t.watcherFunc("/tiga-controller/daemonSet/node", &t.respConfig.Node)
    }()

    defer t.Client.Close()

    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.String(200, "以后再说")
    })
    r.GET("/daemonSet", func(c *gin.Context) {
        c.JSON(200, t.respConfig)
    })
    log.Fatal(r.Run(":" + strconv.Itoa(port)))
}

