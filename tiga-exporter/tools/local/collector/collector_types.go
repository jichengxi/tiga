package collector

type metricValue struct {
	value  float64
	labels []string
}

type metricValues []metricValue
