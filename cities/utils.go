package cities

import "strings"

func ContainsCity(alertData, city string) bool {
	return strings.Contains(strings.ToLower(alertData), strings.ToLower(city))
}

func ContainsCityArray(cities []string, city string) bool {
	cityLower := strings.ToLower(city)
	for _, c := range cities {
		if strings.Contains(strings.ToLower(c), cityLower) {
			return true
		}
	}
	return false
}
