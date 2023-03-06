package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "http://apiserver:3000/update-token"
)

type retryableTransport struct {
	base     http.RoundTripper
	attempts int
	waitTime time.Duration
}

func (rt retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for count := 1; count <= rt.attempts; count++ {
		resp, err = rt.base.RoundTrip(req)
		log.Println(resp.Status)

		if !rt.shouldRetry(resp, err) {
			return resp, err
		}

		req.Header.Set("Authorization", "valid_token")

		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(rt.waitTime):
		}

	}
	return resp, err
}

func (rt retryableTransport) shouldRetry(resp *http.Response, err error) bool {
	return resp.StatusCode == http.StatusUnauthorized
}

func main() {
	client := &http.Client{
		Transport: &retryableTransport{
			base:     http.DefaultTransport,
			attempts: 3,
			waitTime: time.Duration(3) * time.Second,
		},
	}
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, resp.Body)
	defer resp.Body.Close()
	log.Println()
	for k, v := range resp.Header {
		log.Println("header: ", k)
		log.Println("value: ", v)
	}

}
