package main

import (
    "fmt"
    gofstab "github.com/deniswernert/go-fstab"
    "os"
    "path/filepath"
    "strings"
    "syscall"
)

func DirSize(path string) (int64, error) {
   var size int64
   err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
       if !info.IsDir() {
           size += info.Size()
       }
       return err
   })
   return size, err
}


func DiskUsage(path string) (uint64, uint64, uint64) {
    fs := syscall.Statfs_t{}
    err := syscall.Statfs(path, &fs)
    if err != nil {
        return 0, 0, 0
    }
    All := fs.Blocks * uint64(fs.Bsize)
    Free := fs.Bfree * uint64(fs.Bsize)
    Used := All - Free

    return All, Free, Used
}

func main() {
    procs ,_ := gofstab.ParseProc()
    //mounts, _ := gofstab.ParseSystem()
    //fmt.Println(procs)
    //fmt.Println("-------------")
    //fmt.Println(mounts)
    for _, val := range procs {
        if val.Spec == "sysfs" || val.Spec == "proc" || val.Spec == "devtmpfs" || val.Spec == "securityfs" ||
            val.Spec == "devpts" || val.Spec == "cgroup" || val.Spec == "configfs" || val.Spec == "systemd-1" ||
            val.Spec == "debugfs" || val.Spec == "mqueue" || val.Spec == "hugetlbfs" || val.Spec == "binfmt_misc" ||
            val.Spec == "fusectl" || val.Spec == "/sys/fs/bpf" ||val.File == "/dev/shm" || val.File == "/run" ||
            val.File == "/sys/fs/pstore" || val.File == "/sys/fs/cgroup" || val.Spec == "overlay" ||
            strings.Contains(val.File, "/run/user/") {
            continue
        }
        fmt.Println(val)
    }
    _, _, c := DiskUsage("/data/docker/overlay2/055a7c89fa43473dca95ad851439f75db045cfd4824c7386a5a8bb8f89d6a189/merged")
    fmt.Println(c)
}