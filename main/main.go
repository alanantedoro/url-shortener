package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	urlshortener "github.com/alanantedoro/url-shortener"
	"github.com/boltdb/bolt"
)


func main() {
	var yamlFile, jsonFile, boltDB string
	 flag.StringVar(&yamlFile,"yaml","","a yaml file in the format of 'path and url'")
	flag.StringVar(&jsonFile,"json","","a json file in the format of 'path and url'")
	flag.StringVar(&boltDB,"bolt","","a json file in the format of 'path and url'")

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
	} else if(boltDB != ""){
		db, err := bolt.Open(boltDB, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Just for setting up a record in the db in the first run
		// db.Update(func(tx *bolt.Tx) error {
		// 	b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		// 	if err != nil {
		// 		return fmt.Errorf("create bucket: %s", err)
		// 	}
		
		// 	if err := b.Put([]byte("/bolt"), []byte("https://pkg.go.dev/github.com/boltdb/bolt")); err != nil {
		// 		return fmt.Errorf("put data: %s", err)
		// 	}
		
		// 	return nil
		// })

		boltHandler := urlshortener.BOLThanlder(db,mapHandler)
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", boltHandler)

		defer db.Close()
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


