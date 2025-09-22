package status_checker

import "net/http"

func getPageStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func GetPagesStatus(pages []string) ([]int, error) {
	results := make([]int, 0, len(pages))
	for _, p := range pages {
		code, err := getPageStatus(p)
		if err != nil {
			results = append(results, 0) // treat as down
			continue
		}
		results = append(results, code)
	}
	return results, nil
}
