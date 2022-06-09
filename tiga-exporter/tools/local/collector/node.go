package collector

import (
	"fmt"
	gofstab "github.com/deniswernert/go-fstab"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"strings"
	"sync"
	"tiga-exporter/public/types/local"
	utils "tiga-exporter/tools/local"
	"tiga-exporter/tools/local/docker"
)

func diskUsageCollecter() []local.DiskUsageStat {
	var diskInfos []local.DiskUsageStat
	loopMap := make(map[string]string)
	mounts, _ := gofstab.ParseProc()
	for _, val := range mounts {
		if val.Spec == "sysfs" || val.Spec == "proc" || val.Spec == "devtmpfs" || val.Spec == "securityfs" ||
			val.Spec == "devpts" || val.Spec == "cgroup" || val.Spec == "configfs" || val.Spec == "systemd-1" ||
			val.Spec == "debugfs" || val.Spec == "mqueue" || val.Spec == "hugetlbfs" || val.Spec == "binfmt_misc" ||
			val.Spec == "fusectl" || val.Spec == "/sys/fs/bpf" || val.File == "/dev/shm" || val.File == "/run" ||
			val.File == "/sys/fs/pstore" || val.File == "/sys/fs/cgroup" ||
			strings.Contains(val.File, "/run/user/") {
			continue
		}
		fmt.Println(val)
		if strings.Contains(val.Spec, "/dev/loop") || val.Spec == "overlay" {
			loopMap[val.File] = val.Spec
			continue
		}
		diskInfo := utils.DiskUsage(val.File)
		devUsage := local.DevUsage{Device: val.Spec, Mounting: val.File, DiskUsage: diskInfo}
		a := local.DiskUsageStat{DiskUsageList: []local.DevUsage{devUsage}}
		diskInfos = append(diskInfos, a)
	}
	fmt.Println("loopMap: ", loopMap)
	var ContainerInfoList []*local.ContainerInfo
	docker.DCli.GetContainerInfoList(&ContainerInfoList)
	for _, i := range ContainerInfoList {
		var diskUsageList []local.DevUsage
		for k, v := range i.ContainerMounts {
			if v == "/data" {
				diskInfo := utils.DiskUsage(k)
				devUsage := local.DevUsage{Device: loopMap[k], Mounting: k, DiskUsage: diskInfo}
				diskUsageList = append(diskUsageList, devUsage)
			}
		}
		diskInfo := utils.DiskUsage(i.MergedDir)
		devUsage := local.DevUsage{Device: loopMap[i.MergedDir], Mounting: i.MergedDir, DiskUsage: diskInfo}
		diskUsageList = append(diskUsageList, devUsage)

		a := local.DiskUsageStat{
			PodName:       i.KubernetesPodName,
			ContainerName: i.KubernetesContainerName,
			NameSpace:     i.KubernetesPodNamespace,
			Image:         i.Image,
			DiskUsageList: diskUsageList,
		}
		diskInfos = append(diskInfos, a)
	}
	return diskInfos
}

// 数据采集函数
func nodeProvider() local.NodeCollecterInfo {
	var collecterInfo local.NodeCollecterInfo
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		panic(err)
	}
	//cpuInfo, _ := fs.CPUInfo()
	cpuProcessStats, _ := fs.Stat()
	collecterInfo.CpuProcessStats = cpuProcessStats
	//fmt.Println("cpuProcessStats: ", cpuProcessStats)
	memoryStats, _ := fs.Meminfo()
	collecterInfo.MemoryStats = memoryStats
	//fmt.Println("memoryStats: ", memoryStats)
	// /proc/net/dev
	netDevStats, _ := fs.NetDev()
	collecterInfo.NetDevStats = netDevStats
	//fmt.Println("netDevStats: ", netDevStats)
	// /proc/net/tcp
	tcpStats, _ := fs.NetTCP()
	collecterInfo.TcpStats = tcpStats
	//fmt.Println("tcpStats: ", tcpStats)
	//for _, i := range tcpStats {
	//	fmt.Println("tcpStats: ", *i)
	//}
	// /proc/net/udp
	udpStats, _ := fs.NetUDP()
	collecterInfo.UdpStats = udpStats
	//fmt.Println("udpStats: ", udpStats)
	// /proc/net/ip_vs_stats
	ipvsStats, _ := fs.IPVSStats()
	collecterInfo.IpvsStats = ipvsStats
	//fmt.Println("ipvsStats: ", ipvsStats)
	ipvsBackendStats, _ := fs.IPVSBackendStatus()
	collecterInfo.IpvsBackendStats = ipvsBackendStats
	//fmt.Println("ipvsBackendStats: ", ipvsBackendStats)
	// /proc/loadavg
	loadStats, _ := fs.LoadAvg()
	collecterInfo.LoadStats = loadStats
	//fmt.Println("loadStats: ", loadStats)
	// /proc/buddyinfo
	buddyStats, _ := fs.BuddyInfo()
	collecterInfo.BuddyStats = buddyStats
	//fmt.Println("buddyStats: ", buddyStats)

	blockDeviceFs, err := blockdevice.NewFS("/proc", "/sys")
	if err != nil {
		panic(err)
	}
	diskStats, _ := blockDeviceFs.ProcDiskstats()
	//devices, _ := blockDeviceFs.SysBlockDevices()
	//fmt.Println("devices: ", devices)

	collecterInfo.DiskStats = diskStats
	//fmt.Println("diskStats: ", diskStats)

	diskUsageList := diskUsageCollecter()
	collecterInfo.DiskUsageStats = diskUsageList
	//fmt.Println("diskUsageList: ", diskUsageList)

	return collecterInfo
}

type nodeMetric struct {
	name      string
	help      string
	labels    []string
	valueType prometheus.ValueType
	getValues func(s *local.NodeCollecterInfo) metricValues
}

func (cm *nodeMetric) desc(baseLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(cm.name, cm.help, baseLabels, nil)
}

type NodeCollector struct {
	errors      prometheus.Gauge
	nodeMetrics []nodeMetric
	mutex       sync.Mutex
}

func NewNodeCollector(metrics *[]string) *NodeCollector {
	var nodeMetrics []nodeMetric

	for _, i := range *metrics {
		if v, ok := NodeMetricMap[i]; ok {
			nodeMetrics = append(nodeMetrics, v)
		}
	}

	return &NodeCollector{
		errors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "local_node",
			Name:      "scrape_error",
			Help:      "1 if there was an error while getting node metrics, 0 otherwise",
		}),
		nodeMetrics: nodeMetrics,
	}
}

// Describe 实现采集器Describe接口
func (n *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	n.errors.Describe(ch)
	for _, cm := range n.nodeMetrics {
		ch <- cm.desc([]string{})
	}
}

// Collect 实现采集器Collect接口,真正采集动作
func (n *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	n.collectNodeInfo(ch)
	n.mutex.Unlock()
}

func (n *NodeCollector) collectNodeInfo(ch chan<- prometheus.Metric) {
	nodeData := nodeProvider()
	for _, nm := range n.nodeMetrics {
		if len(nm.labels) == 0 {
			desc := nm.desc(nil)
			for _, metricV := range nm.getValues(&nodeData) {
				ch <- prometheus.MustNewConstMetric(desc, nm.valueType, metricV.value)
			}
		} else {
			desc := nm.desc(nm.labels)
			for _, metricV := range nm.getValues(&nodeData) {
				ch <- prometheus.MustNewConstMetric(desc, nm.valueType, metricV.value, metricV.labels...)
			}
		}
	}
}
