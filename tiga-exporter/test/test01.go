package main

import (
    "fmt"
    "github.com/prometheus/procfs"
)

func main() {
    fs, err := procfs.NewFS("/proc")
    if err != nil {
        panic(err)
    }
    //cpuInfo, _ := fs.CPUInfo()
    cpuProcessStats, _ := fs.Stat()
    fmt.Println("cpuProcessStats: ", cpuProcessStats)
    memoryStats, _ := fs.Meminfo()
    fmt.Println("memoryStats: ", memoryStats)
    netStats, _ := fs.NetDev()
    fmt.Println("netStats: ", netStats)
    tcpStats, _ := fs.NetTCP()
    fmt.Println("tcpStats: ", tcpStats)
    udpStats, _ := fs.NetUDP()
    fmt.Println("udpStats: ", udpStats)
    ipvsStats, _ := fs.IPVSStats()
    fmt.Println("ipvsStats: ", ipvsStats)
    ipvsBackendStats, _ := fs.IPVSBackendStatus()
    fmt.Println("ipvsBackendStats: ", ipvsBackendStats)





}
