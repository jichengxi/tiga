package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"tiga-exporter/public/types/local"
)

var (
	NodeMetricMap map[string]nodeMetric
)

func init() {
	NodeMetricMap = make(map[string]nodeMetric)
	NodeMetricMap["local_node_cpu_seconds_total"] = nodeMetric{
		name:      "local_node_cpu_seconds_total",
		help:      "cpu总使用时间.",
		valueType: prometheus.CounterValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			cpuTotal := s.CpuProcessStats.CPUTotal
			cpuSecondsTotal := cpuTotal.User + cpuTotal.Nice + cpuTotal.System + cpuTotal.Idle +
				cpuTotal.Iowait + cpuTotal.IRQ + cpuTotal.SoftIRQ + cpuTotal.Guest + cpuTotal.GuestNice
			return metricValues{
				{
					value: cpuSecondsTotal,
				},
			}
		},
	}
	NodeMetricMap["local_node_cpu_per_mode_seconds_total"] = nodeMetric{
		name:      "local_node_cpu_per_mode_seconds_total",
		help:      "cpu每个模式所使用的总时间.",
		valueType: prometheus.CounterValue,
		labels:    []string{"mode"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			cpuTotal := s.CpuProcessStats.CPUTotal
			return metricValues{
				{labels: []string{"user"}, value: cpuTotal.User},
				{labels: []string{"nice"}, value: cpuTotal.Nice},
				{labels: []string{"system"}, value: cpuTotal.System},
				{labels: []string{"idle"}, value: cpuTotal.Idle},
				{labels: []string{"iowait"}, value: cpuTotal.Iowait},
				{labels: []string{"irq"}, value: cpuTotal.IRQ},
				{labels: []string{"softirq"}, value: cpuTotal.SoftIRQ},
				{labels: []string{"steal"}, value: cpuTotal.Steal},
				{labels: []string{"guest"}, value: cpuTotal.Guest},
				{labels: []string{"guestnice"}, value: cpuTotal.GuestNice},
			}
		},
	}
	NodeMetricMap["local_node_per_cpu_seconds_total"] = nodeMetric{
		name:      "local_node_per_cpu_seconds_total",
		help:      "每个CPU总使用时间.",
		valueType: prometheus.CounterValue,
		labels:    []string{"cpu"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			cpu := s.CpuProcessStats.CPU
			for i, v := range cpu {
				var mV metricValue
				mV.labels = []string{strconv.Itoa(i + 1)}
				mV.value = v.User + v.Nice + v.System + v.Idle +
					v.Iowait + v.IRQ + v.SoftIRQ + v.Guest + v.GuestNice
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_per_cpu_mode_seconds_total"] = nodeMetric{
		name:      "local_node_per_cpu_mode_seconds_total",
		help:      "每个CPU每个模式总使用时间.",
		valueType: prometheus.CounterValue,
		labels:    []string{"cpu", "mode"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			cpu := s.CpuProcessStats.CPU
			//mode := []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal","guest", "guestnice"}
			for i, v := range cpu {
				mVs := metricValues{
					{labels: []string{strconv.Itoa(i + 1), "user"}, value: v.User},
					{labels: []string{strconv.Itoa(i + 1), "nice"}, value: v.Nice},
					{labels: []string{strconv.Itoa(i + 1), "system"}, value: v.System},
					{labels: []string{strconv.Itoa(i + 1), "idle"}, value: v.Idle},
					{labels: []string{strconv.Itoa(i + 1), "iowait"}, value: v.Iowait},
					{labels: []string{strconv.Itoa(i + 1), "irq"}, value: v.IRQ},
					{labels: []string{strconv.Itoa(i + 1), "softirq"}, value: v.SoftIRQ},
					{labels: []string{strconv.Itoa(i + 1), "steal"}, value: v.Steal},
					{labels: []string{strconv.Itoa(i + 1), "guest"}, value: v.Guest},
					{labels: []string{strconv.Itoa(i + 1), "guestnice"}, value: v.GuestNice},
				}
				mVS = append(mVS, mVs...)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_cpu_irq_numbers_total"] = nodeMetric{
		name:      "local_node_cpu_irq_numbers_total",
		help:      "CPU IRQ 中断的次数.",
		valueType: prometheus.CounterValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			iRQTotal := s.CpuProcessStats.IRQTotal
			return metricValues{
				{
					value: float64(iRQTotal),
				},
			}
		},
	}
	NodeMetricMap["local_node_cpu_soft_irq_numbers_total"] = nodeMetric{
		name:      "local_node_cpu_soft_irq_numbers_total",
		help:      "CPU IRQ 软中断的次数.",
		valueType: prometheus.CounterValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			softIRQTotal := s.CpuProcessStats.SoftIRQTotal
			return metricValues{
				{
					value: float64(softIRQTotal),
				},
			}
		},
	}
	NodeMetricMap["local_node_cpu_context_switches_total"] = nodeMetric{
		name:      "local_node_cpu_context_switches_total",
		help:      "CPU上下文切换发生的次数.",
		valueType: prometheus.CounterValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(s.CpuProcessStats.ContextSwitches),
				},
			}
		},
	}
	NodeMetricMap["local_node_process_created_total"] = nodeMetric{
		name:      "local_node_process_created_total",
		help:      "进程创建的次数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(s.CpuProcessStats.ProcessCreated),
				},
			}
		},
	}
	NodeMetricMap["local_node_process_running_total"] = nodeMetric{
		name:      "local_node_process_running_total",
		help:      "当前运行的进程数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(s.CpuProcessStats.ProcessesRunning),
				},
			}
		},
	}
	NodeMetricMap["local_node_process_blocked_total"] = nodeMetric{
		name:      "local_node_process_blocked_total",
		help:      "当前阻塞的进程数(等待 IO).",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(s.CpuProcessStats.ProcessesBlocked),
				},
			}
		},
	}
	NodeMetricMap["local_node_memory_total"] = nodeMetric{
		name:      "local_node_memory_total",
		help:      "总内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.MemTotal),
				},
			}
		},
	}
	NodeMetricMap["local_node_memory_free"] = nodeMetric{
		name:      "local_node_memory_free",
		help:      "空闲内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.MemFree),
				},
			}
		},
	}
	NodeMetricMap["local_node_memory_available"] = nodeMetric{
		name:      "local_node_memory_available",
		help:      "可用内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.MemAvailable),
				},
			}
		},
	}
	NodeMetricMap["local_node_memory_buffer"] = nodeMetric{
		name:      "local_node_memory_buffer",
		help:      "buffer内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.Buffers),
				},
			}
		},
	}
	NodeMetricMap["local_node_memory_cache"] = nodeMetric{
		name:      "local_node_memory_cache",
		help:      "cache内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.Cached),
				},
			}
		},
	}
	NodeMetricMap["local_node_swap_memory_total"] = nodeMetric{
		name:      "local_node_swap_memory_total",
		help:      "swap分区总内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.SwapTotal),
				},
			}
		},
	}
	NodeMetricMap["local_node_swap_memory_free"] = nodeMetric{
		name:      "local_node_swap_memory_free",
		help:      "swap分区空闲内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.SwapFree),
				},
			}
		},
	}
	NodeMetricMap["local_node_swap_memory_cache"] = nodeMetric{
		name:      "local_node_swap_memory_cache",
		help:      "swap分区cache内存.",
		valueType: prometheus.GaugeValue,
		labels:    []string{},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{
					value: float64(*s.MemoryStats.SwapCached),
				},
			}
		},
	}
	NodeMetricMap["local_node_load_average"] = nodeMetric{
		name:      "local_node_load_average",
		help:      "node平均负载.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"load"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			return metricValues{
				{labels: []string{"load1"}, value: s.LoadStats.Load1},
				{labels: []string{"load5"}, value: s.LoadStats.Load5},
				{labels: []string{"load15"}, value: s.LoadStats.Load15},
			}
		},
	}
	NodeMetricMap["local_node_memory_buddy"] = nodeMetric{
		name:      "local_node_memory_buddy",
		help:      "node buddy 物理内存debug信息.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"node", "zone", "page_block_number"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.BuddyStats {
				for a, b := range i.Sizes {
					x := 2 << a
					mV := metricValue{labels: []string{i.Node, i.Zone, strconv.Itoa(x)}, value: b}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_transmit_bytes_total"] = nodeMetric{
		name:      "local_node_network_transmit_bytes_total",
		help:      "出向传输字节数",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.TxBytes)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_transmit_packets_total"] = nodeMetric{
		name:      "local_node_network_transmit_packets_total",
		help:      "出向传输数据包的累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.TxPackets)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_transmit_errors_total"] = nodeMetric{
		name:      "local_node_network_transmit_errors_total",
		help:      "出向传输时错误数据包累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.TxErrors)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_transmit_dropped_total"] = nodeMetric{
		name:      "local_node_network_transmit_dropped_total",
		help:      "出向传输丢弃数据包的累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.TxDropped)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_receive_bytes_total"] = nodeMetric{
		name:      "local_node_network_receive_bytes_total",
		help:      "入向传输字节数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.RxBytes)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_receive_packets_total"] = nodeMetric{
		name:      "local_node_network_receive_packets_total",
		help:      "入向传输数据包的累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.RxBytes)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_receive_errors_total"] = nodeMetric{
		name:      "local_node_network_receive_errors_total",
		help:      "入向传输时错误数据包累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.RxErrors)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_receive_dropped_total"] = nodeMetric{
		name:      "local_node_network_receive_dropped_total",
		help:      "入向传输丢弃数据包的累积计数.",
		valueType: prometheus.CounterValue,
		labels:    []string{"interface"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, v := range s.NetDevStats {
				mV := metricValue{labels: []string{v.Name}, value: float64(v.RxDropped)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_tcp_connect_status"] = nodeMetric{
		name:      "local_node_network_tcp_connect_status",
		help:      "TCP套接字连接状态.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"status"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			tcpStatus := map[string]float64{
				"TCP_ESTABLISHED": 0, // 1
				"TCP_SYN_SENT":    0, // 2
				"TCP_SYN_RECV":    0, // 3
				"TCP_FIN_WAIT1":   0, // 4
				"TCP_FIN_WAIT2":   0, // 5
				"TCP_TIME_WAIT":   0, // 6
				"TCP_CLOSE":       0, // 7
				"TCP_CLOSE_WAIT":  0, // 8
				"TCP_LAST_ACL":    0, // 9
				"TCP_LISTEN":      0, // 10
				"TCP_CLOSING":     0, // 11
			}
			for _, i := range s.TcpStats {
				switch i.St {
				case uint64(1):
					tcpStatus["TCP_ESTABLISHED"] += 1
				case uint64(2):
					tcpStatus["TCP_SYN_SENT"] += 1
				case uint64(3):
					tcpStatus["TCP_SYN_RECV"] += 1
				case uint64(4):
					tcpStatus["TCP_FIN_WAIT1"] += 1
				case uint64(5):
					tcpStatus["TCP_FIN_WAIT2"] += 1
				case uint64(6):
					tcpStatus["TCP_TIME_WAIT"] += 1
				case uint64(7):
					tcpStatus["TCP_CLOSE"] += 1
				case uint64(8):
					tcpStatus["TCP_CLOSE_WAIT"] += 1
				case uint64(9):
					tcpStatus["TCP_LAST_ACL"] += 1
				case uint64(10):
					tcpStatus["TCP_LISTEN"] += 1
				case uint64(11):
					tcpStatus["TCP_CLOSING"] += 1
				}
			}
			for k, v := range tcpStatus {
				mV := metricValue{labels: []string{k}, value: v}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_network_udp_connect_status"] = nodeMetric{
		name:      "local_node_network_udp_connect_status",
		help:      "UDP套接字连接状态.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"status"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			udpStatus := map[string]float64{
				"UDP_ESTABLISHED": 0, // 1
				"UDP_SYN_SENT":    0, // 2
				"UDP_SYN_RECV":    0, // 3
				"UDP_FIN_WAIT1":   0, // 4
				"UDP_FIN_WAIT2":   0, // 5
				"UDP_TIME_WAIT":   0, // 6
				"UDP_CLOSE":       0, // 7
				"UDP_CLOSE_WAIT":  0, // 8
				"UDP_LAST_ACL":    0, // 9
				"UDP_LISTEN":      0, // 10
				"UDP_CLOSING":     0, // 11
			}
			for _, i := range s.UdpStats {
				switch i.St {
				case uint64(1):
					udpStatus["UDP_ESTABLISHED"] += 1
				case uint64(2):
					udpStatus["UDP_SYN_SENT"] += 1
				case uint64(3):
					udpStatus["UDP_SYN_RECV"] += 1
				case uint64(4):
					udpStatus["UDP_FIN_WAIT1"] += 1
				case uint64(5):
					udpStatus["UDP_FIN_WAIT2"] += 1
				case uint64(6):
					udpStatus["UDP_TIME_WAIT"] += 1
				case uint64(7):
					udpStatus["UDP_CLOSE"] += 1
				case uint64(8):
					udpStatus["UDP_CLOSE_WAIT"] += 1
				case uint64(9):
					udpStatus["UDP_LAST_ACL"] += 1
				case uint64(10):
					udpStatus["UDP_LISTEN"] += 1
				case uint64(11):
					udpStatus["UDP_CLOSING"] += 1
				}
			}
			for k, v := range udpStatus {
				mV := metricValue{labels: []string{k}, value: v}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_reads"] = nodeMetric{
		name:      "local_node_disk_io_reads",
		help:      "成功完成的读取次数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.ReadIOs)}
				mVS = append(mVS, mV)
			}

			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_read_ticks"] = nodeMetric{
		name:      "local_node_disk_io_read_ticks",
		help:      "所有读取花费的总毫秒数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.ReadTicks)}
				mVS = append(mVS, mV)
			}

			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_writes"] = nodeMetric{
		name:      "local_node_disk_io_writes",
		help:      "成功完成的写入次数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.WriteIOs)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_write_ticks"] = nodeMetric{
		name:      "local_node_disk_io_write_ticks",
		help:      "所有写入花费的总毫秒数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.WriteTicks)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_progress"] = nodeMetric{
		name:      "local_node_disk_io_progress",
		help:      "当前正在进行的 I/O 数量.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.IOsInProgress)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_io_ticks"] = nodeMetric{
		name:      "local_node_disk_io_ticks",
		help:      "花费在 I/O 上的毫秒数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"device"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskStats {
				mV := metricValue{labels: []string{i.Info.DeviceName}, value: float64(i.IOsTotalTicks)}
				mVS = append(mVS, mV)
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_total"] = nodeMetric{
		name:      "local_node_disk_total",
		help:      "文件系统磁盘总量.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.AllSize),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_free_total"] = nodeMetric{
		name:      "local_node_disk_free_total",
		help:      "文件系统空闲磁盘总量.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.FreeSize),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_usage_total"] = nodeMetric{
		name:      "local_node_disk_usage_total",
		help:      "文件系统已使用磁盘总量.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.UsedSize),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_inode_total"] = nodeMetric{
		name:      "local_node_disk_inode_total",
		help:      "文件系统inode总数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.AllInodes),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_free_inode_total"] = nodeMetric{
		name:      "local_node_disk_free_inode_total",
		help:      "文件系统空闲inode总数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.FreeInodes),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
	NodeMetricMap["local_node_disk_usage_inode_total"] = nodeMetric{
		name:      "local_node_disk_usage_inode_total",
		help:      "文件系统已使用inode总数.",
		valueType: prometheus.GaugeValue,
		labels:    []string{"pod_name", "container_name", "namespace", "image", "device", "mount"},
		getValues: func(s *local.NodeCollecterInfo) metricValues {
			var mVS metricValues
			for _, i := range s.DiskUsageStats {
				for _, disk := range i.DiskUsageList {
					mV := metricValue{
						labels: []string{i.PodName, i.ContainerName, i.NameSpace, i.Image, disk.Device, disk.Mounting},
						value:  float64(disk.DiskUsage.UsedInodes),
					}
					mVS = append(mVS, mV)
				}
			}
			return mVS
		},
	}
}
