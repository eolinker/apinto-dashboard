package v2

func dealAPPAppend(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		if key == "additional" || key == "auth" {
			vs, ok := value.([]interface{})
			if !ok {
				continue
			}
			params := make([]interface{}, 0, len(vs))
			for _, v := range vs {
				val, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				vv, ok := val["config"]
				if ok {
					params = append(params, vv)
				}
			}
			result[key] = params
			continue
		}
	}
	return result
}
