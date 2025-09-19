package status_checker

import (
	"net/http"
)

func GetMainPage(homePageUrl string) (int, error) {

	resp, err := http.Get(homePageUrl)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}
