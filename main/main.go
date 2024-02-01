package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	urlshortener "github.com/alanantedoro/url-shortener"
)

func main() {
	sourceFile := flag.String("yaml","urls.yaml","a yaml file in the format of 'path and url'")
	flag.Parse()
	mux := defaultMux()

	yamlContent, err := ioutil.ReadFile(*sourceFile)
	if err != nil {
		panic(err)
	}

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	yamlHandler, err := urlshortener.YAMLHandler([]byte(yamlContent), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

