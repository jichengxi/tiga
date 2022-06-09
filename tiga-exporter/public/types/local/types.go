package local

import (
	dclient "github.com/fsouza/go-dockerclient"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
)

var (
	LsblkInfos map[string]LsblkInfo
)

type LsblkInfo struct {
	Name       string
	MountPoint string
}

type DiskByte struct {
	ReadBytes  uint64
	WriteBytes uint64
	SyncBytes  uint64
	AsyncBytes uint64
	TotalBytes uint64
}

type DiskUsage struct {
	AllSize    uint64
	FreeSize   uint64
	UsedSize   uint64
	AllInodes  uint64
	FreeInodes uint64
	UsedInodes uint64
}

type ContainerCollecterInfo struct {
	cgroups.Stats
	Networks map[string]dclient.NetworkStats
	Disks    map[string]DiskByte
}

type ContainerInfo struct {
	ContainerId             string            `json:"container_id" yaml:"container_id"`
	Pid                     int               `json:"pid" yaml:"pid"`
	Name                    string            `json:"name" yaml:"name"`
	HostName                string            `json:"host_name" yaml:"host_name"`
	Image                   string            `json:"image" yaml:"image"`
	KubernetesDockerType    string            `json:"kubernetes_docker_type" yaml:"kubernetes_docker_type"`
	KubernetesPodNamespace  string            `json:"kubernetes_pod_namespace" yaml:"kubernetes_pod_namespace"`
	KubernetesPodName       string            `json:"kubernetes_pod_name" yaml:"kubernetes_pod_name"`
	KubernetesContainerName string            `json:"kubernetes_container_name" yaml:"kubernetes_container_name"`
	KubernetesPodUid        string            `json:"kubernetes_pod_uid" yaml:"kubernetes_pod_uid"`
	KubernetesSandboxId     string            `json:"kubernetes_sandbox_id" yaml:"kubernetes_sandbox_id"`
	ContainerLabelPodIp     string            `json:"container_label_pod_ip" yaml:"container_label_pod_ip"`
	ContainerMounts         map[string]string `json:"container_mounts" yaml:"container_mounts"`
	MergedDir               string            `json:"merged_dir" yaml:"merged_dir"`
	ContainerCollecterInfo  ContainerCollecterInfo
}

type DevUsage struct {
	Device    string
	Mounting  string
	DiskUsage DiskUsage
}

type DiskUsageStat struct {
	PodName       string
	ContainerName string
	NameSpace     string
	Image         string
	DiskUsageList []DevUsage
}

type NodeCollecterInfo struct {
	CpuProcessStats  procfs.Stat                `json:"cpu_process_stats"`
	MemoryStats      procfs.Meminfo             `json:"memory_stats"`
	NetDevStats      procfs.NetDev              `json:"net_dev_stats"`
	TcpStats         procfs.NetTCP              `json:"tcp_stats"`
	UdpStats         procfs.NetUDP              `json:"udp_stats"`
	IpvsStats        procfs.IPVSStats           `json:"ipvs_stats"`
	IpvsBackendStats []procfs.IPVSBackendStatus `json:"ipvs_backend_stats"`
	LoadStats        *procfs.LoadAvg            `json:"load_stats"`
	BuddyStats       []procfs.BuddyInfo         `json:"buddy_stats"`
	DiskStats        []blockdevice.Diskstats    `json:"disk_stats"`
	DiskUsageStats   []DiskUsageStat
}

func init() {
	LsblkInfos = make(map[string]LsblkInfo)
}
