package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/prometheus/common/log"
	"strings"
	"tiga-controller/ControllerTypes"
	"time"
)

type RdClient struct {
	*redis.Client
}

func NewAloneClient(ip string, db int) RdClient {
	client := redis.NewClient(&redis.Options{
		Addr:        ip,
		DB:          db,
		DialTimeout: time.Second * 3,
		ReadTimeout: time.Second * 3,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return RdClient{client}
}

func NewSentinelClient(ip string, master string, password string, db int) RdClient {
	ipList := strings.Split(ip, ",")
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: ipList,
		MasterName:    master,
		Password:      password,
		DB:            db,
		DialTimeout:   time.Second * 3,
		ReadTimeout:   time.Second * 3,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return RdClient{client}
}

func (t *RdClient) watcherFuc(hKey string, fieldKey string, metricsConfig interface{}) {
	if hKey == "local" {
		if fieldKey == "container" {
			config := metricsConfig.(*ControllerTypes.LocalMetrics)
			for {
				var localMetrics ControllerTypes.LocalMetrics
				containerStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(containerStr), &localMetrics)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = localMetrics
				time.Sleep(time.Second * 30)
			}

		} else if fieldKey == "node" {
			config := metricsConfig.(*ControllerTypes.LocalMetrics)
			for {
				var localMetrics ControllerTypes.LocalMetrics
				containerStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(containerStr), &localMetrics)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = localMetrics
				time.Sleep(time.Second * 30)
			}
		} else {
			return
		}

	} else if hKey == "remote" {
		if fieldKey == "enable" {
			config := metricsConfig.(*bool)
			for {
				var enableConfig bool
				urlStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(urlStr), &enableConfig)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = enableConfig
				time.Sleep(time.Second * 30)
			}
		} else if fieldKey == "url" {
			config := metricsConfig.(*map[string]string)
			for {
				var urlConfig map[string]string
				urlStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(urlStr), &urlConfig)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = urlConfig
				time.Sleep(time.Second * 30)
			}
		} else if fieldKey == "metrics" {
			config := metricsConfig.(*[]string)
			for {
				var metrics []string
				urlStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(urlStr), &metrics)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = metrics
				time.Sleep(time.Second * 30)
			}
		} else if fieldKey == "blacklist" {
			config := metricsConfig.(*[]map[string]string)
			for {
				var blacklistConfig []map[string]string
				urlStr, err := t.HGet(hKey, fieldKey).Result()
				if err != nil {
					if err != redis.Nil {
						log.With("key", hKey).Error(err.Error())
					}
				} else {
					err = json.Unmarshal([]byte(urlStr), &blacklistConfig)
					if err != nil {
						log.With("key", hKey).Error(err.Error())
					}
				}
				*config = blacklistConfig
				time.Sleep(time.Second * 30)
			}
		} else {
			return
		}
	} else {
		log.With("key", hKey).Error("redis HKey is not local or remote.")
		return
	}
}

func (t *RdClient) Watcher() {
	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("local", "container", &ControllerTypes.ResultData.Local.Container)
		defer ControllerTypes.WorkGroup.Done()
	}()

	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("local", "node", &ControllerTypes.ResultData.Local.Node)
		defer ControllerTypes.WorkGroup.Done()
	}()

	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("remote", "enable", &ControllerTypes.ResultData.Remote.Enable)
		defer ControllerTypes.WorkGroup.Done()
	}()

	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("remote", "url", &ControllerTypes.ResultData.Remote.Url)
		defer ControllerTypes.WorkGroup.Done()
	}()

	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("remote", "metrics", &ControllerTypes.ResultData.Remote.Metrics)
		defer ControllerTypes.WorkGroup.Done()
	}()

	go func() {
		ControllerTypes.WorkGroup.Add(1)
		t.watcherFuc("remote", "blacklist", &ControllerTypes.ResultData.Remote.Blacklist)
		defer ControllerTypes.WorkGroup.Done()
	}()
}
