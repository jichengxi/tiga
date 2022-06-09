package docker

import (
	"context"
	dclient "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"net"
	"strings"
	"sync"
	"tiga-exporter/public/types/local"
	"time"
)

var (
	DCli DClient
	once sync.Once
)

type DClient struct {
	*dclient.Client
}

//func (t *DClient) NewClient() {
//	t.once.Do(func() {
//		//t.Client, t.clientErr = dclient.NewClient("unix:///var/run/docker.sock")
//		t.Client, t.clientErr = dclient.NewClientFromEnv()
//	})
//}

func NewClient() DClient {
	var dClient DClient
	once.Do(func() {
		//t.Client, t.clientErr = dclient.NewClient("unix:///var/run/docker.sock")
		client, err := dclient.NewClientFromEnv()
		if err != nil {
			panic(err)
		}
		dClient.Client = client
	})
	return dClient
}

func (t *DClient) DockerInfo() (dclient.DockerInfo, error) {
	info, err := t.Client.Info()
	if err != nil {
		return dclient.DockerInfo{}, err
	}
	return *info, nil
}

func (t *DClient) GetContainerInfoList(containerInfoList *[]*local.ContainerInfo) {
	// 获取宿主机ip
	var hostIp string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, i := range addrs {
		if ipNet, ok := i.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				if strings.HasPrefix(ipNet.IP.String(), "172.21") ||
					strings.HasPrefix(ipNet.IP.String(), "172.20") ||
					strings.HasPrefix(ipNet.IP.String(), "172.28") {
					hostIp = ipNet.IP.String()
				}
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	containers, err := t.Client.ListContainers(dclient.ListContainersOptions{Context: ctx})
	if err != nil {
		panic(err)
	}

	for _, i := range containers {
		container, err := t.Client.InspectContainerWithOptions(dclient.InspectContainerOptions{Context: ctx, ID: i.ID})
		if err != nil {
			panic(err)
		}
		cI := local.ContainerInfo{}
		cI.ContainerMounts = make(map[string]string)
		for _, mount := range container.Mounts {
			cI.ContainerMounts[mount.Source] = mount.Destination
		}
		var containerHasNet *dclient.Container
		var containerPodIp string
		containerHasNet = container
		if len(container.NetworkSettings.Networks) == 0 {
			networkMode := strings.Split(container.HostConfig.NetworkMode, ":")[1]
			containerHasNet, err = t.Client.InspectContainerWithOptions(dclient.InspectContainerOptions{
				Context: ctx,
				ID:      networkMode,
			})
			if err != nil {
				glog.Error(err)
			}
		}

		for netKey, netval := range containerHasNet.NetworkSettings.Networks {
			if netKey == "host" {
				containerPodIp = hostIp
				break
			} else {
				if strings.HasPrefix(netval.IPAddress, "172.20") ||
					strings.HasPrefix(netval.IPAddress, "172.21") ||
					strings.HasPrefix(netval.IPAddress, "172.28") {
					containerPodIp = netval.IPAddress
					break
				}
			}
		}

		cI.ContainerId = container.ID
		cI.Pid = container.State.Pid
		cI.Name = strings.TrimPrefix(container.Name, "/")
		cI.HostName = container.Config.Hostname
		cI.Image = container.Config.Image
		cI.KubernetesDockerType = container.Config.Labels["io.kubernetes.docker.type"]
		cI.KubernetesPodNamespace = container.Config.Labels["io.kubernetes.pod.namespace"]
		cI.KubernetesPodName = container.Config.Labels["io.kubernetes.pod.name"]
		cI.KubernetesContainerName = container.Config.Labels["io.kubernetes.container.name"]
		cI.KubernetesPodUid = container.Config.Labels["io.kubernetes.pod.uid"]
		cI.ContainerLabelPodIp = containerPodIp
		cI.KubernetesSandboxId = container.Config.Labels["io.kubernetes.sandbox.id"]
		cI.MergedDir = container.GraphDriver.Data["MergedDir"]

		*containerInfoList = append(*containerInfoList, &cI)
	}
	//return containerInfoList
}

func init() {
	DCli = NewClient()
}
