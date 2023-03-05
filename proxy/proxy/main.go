package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	baseURL = "http://apiserver:3000/"
)

func initServer() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Get("/", ProxyHandler)
	r.Route("/{path:^(.*)?$}", func(r chi.Router) {
		r.Get("/", ProxyHandler)
	})
	return r
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	var resp *http.Response
	var err error
	if path == "" {
		resp, err = http.Get(baseURL)
	} else {
		resp, err = http.Get(baseURL + path)
	}
	if err != nil {
		log.Println("client get error: ", err)
	}
	log.Printf("Response status: %#v", resp.Status)
	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	defer resp.Body.Close()
}

func main() {
	log.Println("open at :3001")
	log.Println(http.ListenAndServe(":3001", initServer()))
}
