package remote

import (
	"fmt"
	"github.com/golang/glog"
	"regexp"
	"strconv"
	"strings"
	"tiga-exporter/public/types/remote"
	"tiga-exporter/tools"
)

// ParseRemoteConfig 从配置文件读取配置，并解析输出结果
func ParseRemoteConfig(remoteConfig *tools.Remote, metricType *map[string]*remote.MetricTitle, metricInfoData *remote.MetricData) {
	// 如果外部url等于0，就退出不采集
	if len(remoteConfig.Url) == 0 {
		return
	}
	// 如果监控项没有，就退出不采集
	metricList := remoteConfig.Metrics
	blackList := remoteConfig.Blacklist
	if len(metricList) == 0 {
		return
	}

	// type string是监控项标题
	metricTitle := make(map[string]*remote.MetricTitle)
	// Info1 string是监控项
	metricInfoByLabelsMap := make(map[string]remote.InfoByLabelsStruct)
	// Info2 string是监控项
	metricInfoNotValueMap := make(map[string]remote.MetricValue)

	// metricData
	metricData := remote.MetricData{}
	// 带标签行正则匹配
	rex1 := regexp.MustCompile(`(.*){(.*)} (.*)`)
	// 无标签行正则匹配
	rex2 := regexp.MustCompile(`(.*) (.*)`)
	// 标签细节匹配
	rexLabel := regexp.MustCompile(`(.*)="(.*)"`)

	for urlKey, urlVal := range remoteConfig.Url {
		fmt.Println(urlKey, urlVal)
		httpResp := tools.HttpGet(urlVal)
		respBodyArr := strings.Split(string(httpResp), "\n")

		// 循环网页每一行
		for _, row := range respBodyArr {
			// 解析help行
			if strings.HasPrefix(row, "# HELP") {
				tempJ := strings.SplitN(row, " ", 4)
				metricNewName := urlKey + "_" + tempJ[2]
				if tools.IsExist(metricNewName, metricList) {
					metricTitle[metricNewName] = &remote.MetricTitle{Name: metricNewName, Help: tempJ[3]}
				}
				continue
			}

			// 解析type行
			// Todo 完全可以不用， 但是为了照顾不规范的Help
			if strings.HasPrefix(row, "# TYPE") {
				tempJ := strings.Split(row, " ")
				metricNewName := urlKey + "_" + tempJ[2]
				if tools.IsExist(metricNewName, metricList) {
					if _, ok := metricTitle[metricNewName]; !ok {
						metricTitle[metricNewName] = &remote.MetricTitle{Name: metricNewName}
					}

				}
				continue
			}

			// 解析带标签的行
			rep1 := rex1.FindStringSubmatch(row)
			if len(rep1) != 0 {
				metricNewName := urlKey + "_" + rep1[1]
				// 判断当前行的值在不在配置文件中
				if tools.IsExist(metricNewName, metricList) {
					labelsKv := make(map[string]string)
					LabelsSplits := strings.Split(rep1[2], ",")
					for _, LabelsSplit := range LabelsSplits {
						repLabel := rexLabel.FindStringSubmatch(LabelsSplit)
						labelsKv[repLabel[1]] = repLabel[2]
					}
					if BlackListExist(blackList, labelsKv) {
						continue
					}

					InfoByLabels := metricInfoByLabelsMap[metricNewName]
					var labelsValue []string

					if len(InfoByLabels.Labels) == 0 {
						for labelK, labelV := range labelsKv {
							InfoByLabels.Labels = append(InfoByLabels.Labels, labelK)
							labelsValue = append(labelsValue, labelV)
						}
					} else {
						for _, labelV := range labelsKv {
							labelsValue = append(labelsValue, labelV)
						}
					}
					metricValue, err := strconv.ParseFloat(rep1[3], 64)
					if err != nil {
						glog.Error(err.Error())
						metricValue = -1
					}

					InfoByLabels.MetricValue = append(InfoByLabels.MetricValue, remote.MetricValue{
						LabelsValue: labelsValue,
						Value:       metricValue,
					})

					//x := metricInfo11[metricNewName]
					//var labelsValue []string
					//
					//// TODO: 待优化 写的太垃圾了
					//if len(x.Labels) == 0 {
					//	for labelK, labelV := range labelsKv {
					//		x.Labels = append(x.Labels, labelK)
					//		labelsValue = append(labelsValue, labelV)
					//	}
					//} else {
					//	for _, labelV := range labelsKv {
					//		labelsValue = append(labelsValue, labelV)
					//	}
					//}
					//metricInfo11[metricNewName] = x
					//value, err := strconv.ParseFloat(rep1[3], 64)
					//if err != nil {
					//	panic(err)
					//}
					//metricInfo11Info := metricInfo11[metricNewName]
					//metricInfo11Info.MetricInfo1s = append(metricInfo11Info.MetricInfo1s, remote.MetricInfo1{
					//	LabelsValue: labelsValue,
					//	Value:       value,
					//})
					//metricInfo11[metricNewName] = metricInfo11Info
				}
				continue
			}

			// 解析不带参数的行
			rep2 := rex2.FindStringSubmatch(row)
			if len(rep2) != 0 {
				metricNewName := urlKey + "_" + rep2[1]
				if tools.IsExist(metricNewName, metricList) {
					metricValue, err := strconv.ParseFloat(rep2[2], 64)
					if err != nil {
						glog.Error(err.Error())
						metricValue = -1
					}
					metricInfoNotValueMap[metricNewName] = remote.MetricValue{Value: metricValue}
				}
				continue
			}
		}

	}
	// 有参数的数据在info1
	metricData.InfoByLabelsMap = metricInfoByLabelsMap
	// 没参数的数据在info2
	metricData.InfoNotLabelsMap = metricInfoNotValueMap
	*metricType = metricTitle
	*metricInfoData = metricData
}
