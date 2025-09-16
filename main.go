package main 

import (
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func getUserName(w http.ResponseWriter, r * http.Request){
	userName := chi.URLParam(r,"userName")

	w.Write([]byte("Hello "+userName))
}


func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("root."))
	})
	
	r.Get("/{userName}",getUserName)
	
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
