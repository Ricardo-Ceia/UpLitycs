package main 

import (
    "net/http"
    "fmt"
		"log"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/v5/middleware"
		"time"
    "uplytics/db"
)

func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	conn := db.OpenDB()
	defer conn.Close()
	
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
