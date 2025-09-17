package backend

import (
	"net/http"
)


func GetMainPage(homePageUrl string) (int,error){

	resp,err := http.Get(homePageUrl)
	defer resp.Body.Close()

	if err!=nil{
		return 0,err
	}


	return resp.StatusCode,nil
}
