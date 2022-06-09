package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type Collector struct {
	gaugeDesc *prometheus.Desc
	mutex     sync.Mutex
}

func NewCollector() prometheus.Collector {
	return &Collector{
		gaugeDesc: prometheus.NewDesc(
			"container_memory_max_usage_bytes",
			"Maximum memory usage recorded in bytes",
			[]string{"container_label_pod_ip"},
			nil,
		),
	}
}

//实现采集器Describe接口
func (n *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.gaugeDesc
}

//实现采集器Collect接口,真正采集动作
func (n *Collector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	ch <- prometheus.MustNewConstMetric(n.gaugeDesc, prometheus.GaugeValue, float64(runtime.NumGoroutine()), "172.28.154.89")
}

func init() {
	prometheus.MustRegister(NewCollector())
}

func main() {
	//registry := collector.NewRegistry()
	//registry.MustRegister(metrics) // 注册指标
	http.Handle("/aaaaa", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":10000", nil))
}
