package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"tiga-exporter/collect"
	"tiga-exporter/public/types/local"
	localUtils "tiga-exporter/tools/local"
	"tiga-exporter/tools/local/docker"
)

// 数据采集函数
func containerProvider() []*local.ContainerInfo {
	var ContainerInfoList []*local.ContainerInfo
	//dCli := docker.DClient{}
	//dCli.NewClient()
	docker.DCli.GetContainerInfoList(&ContainerInfoList)
	localUtils.GetHostLsblk(&local.LsblkInfos)
	collect.Collecter("/", &ContainerInfoList)
	return ContainerInfoList
}

type containerMetric struct {
	name        string
	help        string
	valueType   prometheus.ValueType
	extraLabels []string
	getValues   func(s local.ContainerCollecterInfo) metricValues
}

func (cm *containerMetric) desc(baseLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(cm.name, cm.help, append(baseLabels, cm.extraLabels...), nil)
}

// ContainerCollector implements collector.Collector.
type ContainerCollector struct {
	errors           prometheus.Gauge
	containerMetrics []containerMetric
	mutex            sync.Mutex
}

func NewContainerCollector(metrics *[]string) *ContainerCollector {
	var containerMetrics []containerMetric

	for _, i := range *metrics {
		if v, ok := ContainerMetricMap[i]; ok {
			containerMetrics = append(containerMetrics, v)
		}
	}

	return &ContainerCollector{
		errors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "local_container",
			Name:      "scrape_error",
			Help:      "1 if there was an error while getting containers metrics, 0 otherwise",
		}),
		containerMetrics: containerMetrics,
	}
}

// Describe 实现采集器Describe接口
func (n *ContainerCollector) Describe(ch chan<- *prometheus.Desc) {
	n.errors.Describe(ch)
	for _, cm := range n.containerMetrics {
		ch <- cm.desc([]string{})
	}
}

// Collect 实现采集器Collect接口,真正采集动作
func (n *ContainerCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	n.collectContainersInfo(ch)
	n.mutex.Unlock()
}

func (n *ContainerCollector) collectContainersInfo(ch chan<- prometheus.Metric) {
	containers := containerProvider()
	for _, container := range containers {
		baseLabels := []string{"kubernetes_docker_type",
			"kubernetes_pod_namespace",
			"kubernetes_pod_name",
			"kubernetes_container_name",
			"kubernetes_pod_uid",
			"kubernetes_sandbox_id",
			"kubernetes_docker_uid",
			"container_label_pod_ip",
			"name",
			"image",
			"host_name",
		}

		baseLabelsValues := []string{container.KubernetesDockerType,
			container.KubernetesPodNamespace,
			container.KubernetesPodName,
			container.KubernetesContainerName,
			container.KubernetesPodUid,
			container.KubernetesSandboxId,
			container.ContainerId,
			container.ContainerLabelPodIp,
			container.Name,
			container.Image,
			container.HostName,
		}
		for _, cm := range n.containerMetrics {
			desc := cm.desc(baseLabels)
			for _, metricV := range cm.getValues(container.ContainerCollecterInfo) {
				ch <- prometheus.MustNewConstMetric(desc, cm.valueType, metricV.value, append(baseLabelsValues, metricV.labels...)...)
			}
		}
	}

}
