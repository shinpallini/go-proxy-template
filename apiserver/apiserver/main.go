package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiredIn   int    `json:"expired_in"`
}
type Sample struct {
	Aaa string `json:"aaa"`
	Bbb string `json:"bbb"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"hello": "world"}`))
	})
	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		user := &User{
			Name:        "Shin",
			Age:         31,
			Description: "Men",
		}
		b, err := json.Marshal(user)
		if err != nil {
			log.Println("Cannot marshal struct. ", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})
	r.Get("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
	r.Get("/update-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "valid_token" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("validation OK!"))
			return
		} else {
			w.Header().Add("WWW-Authenticate", `Bearer realm="SECRET AREA"`)
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
	r.Post("/token", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			log.Printf("Key: %#v\nValue: %#v", k, v)
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
		}
		var sample Sample
		err = json.Unmarshal(body, &sample)
		if err != nil {
			log.Println(err)
		}
		log.Printf("body content: %#v", sample)
		token := &Token{
			AccessToken: "1234567890",
			ExpiredIn:   -1,
		}
		b, err := json.Marshal(token)
		if err != nil {
			log.Println("Cannot marshal token: ", err)
		}
		w.Write(b)
	})
	http.ListenAndServe(":3000", r)
}
