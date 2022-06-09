package tools

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"time"
)

type Remote struct {
	Enable    bool              `json:"enable"`
	Url       map[string]string `json:"url"`
	Metrics   []string          `json:"metrics"`
	Blacklist map[string]string `json:"blacklist"`
}

type LocalMetricConfig struct {
	Enable  bool     `json:"enable"`
	Metrics []string `json:"metrics"`
}

type Local struct {
	Container LocalMetricConfig `json:"container"`
	Node      LocalMetricConfig `json:"node"`
}

type TigaConfig struct {
	Remote Remote `json:"remote"`
	Local  Local  `json:"local"`
}

func IsExist(a string, b []string) bool {
	for _, i := range b {
		if i == a {
			return true
		}
	}
	return false
}

// CompareSlice 判断两个数组是否一直
func CompareSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, i := range a {
		res := IsExist(i, b)
		if res {
			continue
		} else {
			return false
		}
	}
	for _, i := range b {
		res := IsExist(i, a)
		if res {
			continue
		} else {
			return false
		}
	}
	return true
}

// HttpGet http get 请求
func HttpGet(url string) []byte {
	client := &http.Client{Timeout: 120 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err.Error())
	}
	response, err := client.Do(request)
	if err != nil {
		glog.Error(err.Error())
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Error(err.Error())
	}
	return data
}

// RequestTigaConfig 请求controller并解析配置文件到结构体
func RequestTigaConfig(url string) TigaConfig {
	//read, err := ioutil.ReadFile("remote-config/config.yaml")
	//read, err := ioutil.ReadFile("config.yaml")
	//read, err := utils.FileRead("remote-config/config.yaml")
	//read, err := utils.FileRead("config.yaml")
	//if err != nil {
	//	glog.Fatal(err)
	//}
	controllerConfig := HttpGet(url)
	glog.V(7).Infoln("读取到的配置信息: \n", string(controllerConfig))
	var tConfig TigaConfig
	err := json.Unmarshal(controllerConfig, &tConfig)
	if err != nil {
		glog.Fatal(err)
	}
	defer glog.Flush()
	return tConfig
}
