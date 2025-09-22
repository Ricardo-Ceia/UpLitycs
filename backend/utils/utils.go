package utils

import (
	"strings"
)

var statusMap = map[int]string{
	200: "up",
	201: "up",
	301: "redirect",
	302: "redirect",
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	500: "server error",
	502: "bad gateway",
	503: "service unavailable",
	504: "gateway timeout",
}

func MapStatusCode(code int) string {
	if status, ok := statusMap[code]; ok {
		return status
	}
	return "unknown"
}

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
