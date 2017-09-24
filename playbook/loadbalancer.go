package playbook

type LoadBalancer Component

func (l *LoadBalancer) getEndpoint(name string) string {
	vips, find := l.Property["vips"]
	if !find {
		return ""
	}

	vs, ok := vips.([]interface{})
	if !ok {
		return ""
	}

	for _, vip := range vs {
		v, ok := vip.(map[string]interface{})
		if !ok {
			return ""
		}

		n, ok := v["type"].(string)
		if !ok {
			return ""
		}

		if n == name {
			ip, ok := v["vip"].(string)
			if !ok {
				return ""
			}

			return ip
		}
	}

	return ""
}
