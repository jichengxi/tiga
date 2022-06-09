package test

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRedisCli(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		//Addr: "120.24.98.32:26379",
		Addr:        "192.168.1.21:6379",
		DialTimeout: time.Second * 3,
		ReadTimeout: time.Second * 3,
	})

	_, err := client.Ping().Result()
	//t.Log(result, err)
	if err != nil {
		panic(err)
	}
	result, err := client.HGetAll("S").Result()
	if err != nil {
		panic(err)
	}
	t.Log(result)
}
