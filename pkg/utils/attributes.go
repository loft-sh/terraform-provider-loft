package utils

func HasKeys(in map[string]interface{}) bool {
	return len(in) > 0
}

func HasValue(in []interface{}) bool {
	return len(in) > 0 && in[0] != nil
}

func AttributesToMap(rawMap map[string]interface{}) map[string]string {
	strMap := map[string]string{}
	for k, v := range rawMap {
		strMap[k] = v.(string)
	}

	if len(strMap) == 0 {
		return nil
	}

	return strMap
}

func MapToAttributes(rawMap map[string]string) map[string]interface{} {
	attr := map[string]interface{}{}
	for k, v := range rawMap {
		attr[k] = v
	}

	if len(attr) > 0 {
		return nil
	}

	return attr
}
