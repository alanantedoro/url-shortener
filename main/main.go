package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshortener "github.com/alanantedoro/url-shortener"
)


func main() {
	var yamlFile, jsonFile string
	 flag.StringVar(&yamlFile,"yaml","","a yaml file in the format of 'path and url'")
	flag.StringVar(&jsonFile,"json","","a json file in the format of 'path and url'")
	flag.Parse()
	mux := defaultMux()



	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)
	
	if(yamlFile != ""){
		yamlContent, err := os.ReadFile(yamlFile)
		if err != nil {
			panic(err)
		}
		yamlHandler, err := urlshortener.YAMLHandler([]byte(yamlContent), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	} else if jsonFile != ""{
		jsonContent, err := os.ReadFile(jsonFile)
		if err != nil {
			panic(err)
		}
		
		jsonHandler, err := urlshortener.JSONHandler(jsonContent, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	}

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

