package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"tiga-exporter/public/types/local"
)

var (
	ContainerMetricMap map[string]containerMetric
)

func init() {
	ContainerMetricMap = make(map[string]containerMetric)
	// cpu
	ContainerMetricMap["local_container_cpu_usage_seconds_total"] = containerMetric{
		name:      "local_container_cpu_usage_seconds_total",
		help:      "Cumulative cpu time consumed in seconds./ 针对每个CPU累计消耗的CPU时间。如果有多个CPU，则总的CPU时间需要把各个CPU耗费的时间相加.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.CpuUsage.TotalUsage)}}
		},
	}
	ContainerMetricMap["local_container_cpu_user_seconds_total"] = containerMetric{
		name:      "local_container_cpu_user_seconds_total",
		help:      "Cumulative user cpu time consumed in seconds./ 累计消耗的用户（user）CPU时间.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.CpuUsage.UsageInUsermode)}}
		},
	}
	ContainerMetricMap["local_container_cpu_system_seconds_total"] = containerMetric{
		name:      "local_container_cpu_system_seconds_total",
		help:      "Cumulative system cpu time consumed in seconds./ 累计消耗的系统（system）CPU时间.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.CpuUsage.UsageInKernelmode)}}
		},
	}
	ContainerMetricMap["local_container_cpu_cfs_throttled_seconds_total"] = containerMetric{
		name:      "local_container_cpu_cfs_throttled_seconds_total",
		help:      "Total time duration the containers has been throttled./ 被限制的总持续时间.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.ThrottlingData.ThrottledTime)}}
		},
	}
	ContainerMetricMap["local_container_cpu_cfs_periods_total"] = containerMetric{
		name:      "local_container_cpu_cfs_periods_total",
		help:      "Number of elapsed enforcement period intervals./ 已经执行的CPU时间周期数.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.ThrottlingData.ThrottledTime)}}
		},
	}
	ContainerMetricMap["local_container_cpu_cfs_throttled_periods_total"] = containerMetric{
		name:      "local_container_cpu_cfs_throttled_periods_total",
		help:      "Number of throttled period intervals./ 被限制或节流的CPU时间周期数.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.CpuStats.ThrottlingData.ThrottledPeriods)}}
		},
	}

	// memory
	ContainerMetricMap["local_container_memory_usage_bytes"] = containerMetric{
		name:      "local_container_memory_usage_bytes",
		help:      "Current memory usage in bytes, including all memory regardless of when it was accessed. / 当前内存使用量.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.MemoryStats.Usage.Usage)}}
		},
	}
	ContainerMetricMap["local_container_memory_max_usage_bytes"] = containerMetric{
		name:      "local_container_memory_max_usage_bytes",
		help:      "Maximum memory usage recorded in bytes. / 最大内存使用量.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.MemoryStats.Usage.MaxUsage)}}
		},
	}
	ContainerMetricMap["local_container_memory_failcnt"] = containerMetric{
		name:      "local_container_memory_failcnt",
		help:      "Number of memory usage hits limits. / 内存使用达到限制的次数",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.MemoryStats.Usage.Failcnt)}}
		},
	}
	ContainerMetricMap["local_container_memory_rss"] = containerMetric{
		name:      "local_container_memory_rss",
		help:      "Size of RSS in bytes. / 内存RSS的大小",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.MemoryStats.Stats["rss"])}}
		},
	}
	ContainerMetricMap["local_container_memory_swap"] = containerMetric{
		name:      "local_container_memory_swap",
		help:      "Container swap usage in bytes. / 虚拟内存使用量.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			return metricValues{{value: float64(s.MemoryStats.SwapUsage.Usage)}}
		},
	}
	// TODO: container_memory_working_set_bytes 计算不准，待考量
	ContainerMetricMap["local_container_memory_working_set_bytes"] = containerMetric{
		name:      "local_container_memory_working_set_bytes",
		help:      "Current working set in bytes. / 当前内存工作集（working set）使用量.",
		valueType: prometheus.GaugeValue,
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			// working = rss + cache
			return metricValues{{value: float64(s.MemoryStats.Stats["rss"]) + float64(s.MemoryStats.Cache)}}
		},
	}

	// network transmit
	ContainerMetricMap["local_container_network_transmit_bytes_total"] = containerMetric{
		name:        "local_container_network_transmit_bytes_total",
		help:        "Cumulative count of bytes transmitted. / 出向流量大小，单位字节.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.TxBytes)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_transmit_packets_total"] = containerMetric{
		name:        "local_container_network_transmit_packets_total",
		help:        "Cumulative count of packets transmitted. / 出向传输数据包的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.TxPackets)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_transmit_errors_total"] = containerMetric{
		name:        "local_container_network_transmit_errors_total",
		help:        "Cumulative count of errors encountered while transmitting. / 出向传输时遇到的错误累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.TxErrors)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_transmit_packets_dropped_total"] = containerMetric{
		name:        "local_container_network_transmit_packets_dropped_total",
		help:        "Cumulative count of packets dropped while transmitting. / 出向传输时丢弃的数据包的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.TxDropped)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}

	// network receive
	ContainerMetricMap["local_container_network_receive_bytes_total"] = containerMetric{
		name:        "local_container_network_receive_bytes_total",
		help:        "Cumulative count of bytes received. / 入向流量大小，单位字节.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.RxBytes)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_receive_packets_total"] = containerMetric{
		name:        "local_container_network_receive_packets_total",
		help:        "Cumulative count of packets received. / 入向传输数据包的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.RxPackets)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_receive_errors_total"] = containerMetric{
		name:        "local_container_network_receive_errors_total",
		help:        "Cumulative count of errors encountered while receiving. / 入向传输时遇到的错误累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.RxErrors)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_network_receive_packets_dropped_total"] = containerMetric{
		name:        "local_container_network_receive_packets_dropped_total",
		help:        "Cumulative count of packets dropped while receiving. / 入向传输时丢弃的数据包的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"interface"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Networks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.RxDropped)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}

	// TODO: 待重写，取值太垃圾
	// fs
	ContainerMetricMap["local_container_fs_reads_total"] = containerMetric{
		name:        "local_container_fs_reads_total",
		help:        "Cumulative count of reads completed. / 已完成读取的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"device"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Disks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.ReadBytes)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
	ContainerMetricMap["local_container_fs_writes_total"] = containerMetric{
		name:        "local_container_fs_writes_total",
		help:        "Cumulative count of writes completed. / 已完成写入的累积计数.",
		valueType:   prometheus.GaugeValue,
		extraLabels: []string{"device"},
		getValues: func(s local.ContainerCollecterInfo) metricValues {
			var mVs metricValues
			for k, v := range s.Disks {
				var mV metricValue
				mV.labels = []string{k}
				mV.value = float64(v.WriteBytes)
				mVs = append(mVs, mV)
			}
			return mVs
		},
	}
}
