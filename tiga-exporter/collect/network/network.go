package network

import (
    dclient "github.com/fsouza/go-dockerclient"
    "io/ioutil"
    "regexp"
    "strconv"
    "strings"
)

func GetNetworkData(pid int) map[string]dclient.NetworkStats {
    filePath := "/proc/" + strconv.Itoa(pid) + "/net/dev"
    networkStats := make(map[string]dclient.NetworkStats)
    rex := regexp.MustCompile(`(.*):\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`)
    file, err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }
    netLines := strings.Split(string(file), "\n")
    for _, i := range netLines {
        netTemp := strings.Index(i, ":")
        if netTemp != -1 {
            netInfoSource := rex.FindStringSubmatch(i)
            //fmt.Println(netInfoSource)
            netStats := dclient.NetworkStats{}

            // Receive RxBytes
            if rxBytes, err := strconv.ParseUint(netInfoSource[2], 10, 64); err != nil {
                netStats.RxBytes = 0
            } else {
                netStats.RxBytes = rxBytes
            }

            // Receive RxPackets
            if rxPackets, err := strconv.ParseUint(netInfoSource[3], 10, 64); err != nil {
                netStats.RxPackets = 0
            } else {
                netStats.RxPackets = rxPackets
            }

            // Receive RxErrors
            if rxErrors, err := strconv.ParseUint(netInfoSource[4], 10, 64); err != nil {
                netStats.RxErrors = 0
            } else {
                netStats.RxErrors = rxErrors
            }

            // Receive RxDropped
            if rxDropped, err := strconv.ParseUint(netInfoSource[5], 10, 64); err != nil {
                netStats.RxDropped = 0
            } else {
                netStats.RxDropped = rxDropped
            }

            // Transmit TxBytes
            if txBytes, err := strconv.ParseUint(netInfoSource[10], 10, 64); err != nil {
                netStats.TxBytes = 0
            } else {
                netStats.TxBytes = txBytes
            }

            // Transmit TxPackets
            if txPackets, err := strconv.ParseUint(netInfoSource[11], 10, 64); err != nil {
                netStats.TxPackets = 0
            } else {
                netStats.TxPackets = txPackets
            }

            // Transmit TxErrors
            if txErrors, err := strconv.ParseUint(netInfoSource[12], 10, 64); err != nil {
                netStats.TxErrors = 0
            } else {
                netStats.TxErrors = txErrors
            }

            // Transmit TxDropped
            if txDropped, err := strconv.ParseUint(netInfoSource[13], 10, 64); err != nil {
                netStats.TxDropped = 0
            } else {
                netStats.TxDropped = txDropped
            }
            //fmt.Println(netStats)
            networkStats[strings.TrimSpace(netInfoSource[1])] = netStats
        }
    }
    //fmt.Println(networkStats)
    return networkStats
}
