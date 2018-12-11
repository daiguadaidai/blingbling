package common

// 合并多个map
func CombindMapStrInt(maps ...map[string]int) map[string]int {
	dMap := make(map[string]int)
	for _, m := range maps {
		for k, v := range m {
			dMap[k] = v
		}
	}
	return dMap
}
