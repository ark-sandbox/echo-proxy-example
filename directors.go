package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

func SimpleForwarder(tgtUrl string, ctx echo.Context, args []string) func(*http.Request) {
	tgtUrl = "http://www.google.co.in"
	target, err := url.Parse(tgtUrl)
	if err != nil {
		log.Fatal(err)
	}
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = target.Path, target.RawPath
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return director
}
