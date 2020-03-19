package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Mount("/subdomain", subrouter())
	http.ListenAndServe(":3000", r)
}

func subrouter() *chi.Mux {
	r := chi.NewRouter()

	// RESTy routes for "articles" resource
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("subroute"))
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write(append([]byte("subrote with another thing: "), []byte(chi.URLParam(r, "id"))...))
			})
		})
	})
	return r
}
