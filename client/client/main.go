package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://proxy:3001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Status)
	io.Copy(os.Stdout, resp.Body)
	defer resp.Body.Close()
	log.Println()
	for k, v := range resp.Header {
		log.Println("header: ", k)
		log.Println("value: ", v)
	}

}
