package memory

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"tiga-exporter/public/types/local"
	utils "tiga-exporter/tools/local"
)

func GetMemoryData(rootfs string, container local.ContainerInfo, stats *cgroups.Stats) {
	// map[active_anon:36864 active_file:0 cache:0 hierarchical_memory_limit:8201248768 hierarchical_memsw_limit:9223372036854771712 inactive_anon:0 inactive_file:0 mapped_file:0 pgfault:2329 pgmajfault:0 pgpgin:1036 pgpgout:1027 rss:36864 rss_huge:0 swap:0 total_active_anon:36864 total_active_file:0 total_cache:0 total_inactive_anon:0 total_inactive_file:0 total_mapped_file:0 total_pgfault:0 total_pgmajfault:0 total_pgpgin:0 total_pgpgout:0 total_rss:36864 total_rss_huge:0 total_swap:0 total_unevictable:0 unevictable:0]
	memoryGroup := fs.MemoryGroup{}
	memoryMountDir, err := utils.GetThisCGroupDir("memory", container.Pid)
	if err != nil {
		panic(err)
	}
	//fmt.Println("memory: ", memoryMountDir)
	rootMount, _ := cgroups.FindCgroupMountpoint(rootfs, "memory")
	//fmt.Println("mount: ", rootMount)
	memoryPath := rootMount + memoryMountDir

	err = memoryGroup.GetStats(memoryPath, stats)
	if err != nil {
		_ = fmt.Errorf(err.Error() + "\n")
	}
}
