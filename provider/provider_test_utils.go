package provider

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			str := v.(string)
			vs = append(vs, &str)
		}
	}
	return vs
}
