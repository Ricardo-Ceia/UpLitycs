package main 

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/middlware"
)

func main(){
	r := chi.NewRouter()
	r.Use(middlware.Logger)
	r.Use(middlware.Recoverer)

	r.Get("/",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte(
			"root."
		))
	})

	http.lisenAndServe(":3333",r)
}
