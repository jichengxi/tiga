package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"tiga-exporter/public/types/remote"
)

type metricValue struct {
	value  float64
	labels []string
}

type metricValues []metricValue

type containerMetric struct {
	name      string
	help      string
	valueType prometheus.ValueType
	getValues func(s []remote.MetricValue) metricValues
}

func (cm *containerMetric) desc(baseLabels []string) *prometheus.Desc {
	if baseLabels == nil {
		return prometheus.NewDesc(cm.name, cm.help, nil, nil)
	}
	return prometheus.NewDesc(cm.name, cm.help, append(baseLabels), nil)
}

// RemoteCollector implements prometheus.Collector.
type RemoteCollector struct {
	errors           prometheus.Gauge
	containerMetrics []containerMetric
	mutex            sync.Mutex
}

func NewRemoteCollector(rM *map[string]*remote.MetricTitle) *RemoteCollector {
	var cMs []containerMetric
	for _, v := range *rM {
		cMs = append(cMs, containerMetric{
			name:      v.Name,
			help:      v.Help,
			valueType: prometheus.GaugeValue,
			getValues: func(s []remote.MetricValue) metricValues {
				mVs := metricValues{}
				for _, i := range s {
					mVs = append(mVs, metricValue{labels: i.LabelsValue, value: i.Value})
				}
				return mVs
			},
		})
	}
	return &RemoteCollector{
		errors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "remote_container",
			Name:      "scrape_error",
			Help:      "1 if there was an error while getting containers metrics, 0 otherwise",
		}),
		containerMetrics: cMs,
	}
}

func (n *RemoteCollector) Describe(ch chan<- *prometheus.Desc) {
	n.errors.Describe(ch)
	for _, cm := range n.containerMetrics {
		ch <- cm.desc([]string{})
	}
}

func (n *RemoteCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	n.collectContainersInfo(ch)
	n.mutex.Unlock()
}

func (n *RemoteCollector) collectContainersInfo(ch chan<- prometheus.Metric) {
	metricData := remote.MetricData{}
	infoByLabelsMap := metricData.InfoByLabelsMap
	infoNotLabelsMap := metricData.InfoNotLabelsMap
	for _, cm := range n.containerMetrics {
		if v, ok := infoByLabelsMap[cm.name]; ok {
			desc := cm.desc(v.Labels)
			metricValueList := v.MetricValue
			for _, metricV := range cm.getValues(metricValueList) {
				ch <- prometheus.MustNewConstMetric(desc, cm.valueType, metricV.value, metricV.labels...)
			}
		}
		if v, ok := infoNotLabelsMap[cm.name]; ok {
			desc := cm.desc(nil)
			metricV := v.Value
			ch <- prometheus.MustNewConstMetric(desc, cm.valueType, metricV)
		}
	}

}
