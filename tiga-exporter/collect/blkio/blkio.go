package blkio

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"tiga-exporter/public/types/local"
	utils "tiga-exporter/tools/local"
)

func GetBlkioData(rootfs string, container local.ContainerInfo, stats *cgroups.Stats) {
	blkioGroup := fs.BlkioGroup{}
	blkioMountDir, err := utils.GetThisCGroupDir("blkio", container.Pid)
	if err != nil {
		panic(err)
	}
	// blkioMountDir:  /kubepods.slice/kubepods-burstable.slice/kubepods-burstable-podd1c54b44_724f_4f4e_af25_cd1e3c15eddc.slice/docker-bd2fe3511cf6212da745edb16913bad049a5fd0b0484071b7e1f39937f38a947.scop
	//fmt.Println("blkio: ", blkioMountDir)
	// rootMount:  /sys/fs/cgroup/blkio
	rootMount, _ := cgroups.FindCgroupMountpoint(rootfs, "blkio")
	//fmt.Println("mount: ", rootMount)
	blkioPath := rootMount + blkioMountDir
	err = blkioGroup.GetStats(blkioPath, stats)
	if err != nil {
		_ = fmt.Errorf(err.Error() + "\n")
	}
	//fmt.Println("-----------")
	//fmt.Println(stats.BlkioStats)
	//fmt.Println("-----------")
}
