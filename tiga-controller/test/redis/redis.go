package main

import (
	"fmt"
	"tiga-controller/watcher/redis"
)

func main() {
	//data := ControllerTypes.RespConfig{
	//	Local: ControllerTypes.LocalConfig{
	//		Container: ControllerTypes.LocalMetrics{
	//			Enable: true,
	//			Metrics: []string{
	//				"local_container_cpu_system_seconds_total",
	//				"local_container_memory_usage_bytes",
	//			},
	//		},
	//		Node: ControllerTypes.LocalMetrics{
	//			Enable:  true,
	//			Metrics: []string{},
	//		},
	//	},
	//	Remote: ControllerTypes.RemoteConfig{
	//		Url: map[string]string{
	//			{"cadvisor": "http://172.17.36.201/metrics/cadvisor.txt"},
	//			{"node": "http://172.17.36.201/metrics/node.txt"},
	//		},
	//		Metrics: []string{
	//			"cadvisor_cadvisor_version_info",
	//			"cadvisor_container_cpu_load_average_10s",
	//		},
	//	},
	//}
	//local := ControllerTypes.LocalConfig{
	//	Container: ControllerTypes.LocalMetrics{
	//		Enable: true,
	//		Metrics: []string{
	//			"local_container_cpu_system_seconds_total",
	//			"local_container_memory_usage_bytes",
	//		},
	//	},
	//	Node: ControllerTypes.LocalMetrics{
	//		Enable:  true,
	//		Metrics: []string{},
	//	},
	//}
	//data := make(map[string]string)
	//data["cadvisor"] = "http://172.17.36.201/metrics/cadvisor.txt"
	//data["node"] = "http://172.17.36.201/metrics/node.txt"
	dataStr := fmt.Sprintf("%t", true)
	//dataStr, err := json.Marshal(data)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(string(dataStr))
	redisCli := redis.NewAloneClient("192.168.1.21:6379", 0)
	err := redisCli.HSet("remote", "enable", dataStr).Err()
	if err != nil {
		panic(err)
	}

	//var a ControllerTypes.LocalMetrics
	//_, err := redisCli.HGet("local", "container").Result()
	//if err != nil {
	//	return
	//}
	//fmt.Println(a)
}
