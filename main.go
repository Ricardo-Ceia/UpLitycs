package main 

import (
	//"context"
	"net/http"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func MyHandler(w http.ResponseWriter,r *http.Request){
	user := r.Context().Value("user").(string)

	w.Write([]byte(fmt.Sprintf("hi %s",user)))
}

func sendHTMLHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","text/html")
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
    	<title>Hello</title>
		</head>
		<body>
    	<h1>Hello, World!</h1>
    	<p>This is a simple HTML string.</p>
		</body>
		</html>`))
}

func main(){
	r := chi.NewRouter()
	r.Use(middleware.Compress(5,"text/html","text/css"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/",sendHTMLHandler)
	
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
