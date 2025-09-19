package backend

import (
	"net/http"
	"log"
)
func StartOnboardingHandler(w http.ResponseWriter, r *http.Request){
	//TODO ADD LOGIN VERIFICATION BEFORE REDIRECTING USER
	log.Println("Reaching the endpoint")
	http.Redirect(w,r,"/onboarding",http.StatusFound)
}
