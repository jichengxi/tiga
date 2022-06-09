package remote

func BlackListExist(blackList map[string]string, labelKv map[string]string) bool {
	for m, n := range blackList {
		if v, ok := labelKv[m]; ok {
			if v == n {
				return true
			}
		}
	}
	return false
}
