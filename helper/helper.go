package helper

func CompareMapKey(X, Y map[string]struct{}) []string {
	result := []string{}

	for key, _ := range X {
		if _, ok := Y[key]; !ok {
			result = append(result, key)
		}
	}

	return result
}
