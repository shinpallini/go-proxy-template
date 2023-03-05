package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerRequest(t *testing.T) {
	s := httptest.NewServer(initServer())
	log.Println("server url: ", s.URL)
	r, err := http.Get(s.URL + "/user")
	if err != nil {
		t.Errorf("http get err should be nil: %v", err)
		return
	}
	defer r.Body.Close()
	var j map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		t.Errorf("json decode err should be nil: %v", err)
		return
	}
	name, ok := j["name"].(string)
	if !ok {
		t.Error("name cannot assersion.")
	}
	if name != "Shin" {
		t.Errorf("result should be Shin, but %s", j["name"])
	}
	age, ok := j["age"].(int)
	if !ok {
		// t.Error("age cannot assersion.")
		a := j["age"].(float64)
		age = int(a)
	}
	if age != 31 {
		t.Errorf("result should be 31, but %f", j["age"])
	}
}

func TestHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ProxyHandler(w, r)
	var j map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &j)
	if err != nil {
		t.Errorf("json decode err should be nil: %v", err)
	}
	if j["hello"] != "world" {
		t.Errorf("result should be Shin, but %s", j["hello"])
	}
}
