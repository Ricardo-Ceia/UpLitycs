package utils

import (
	"strings"
)

func CheckUsername(username string) bool {

	if len(username) < 3 || len(username) > 20 {
		return false
	}
	return true
}

func CheckURLFormat(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

func CheckAlerts(alerts string) bool {
	if alerts != "y" && alerts != "n" && alerts != "yes" && alerts != "no" {
		return false
	}
	return true
}
