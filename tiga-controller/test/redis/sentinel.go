package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"tiga-controller/ControllerTypes"
)

func main() {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"172.22.32.98:26380", "172.22.32.99:26380", "172.22.32.100:26380"},
		MasterName:    "harbor",
		Password:      "Redhat@2016",
		DB:            8,
	})
	data := ControllerTypes.LocalMetrics{
		Enable: true,
	}
	dataStr, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, err = client.HSet("local", "container", dataStr).Result()
	_, err = client.HSet("local", "node", dataStr).Result()
	if err != nil {
		panic(err)
	}
}
