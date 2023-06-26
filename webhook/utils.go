package webhook

import "strings"

// ToHeaders converts a slice of strings to a map of strings.
func ToHeaders(headers []string) map[string]string {
	h := map[string]string{}
	for _, header := range headers {
		kv := strings.Split(header, "=")
		if len(kv) != 2 {
			continue
		}

		h[kv[0]] = kv[1]
	}

	return h
}
