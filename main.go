package main 

import (
    "net/http"
		"log"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/v5/middleware"
		"time"
    "uplytics/db"
		"uplytics/backend"
)

func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	conn := db.OpenDB()
	defer conn.Close()

	testUrl := "https://x.com/home"
	for i:=0;i<4;i++{
		status,err := backend.GetMainPage(testUrl)
		
		if err!=nil{
			log.Println(err)
		}

		log.Println("Status code returned from:%s was %d",testUrl,status)
		time.Sleep(5*time.Second)
	} 
	
	r.NotFound(func(w http.ResponseWriter,r *http.Request){
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	http.ListenAndServe(":3333",r)
}
