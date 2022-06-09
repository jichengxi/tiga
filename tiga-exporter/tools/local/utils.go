package utils

import (
	"context"
	cgs "github.com/opencontainers/runc/libcontainer/cgroups"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"tiga-exporter/public"
	"tiga-exporter/public/types/local"
	"time"
)

const cgroupNamePrefix = "name="

func getControllerPath(subsystem string, cgroups map[string]string) (string, error) {

	if p, ok := cgroups[subsystem]; ok {
		return p, nil
	}

	if p, ok := cgroups[cgroupNamePrefix+subsystem]; ok {
		return p, nil
	}

	return "", cgs.NewNotFoundError(subsystem)
}

func GetThisCGroupDir(subsystem string, pid int) (string, error) {
	pidTemp := ""
	if pid == 0 || pid == 1 {
		pidTemp = "self"
	} else {
		pidTemp = strconv.Itoa(pid)
	}
	cgroups, err := cgs.ParseCgroupFile("/proc/" + pidTemp + "/cgroup")

	if err != nil {
		return "", err
	}

	return getControllerPath(subsystem, cgroups)
}

// RootDir 判断系统根目录
func RootDir() {
	_, err := os.Stat("/rootfs")
	if err == nil {
		err = syscall.Chroot("/rootfs")
		if err != nil {
			panic(err)
		}
		public.Rootfs = "/rootfs"
	}
	public.Rootfs = "/"
}

// FileRead 读取文件
func FileRead(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	return ioutil.ReadAll(f)
}

// CmdTime 执行linux命令获取标准输出和标准错误并设置超时时间
func CmdTime(command string) (string, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer c()
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	out, err := cmd.CombinedOutput()
	result := string(out)
	return result, err
}

// GetHostLsblk 获取主机磁盘映射关系
func GetHostLsblk(Lsblkmap *map[string]local.LsblkInfo) {
	lsblkMap := make(map[string]local.LsblkInfo)
	lsblk, err := CmdTime("lsblk -o MAJ:MIN,MOUNTPOINT,NAME -n -P")
	if err != nil {
		panic(err)
	}
	rex := regexp.MustCompile(`MAJ:MIN="(.*)" MOUNTPOINT="(.*)" NAME="(.*)"$`)
	// MAJ:MIN="2:0" MOUNTPOINT="" NAME="fd0"
	// MAJ:MIN="8:0" MOUNTPOINT="" NAME="sda"
	// MAJ:MIN="8:1" MOUNTPOINT="" NAME="sda1"
	// MAJ:MIN="8:2" MOUNTPOINT="/boot" NAME="sda2"a2"
	// MAJ:MIN="8:3" MOUNTPOINT="" NAME="sda3"
	// MAJ:MIN="253:0" MOUNTPOINT="/" NAME="rhel-root"-root"
	// MAJ:MIN="253:1" MOUNTPOINT="" NAME="rhel-swap"swap"
	// MAJ:MIN="253:2" MOUNTPOINT="/data" NAME="rhel-data"rhel-data"
	// MAJ:MIN="11:0" MOUNTPOINT="" NAME="sr0"
	lsblkTemp := strings.Split(lsblk, "\n")
	//fmt.Println(lsblkTemp)
	for _, i := range lsblkTemp[:len(lsblkTemp)-1] {
		lsblkArr := rex.FindStringSubmatch(i)
		// [MAJ:MIN="7:26" MOUNTPOINT="/rootfs/data/steamer/kubelet/pods/8b65b137-ae7e-11eb-95f2-0050569338f4/volumes/tiduyun~tsp/bank-aggregate-forward-0-1" NAME="loop26" 7:26 /rootfs/data/steamer/kubelet/pods/8b65b137-ae7e-11eb-95f2-0050569338f4/volumes/tiduyun~tsp/bank-aggregate-forward-0-1 loop26]
		if len(lsblkArr) != 0 {
			lsblkMap[lsblkArr[1]] = local.LsblkInfo{Name: lsblkArr[3], MountPoint: lsblkArr[2]}
		}
	}

	*Lsblkmap = lsblkMap
}

//func GetHostLsblk(Lsblkmap *map[string]local.LsblkInfo)  {
//    lsblkMap := make(map[string]local.LsblkInfo)
//}

// DiskUsage 查看挂载目录使用率
func DiskUsage(path string) local.DiskUsage {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return local.DiskUsage{}
	}
	AllSize := fs.Blocks * uint64(fs.Bsize)
	FreeSize := fs.Bfree * uint64(fs.Bsize)
	UsedSize := AllSize - FreeSize
	AllInodes := fs.Files
	FreeInodes := fs.Ffree
	UsedInodes := AllInodes - FreeInodes
	return local.DiskUsage{
		AllSize:    AllSize,
		FreeSize:   FreeSize,
		UsedSize:   UsedSize,
		AllInodes:  AllInodes,
		FreeInodes: FreeInodes,
		UsedInodes: UsedInodes,
	}
}
