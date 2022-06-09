package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"tiga-exporter/public"
	"tiga-exporter/public/types/remote"
	"tiga-exporter/tools"
	localUtils "tiga-exporter/tools/local"
	localCollector "tiga-exporter/tools/local/collector"
	remoteUtils "tiga-exporter/tools/remote"
	remoteCollector "tiga-exporter/tools/remote/collector"
	"time"
)

var (
	TConfig tools.TigaConfig // 配置存储容器

	tempRemoteNames []string // 用来保存上一次监控项的名称列表

	currentLocalContainerMetrics []string // 当前容器部分监控项存储容器
	currentLocalNodeMetrics      []string // 当前节点部分监控项存储容器

	localContainerMetricAll []string // 代码中所有的容器监控项存储容器
	localNodeMetricAll      []string // 代码中所有的节点监控项存储容器

	registerRemoteCollectors    *remoteCollector.RemoteCollector
	registerLocalCollectors     *localCollector.ContainerCollector
	registerLocalNodeCollectors *localCollector.NodeCollector
	enableProfiling             bool

	collectorUrl string
	port         string
)

func init() {
	flag.StringVar(&collectorUrl, "collector", "", "collector的地址")
	flag.StringVar(&port, "port", "10001", "exporter的端口")
	localUtils.RootDir()
	flag.Parse()
	// 所有的容器监控项
	for k := range localCollector.ContainerMetricMap {
		localContainerMetricAll = append(localContainerMetricAll, k)
	}
	// 所有的宿主机监控项
	for k := range localCollector.NodeMetricMap {
		localNodeMetricAll = append(localNodeMetricAll, k)
	}
	enableProfiling = true
}

func main() {
	// 获取配置，写到全局配置TConfig里
	go func(t *tools.TigaConfig, url *string) {
		for {
			*t = tools.RequestTigaConfig(*url)
			time.Sleep(time.Second * 30)
		}
	}(&TConfig, &collectorUrl)

	//解析远程配置，拿到数据
	go func(tr *tools.Remote, mRMT *map[string]*remote.MetricTitle, rMD *remote.MetricData) {
		for {
			if tr.Enable {
				remoteUtils.ParseRemoteConfig(tr, mRMT, rMD)
			}
			time.Sleep(time.Second * 5)
		}
	}(&TConfig.Remote, &public.MetricTitleMap, &public.MetricData)

	//拿到数据的情况下，进行注册，采集展示(数据结果已经放到对应的struct中)
	go func(mtMap *map[string]*remote.MetricTitle, rP *remoteCollector.RemoteCollector, tRN *[]string) {
		for {
			if len(*mtMap) != 0 {
				break
			}
			time.Sleep(time.Second)
		}
		for {
			var metricNames []string // 存储当前监控项名称列表
			for mt := range *mtMap {
				metricNames = append(metricNames, mt)
			}
			if len(*tRN) == 0 {
				// 第一次进来,上一次监控项列表中没有数据
				rP = remoteCollector.NewRemoteCollector(mtMap)
				prometheus.MustRegister(rP)
				*tRN = metricNames
			} else {
				// 比对当前监控项和上一次的监控项是否有差别, 这一步可以检测配置的监控项是否有变更
				res := tools.CompareSlice(*tRN, metricNames)
				if !res {
					prometheus.Unregister(rP)
					rP = remoteCollector.NewRemoteCollector(mtMap)
					prometheus.MustRegister(rP)
					*tRN = metricNames
				}
			}
			// TODO: 目前是每10秒检测一次监控项，并更新注册一次，目标是通过消息通知进行更新注册
			time.Sleep(time.Second * 10)
		}
	}(&public.MetricTitleMap, registerRemoteCollectors, &tempRemoteNames)

	// 采集container数据
	go func(containerConfig *tools.LocalMetricConfig, lPC *localCollector.ContainerCollector, currentMetrics *[]string) {
		time.Sleep(time.Second * 1)
		for {
			if containerConfig.Enable {
				// 判断配置文件中是否有指定监控项，如果没有则监控所有
				tempContainerMetrics := containerConfig.Metrics
				if len(tempContainerMetrics) == 0 {
					tempContainerMetrics = localContainerMetricAll
				}
				// tempLocalContainerMetrics 在第一次启动时，里面没有当前监控项，直接注册并赋值
				if len(*currentMetrics) == 0 {
					lPC = localCollector.NewContainerCollector(&tempContainerMetrics)
					prometheus.MustRegister(lPC)
					*currentMetrics = tempContainerMetrics
				} else {
					// 如果监控项发生变化，监控项发生变化，取消注册后再注册，并把新的值写入当前监控项存储容器
					res := tools.CompareSlice(*currentMetrics, tempContainerMetrics)
					if !res {
						prometheus.Unregister(lPC)
						lPC = localCollector.NewContainerCollector(&tempContainerMetrics)
						prometheus.MustRegister(lPC)
						*currentMetrics = tempContainerMetrics
					}
				}
			}
			time.Sleep(time.Second * 10)
		}
	}(&TConfig.Local.Container, registerLocalCollectors, &currentLocalContainerMetrics)

	// 采集node数据
	go func(nodeConfig *tools.LocalMetricConfig, lPN *localCollector.NodeCollector, currentMetrics *[]string) {
		time.Sleep(time.Second * 1)
		for {
			if nodeConfig.Enable {
				tempLocalMetrics := nodeConfig.Metrics
				if len(tempLocalMetrics) == 0 {
					tempLocalMetrics = localNodeMetricAll
				}

				if len(*currentMetrics) == 0 {
					lPN = localCollector.NewNodeCollector(&tempLocalMetrics)
					prometheus.MustRegister(lPN)
					*currentMetrics = tempLocalMetrics
				} else {
					res := tools.CompareSlice(*currentMetrics, tempLocalMetrics)
					if !res {
						prometheus.Unregister(lPN)
						lPN = localCollector.NewNodeCollector(&tempLocalMetrics)
						prometheus.MustRegister(lPN)
						*currentMetrics = tempLocalMetrics
					}
				}
			}
			time.Sleep(time.Second * 10)
		}
	}(&TConfig.Local.Node, registerLocalNodeCollectors, &currentLocalNodeMetrics)

	mux := http.NewServeMux()

	if enableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}
	mux.Handle("/metrics", promhttp.Handler())
	glog.Fatal(http.ListenAndServe(":"+port, mux))
}
