package utils

func HasValue(in []interface{}) bool {
	return len(in) > 0 && in[0] != nil
}

func AttributesToMap(rawMap map[string]interface{}) map[string]string {
	var strMap map[string]string
	for k, v := range rawMap {
		strMap[k] = v.(string)
	}
	return strMap
}

func MapToAttributes(rawMap map[string]string) map[string]interface{} {
	var attr map[string]interface{}
	for k, v := range rawMap {
		attr[k] = v
	}
	return attr
}
