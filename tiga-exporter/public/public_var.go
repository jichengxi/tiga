package public

import "tiga-exporter/public/types/remote"

var (
	Rootfs         string
	MetricTitleMap map[string]*remote.MetricTitle
	MetricData     remote.MetricData
)
