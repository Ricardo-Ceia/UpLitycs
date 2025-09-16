package main 

import (
	"context"
	"net/http"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func getUserName(w http.ResponseWriter, r * http.Request){
	userName := chi.URLParam(r,"userName")

	w.Write([]byte("Hello "+userName))
}

func MyMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    ctx := context.WithValue(r.Context(), "user", "123")
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func MyHandler(w http.ResponseWriter,r *http.Request){
	user := r.Context().Value("user").(string)

	w.Write([]byte(fmt.Sprintf("hi %s",user)))
}

func main(){
	r := chi.NewRouter()
	r.Use(MyMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/",MyHandler)
	
	//r.Get("/{userName}",getUserName)
	
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
