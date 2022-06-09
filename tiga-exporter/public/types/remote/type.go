package remote

type (
	MetricTitle struct {
		Name string
		Help string
	}

	MetricValue struct {
		LabelsValue []string
		Value       float64
	}
	InfoByLabelsStruct struct {
		Labels      []string
		MetricValue []MetricValue
	}

	//MetricInfo2 struct {
	//	Value float64
	//}

	MetricData struct {
		InfoByLabelsMap  map[string]InfoByLabelsStruct
		InfoNotLabelsMap map[string]MetricValue
	}
)
