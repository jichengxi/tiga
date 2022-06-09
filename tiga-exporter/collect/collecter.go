package collect

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"strconv"
	"tiga-exporter/collect/blkio"
	"tiga-exporter/collect/cpu"
	"tiga-exporter/collect/memory"
	"tiga-exporter/collect/network"
	"tiga-exporter/public/types/local"
)

func Collecter(rootfs string, containerInfo *[]*local.ContainerInfo) {
	for _, i := range *containerInfo {
		stats := local.ContainerCollecterInfo{}
		stats.MemoryStats.Stats = make(map[string]uint64)
		stats.HugetlbStats = make(map[string]cgroups.HugetlbStats)
		//fmt.Println("pid: ", i.Pid, i.ContainerId)
		memory.GetMemoryData(rootfs, *i, &stats.Stats)

		cpu.GetCpuGroup(rootfs, *i, &stats.Stats)

		cpu.GetCpuacctGroup(rootfs, *i, &stats.Stats)

		blkio.GetBlkioData(rootfs, *i, &stats.Stats)
		//fmt.Println(stats.BlkioStats.IoServiceBytesRecursive)
		//fmt.Println(i)
		//fmt.Println("aaaa: ", stats.BlkioStats.IoServiceBytesRecursive)
		diskBytes := make(map[string]local.DiskByte)
		if len(stats.BlkioStats.IoServiceBytesRecursive) != 0 {
			for _, diskValue := range stats.BlkioStats.IoServiceBytesRecursive {
				majMin := strconv.FormatUint(diskValue.Major, 10) + ":" + strconv.FormatUint(diskValue.Minor, 10)
				//fmt.Println("+++++++", majMin)
				lsblkInfo := local.LsblkInfos[majMin]
				var labelName string
				// mountDestination 取目的地挂载目录
				//var mountDestination string
				//mountDestination = i.ContainerMounts[strings.TrimPrefix(lsblkInfo.MountPoint, "/rootfs")]
				//if len(mountDestination) != 0 {
				//    labelName = mountDestination
				//} else {
				//    labelName = lsblkInfo.Name
				//}
				labelName = lsblkInfo.Name
				switch diskValue.Op {
				case "Read":
					x := diskBytes[labelName]
					x.ReadBytes = diskValue.Value
					diskBytes[labelName] = x
				case "Write":
					x := diskBytes[labelName]
					x.WriteBytes = diskValue.Value
					diskBytes[labelName] = x
				case "Sync":
					x := diskBytes[labelName]
					x.SyncBytes = diskValue.Value
					diskBytes[labelName] = x
				case "Async":
					x := diskBytes[labelName]
					x.AsyncBytes = diskValue.Value
					diskBytes[labelName] = x
				case "Total":
					x := diskBytes[labelName]
					x.TotalBytes = diskValue.Value
					diskBytes[labelName] = x
				}
			}
		}
		//fmt.Println(diskBytes)
		//jsons, _ := json.Marshal(stats.BlkioStats)
		//fmt.Println(string(jsons))

		netInfo := network.GetNetworkData(i.Pid)

		//fmt.Println("stats: ", stats.MemoryStats.Stats)
		stats.Disks = diskBytes
		stats.Networks = netInfo
		i.ContainerCollecterInfo = stats
	}
}
