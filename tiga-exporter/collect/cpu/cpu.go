package cpu

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"tiga-exporter/public/types/local"
	utils "tiga-exporter/tools/local"
)

func GetCpuGroup(rootfs string, container local.ContainerInfo, stats *cgroups.Stats) {
	cpuGroup := fs.CpuGroup{}
	cpuMountDir, err := utils.GetThisCGroupDir("cpu", container.Pid)
	if err != nil {
		panic(err)
	}
	rootMount, _ := cgroups.FindCgroupMountpoint(rootfs, "cpu")
	//fmt.Println("cpu", cpuMountDir)
	//fmt.Println("mount: ", rootMount)
	cpuPath := rootMount + cpuMountDir
	err = cpuGroup.GetStats(cpuPath, stats)
	if err != nil {
		_ = fmt.Errorf(err.Error() + "\n")
	}
}

func GetCpusetGroup(path string, stats *cgroups.Stats) {
	cpusetGroup := fs.CpusetGroup{}
	err := cpusetGroup.GetStats(path, stats)
	if err != nil {
		_ = fmt.Errorf(err.Error() + "\n")
	}
}

func GetCpuacctGroup(rootfs string, container local.ContainerInfo, stats *cgroups.Stats) {
	cpuacctGroup := fs.CpuacctGroup{}
	cpuacctMountDir, err := utils.GetThisCGroupDir("cpuacct", container.Pid)
	if err != nil {
		panic(err)
	}
	rootMount, _ := cgroups.FindCgroupMountpoint(rootfs, "cpuacct")
	//fmt.Println("cpuacct", cpuacctMountDir)
	//fmt.Println("mount: ", rootMount)
	cpuacctPath := rootMount + cpuacctMountDir
	err = cpuacctGroup.GetStats(cpuacctPath, stats)
	if err != nil {
		_ = fmt.Errorf(err.Error() + "\n")
	}
}
